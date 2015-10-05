package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"unicode"
)

const (
	TK_TYPE = iota
	TK_NAME
	TK_PAYLOAD
	TK_COLON
	TK_STRING
	TK_NUMBER
	TK_EOF
)

var (
	keywords = map[string]int{
		"packet_type": TK_TYPE,
		"name":        TK_NAME,
		"payload":     TK_PAYLOAD,
	}
)

var (
	SYNTAX_ERROR = errors.New("syntax error")
)

var (
	TOKEN_EOF   = &token{typ: TK_EOF}
	TOKEN_COLON = &token{typ: TK_COLON}
)

type api_expr struct {
	packet_type int
	name        string
	payload     string
}

type token struct {
	typ     int
	literal string
	number  int
}

func syntax_error() {
	log.Fatal(SYNTAX_ERROR)
}

type Lexer struct {
	reader *bytes.Buffer
}

func (lex *Lexer) init(r io.Reader) {
	bts, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println(err)
	}

	// 清除注释
	re := regexp.MustCompile("(?m:^#(.*)$)")
	bts = re.ReplaceAllLiteral(bts, nil)
	// 清除desc
	re = regexp.MustCompile("(?m:^desc:(.*)$)")
	bts = re.ReplaceAllLiteral(bts, nil)
	lex.reader = bytes.NewBuffer(bts)
}

func (lex *Lexer) next() (t *token) {
	defer func() {
		//	log.Println(t)
	}()
	var r rune
	var err error
	for {
		r, _, err = lex.reader.ReadRune()
		if err == io.EOF {
			return TOKEN_EOF
		} else if unicode.IsSpace(r) {
			continue
		}
		break
	}

	var runes []rune
	if unicode.IsLetter(r) {
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
		if tkid, ok := keywords[string(runes)]; ok {
			t.typ = tkid
		} else {
			t.typ = TK_STRING
			t.literal = string(runes)
		}
		return t
	} else if unicode.IsNumber(r) {
		for {
			runes = append(runes, r)
			r, _, err = lex.reader.ReadRune()
			if err == io.EOF {
				break
			} else if unicode.IsNumber(r) {
				continue
			} else {
				lex.reader.UnreadRune()
				break
			}
		}
		t := &token{}
		t.typ = TK_NUMBER
		n, _ := strconv.Atoi(string(runes))
		t.number = n
		return t
	} else if r == ':' {
		return TOKEN_COLON
	} else {
		syntax_error()
	}
	return TOKEN_EOF
}

//////////////////////////////////////////////////////////////
type Parser struct {
	exprs []api_expr
	lexer *Lexer
}

func (p *Parser) init(lex *Lexer) {
	p.lexer = lex
}

func (p *Parser) expr() bool {
	api := api_expr{}
	t := p.lexer.next()

	if t.typ == TK_EOF {
		return false
	}
	if t.typ == TK_TYPE {
		if p.lexer.next().typ == TK_COLON {
			if t := p.lexer.next(); t.typ == TK_NUMBER {
				api.packet_type = t.number
			} else {
				syntax_error()
			}
		}
	}
	if t := p.lexer.next(); t.typ == TK_NAME {
		if p.lexer.next().typ == TK_COLON {
			if t := p.lexer.next(); t.typ == TK_STRING {
				api.name = t.literal
			} else {
				syntax_error()
			}
		}
	}
	if t := p.lexer.next(); t.typ == TK_PAYLOAD {
		if p.lexer.next().typ == TK_COLON {
			if t := p.lexer.next(); t.typ == TK_STRING {
				api.payload = t.literal
			} else {
				syntax_error()
			}
		}
	}

	p.exprs = append(p.exprs, api)
	return true
}

func main() {
	lexer := Lexer{}
	lexer.init(os.Stdin)
	p := Parser{}
	p.init(&lexer)
	for p.expr() {
	}
	log.Println(p.exprs)
}
