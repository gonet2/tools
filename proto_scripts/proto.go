package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"unicode"
)

const (
	SYMBOL = iota
	STRUCT_BEGIN
	STRUCT_END
	DATA_TYPE
	ARRAY_TYPE
)

var (
	datatypes = map[string]bool{
		"integer": true,
		"string":  true,
		"bytes":   true,
		"byte":    true,
		"boolean": true,
		"float":   true,
	}
)

var (
	SYNTAX_ERROR = errors.New("syntax error")
)

type field_info struct {
	name  string
	typ   string
	array bool
}

type struct_info struct {
	name   string
	fields []field_info
}

type token struct {
	typ     int
	literal string
	r       rune
}

func syntax_error(p *Parser) {
	log.Fatal("syntax error @line:", p.lexer.lineno)
}

type Lexer struct {
	reader *bytes.Buffer
	lineno int
}

func (lex *Lexer) init(r io.Reader) {
	bts, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println(err)
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
			return nil
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
				return nil
			}
			if r != '=' {
				lex.reader.UnreadRune()
				return &token{typ: STRUCT_BEGIN}
			}
		}
		return &token{typ: STRUCT_END}
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
			t.typ = DATA_TYPE
		} else if t.literal == "array" {
			t.typ = ARRAY_TYPE
		} else {
			t.typ = SYMBOL
		}

		return t
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

	t := p.match(SYMBOL)
	info.name = t.literal

	p.match(STRUCT_BEGIN)
	p.fields(&info)
	p.info = append(p.info, info)
	return true
}

func (p *Parser) fields(info *struct_info) {
	for {
		t := p.lexer.next()
		if t.typ == STRUCT_END {
			return
		}
		if t.typ != SYMBOL {
			syntax_error(p)
		}

		field := field_info{name: t.literal}
		t = p.lexer.next()
		if t.typ == ARRAY_TYPE {
			field.array = true
			t = p.match(SYMBOL)
			field.typ = t.literal
		} else if t.typ == DATA_TYPE || t.typ == SYMBOL {
			field.typ = t.literal
		} else {
			syntax_error(p)
		}

		info.fields = append(info.fields, field)
	}
}

func main() {
	lexer := Lexer{}
	lexer.init(os.Stdin)
	p := Parser{}
	p.init(&lexer)
	for p.expr() {
	}

	log.Println(p.info)
}
