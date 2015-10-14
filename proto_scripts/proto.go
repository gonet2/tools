package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/codegangsta/cli"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"text/template"
	"unicode"
)

const (
	TK_SYMBOL = iota
	TK_STRUCT_BEGIN
	TK_STRUCT_END
	TK_DATA_TYPE
	TK_ARRAY
	TK_EOF
)

var (
	datatypes = map[string]bool{
		"integer": true,
		"string":  true,
		"bytes":   true,
		"byte":    true,
		"boolean": true,
		"bool":    true,
		"float":   true,
		"float32": true,
		"float64": true,
		"uint8":   true,
		"int8":    true,
		"uint16":  true,
		"int16":   true,
		"uint32":  true,
		"int32":   true,
		"uint64":  true,
		"int64":   true,
		"long":    true,
		"short":   true,
	}
)

var (
	TOKEN_EOF = &token{typ: TK_EOF}
)

type (
	func_info struct {
		T string `json:"t"` // type
		R string `json:"r"` // read
		W string `json:"w"` // write
	}
	lang_type struct {
		Go func_info `json:"go"` // golang
		Cs func_info `json:"cs"` // c#
	}
)

type (
	field_info struct {
		Name  string
		Typ   string
		Array bool
	}
	struct_info struct {
		Name   string
		Fields []field_info
	}
)

type token struct {
	typ     int
	literal string
	r       rune
}

func syntax_error(p *Parser) {
	log.Println("syntax error @line:", p.lexer.lineno)
	log.Println(">> \033[1;31m", p.lexer.lines[p.lexer.lineno-1], "\033[0m <<")
	os.Exit(-1)
}

type Lexer struct {
	reader *bytes.Buffer
	lines  []string
	lineno int
}

func (lex *Lexer) init(r io.Reader) {
	bts, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	// 按行读入源码
	scanner := bufio.NewScanner(bytes.NewBuffer(bts))
	for scanner.Scan() {
		lex.lines = append(lex.lines, scanner.Text())
	}

	// 清除注释
	re := regexp.MustCompile("(?m:^#(.*)$)")
	bts = re.ReplaceAllLiteral(bts, nil)
	lex.reader = bytes.NewBuffer(bts)
	lex.lineno = 1
}

func (lex *Lexer) next() (t *token) {
	defer func() {
		//log.Println(t)
	}()
	var r rune
	var err error
	for {
		r, _, err = lex.reader.ReadRune()
		if err == io.EOF {
			return TOKEN_EOF
		} else if unicode.IsSpace(r) {
			if r == '\n' {
				lex.lineno++
			}
			continue
		}
		break
	}

	if r == '=' {
		for k := 0; k < 2; k++ { // check "==="
			r, _, err = lex.reader.ReadRune()
			if err == io.EOF {
				return TOKEN_EOF
			}
			if r != '=' {
				lex.reader.UnreadRune()
				return &token{typ: TK_STRUCT_BEGIN}
			}
		}
		return &token{typ: TK_STRUCT_END}
	} else if unicode.IsLetter(r) {
		var runes []rune
		for {
			runes = append(runes, r)
			r, _, err = lex.reader.ReadRune()
			if err == io.EOF {
				break
			} else if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' {
				continue
			} else {
				lex.reader.UnreadRune()
				break
			}
		}

		t := &token{}
		t.literal = string(runes)
		if datatypes[t.literal] {
			t.typ = TK_DATA_TYPE
		} else if t.literal == "array" {
			t.typ = TK_ARRAY
		} else {
			t.typ = TK_SYMBOL
		}

		return t
	} else {
		log.Fatal("lex error @line:", lex.lineno)
	}
	return nil
}

func (lex *Lexer) eof() bool {
	for {
		r, _, err := lex.reader.ReadRune()
		if err == io.EOF {
			return true
		} else if unicode.IsSpace(r) {
			if r == '\n' {
				lex.lineno++
			}
			continue
		} else {
			lex.reader.UnreadRune()
			return false
		}
	}
}

//////////////////////////////////////////////////////////////
type Parser struct {
	lexer *Lexer
	info  []struct_info
}

func (p *Parser) init(lex *Lexer) {
	p.lexer = lex
}

func (p *Parser) match(typ int) *token {
	t := p.lexer.next()
	if t.typ != typ {
		syntax_error(p)
	}
	return t
}

func (p *Parser) expr() bool {
	if p.lexer.eof() {
		return false
	}
	info := struct_info{}

	t := p.match(TK_SYMBOL)
	info.Name = t.literal

	p.match(TK_STRUCT_BEGIN)
	p.fields(&info)
	p.info = append(p.info, info)
	return true
}

func (p *Parser) fields(info *struct_info) {
	for {
		t := p.lexer.next()
		if t.typ == TK_STRUCT_END {
			return
		}
		if t.typ != TK_SYMBOL {
			syntax_error(p)
		}

		field := field_info{Name: t.literal}
		t = p.lexer.next()
		if t.typ == TK_ARRAY {
			field.Array = true
			t = p.lexer.next()
		}

		if t.typ == TK_DATA_TYPE || t.typ == TK_SYMBOL {
			field.Typ = t.literal
		} else {
			syntax_error(p)
		}

		info.Fields = append(info.Fields, field)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Protocol Data Structure Generator"
	app.Usage = "handle proto.txt"
	app.Authors = []cli.Author{{Name: "xtaci"}, {Name: "ycs"}}
	app.Version = "1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "file,f", Value: "./proto.txt", Usage: "input proto.txt file"},
		cli.StringFlag{Name: "template,t", Value: "./templates/server/proto.tmpl", Usage: "template file"},
	}
	app.Action = func(c *cli.Context) {
		// parse
		file, err := os.Open(c.String("file"))
		if err != nil {
			log.Fatal(err)
		}
		lexer := Lexer{}
		lexer.init(file)
		p := Parser{}
		p.init(&lexer)
		for p.expr() {
		}

		// load function mapping
		var funcs map[string]lang_type
		f, err := os.Open("func_map.json")
		if err != nil {
			log.Fatal(err)
		}
		if err := json.NewDecoder(f).Decode(&funcs); err != nil {
			log.Fatal(err)
		}

		// use template to generate final output
		funcMap := template.FuncMap{
			"goType": func(t string) string {
				if v, ok := funcs[t]; ok {
					return v.Go.T
				} else {
					return ""
				}
			},
			"goRead": func(t string) string {
				if v, ok := funcs[t]; ok {
					return v.Go.R
				} else {
					return ""
				}
			},
			"goWrite": func(t string) string {
				if v, ok := funcs[t]; ok {
					return v.Go.W
				} else {
					return ""
				}
			},
			"csType": func(t string) string {
				if v, ok := funcs[t]; ok {
					return v.Cs.T
				} else {
					return ""
				}
			},
			"csRead": func(t string) string {
				if v, ok := funcs[t]; ok {
					return v.Cs.R
				} else {
					return ""
				}
			},
			"csWrite": func(t string) string {
				if v, ok := funcs[t]; ok {
					return v.Cs.W
				} else {
					return ""
				}
			},
		}
		tmpl, err := template.New("proto.tmpl").Funcs(funcMap).ParseFiles(c.String("template"))
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(os.Stdout, p.info)
		if err != nil {
			log.Fatal(err)
		}
	}
	app.Run(os.Args)
}
