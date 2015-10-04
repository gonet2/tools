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
)

var (
	datatypes = map[string]bool{
		"integer": true,
		"string":  true,
		"bytes":   true,
	}
)

var (
	SYNTAX_ERROR = errors.New("syntax error")
)

type field struct {
	name  string
	typ   string
	array bool
}

type struct_info struct {
	name   string
	fields []field
}

type token struct {
	typ     int
	literal string
	r       rune
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
	lex.reader = bytes.NewBuffer(bts)
}

func (lex *Lexer) next() (t *token) {
	defer func() {
		log.Println(t)
	}()
	var r rune
	var err error
	for {
		r, _, err = lex.reader.ReadRune()
		if err == io.EOF {
			return nil
		} else if unicode.IsSpace(r) {
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
			continue
		} else {
			lex.reader.UnreadRune()
			return false
		}
	}

	return true
}

//////////////////////////////////////////////////////////////
type Parser struct {
	lexer *Lexer
}

func (p *Parser) init(lex *Lexer) {
	p.lexer = lex
}

func (p *Parser) expr() bool {
	if p.lexer.eof() {
		return false
	}
	if t := p.lexer.next(); t != nil {
		if t.typ != SYMBOL {
			syntax_error()
		}
	}

	if t := p.lexer.next(); t != nil {
		if t.typ != STRUCT_BEGIN {
			syntax_error()
		}
	}
	t := p.fields()
	if t != nil {
		if t.typ != STRUCT_END {
			syntax_error()
		}
	}

	return true
}

func (p *Parser) fields() *token {
	for {
		if t := p.lexer.next(); t != nil {
			if t.typ != SYMBOL {
				return t
			}
		}

		if t := p.lexer.next(); t != nil {
			if t.typ != DATA_TYPE {
				syntax_error()
			}
		}
	}
}

func syntax_error() {
	log.Fatal(SYNTAX_ERROR)
}

func check(t *token) {
	if t == nil {
		log.Fatal(SYNTAX_ERROR)
	}
}

func main() {
	lexer := Lexer{}
	lexer.init(os.Stdin)
	p := Parser{}
	p.init(&lexer)
	for p.expr() {
	}
}
