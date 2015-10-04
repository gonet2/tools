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
	PACKET_TYPE = iota
	NAME
	PAYLOAD
	DESC
	LITERAL
)

var (
	SYNTAX_ERROR = errors.New("syntax error")
)

type api_expr struct {
	packet_type string
	name        string
	payload     string
	desc        string
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

	re := regexp.MustCompile("(?m:^#(.*)$)")
	bts = re.ReplaceAllLiteral(bts, nil)
	lex.reader = bytes.NewBuffer(bts)
}

func (lex *Lexer) keyword() *token {
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

	var runes []rune
	for {
		runes = append(runes, r)
		r, _, err = lex.reader.ReadRune()
		if err == io.EOF {
			break
		} else if r == ':' {
			lex.reader.UnreadRune()
			break
		}
	}
	t := &token{}
	t.literal = string(runes)
	switch t.literal {
	case "packet_type":
		t.typ = PACKET_TYPE
	case "name":
		t.typ = NAME
	case "payload":
		t.typ = PAYLOAD
	case "desc":
		t.typ = DESC
	}
	return t
}

func (lex *Lexer) r() *token {
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

	t := &token{}
	t.r = r
	return t
}

func (lex *Lexer) str() *token {
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

	var runes []rune
	for {
		runes = append(runes, r)
		r, _, err = lex.reader.ReadRune()
		if err == io.EOF {
			break
		} else if r == '\r' || r == '\n' {
			break
		}
	}

	t := &token{}
	t.literal = string(runes)
	return t
}

func (lex *Lexer) eof() bool {
	for {
		r, _, err := lex.reader.ReadRune()
		if err == io.EOF {
			log.Println(r, err)
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
	exprs []api_expr
	lexer *Lexer
}

func (p *Parser) init(lex *Lexer) {
	p.lexer = lex
}

func (p *Parser) match(r rune) {
	t := p.lexer.r()
	check(t)
	if t.r != r {
		log.Fatal(SYNTAX_ERROR)
	}
}

func (p *Parser) packet_type() {
	t := p.lexer.keyword()
	check(t)
	if t.typ != PACKET_TYPE {
		log.Fatal(SYNTAX_ERROR)
	}
}

func (p *Parser) name() {
	t := p.lexer.keyword()
	check(t)
	if t.typ != NAME {
		log.Fatal(SYNTAX_ERROR)
	}
}
func (p *Parser) payload() {
	t := p.lexer.keyword()
	check(t)
	if t.typ != PAYLOAD {
		log.Fatal(SYNTAX_ERROR)
	}
}
func (p *Parser) desc() {
	t := p.lexer.keyword()
	check(t)
	if t.typ != DESC {
		log.Fatal(SYNTAX_ERROR)
	}
}

func (p *Parser) literal() string {
	t := p.lexer.str()
	check(t)
	return t.literal
}

func (p *Parser) expr() bool {
	api := api_expr{}
	if p.lexer.eof() {
		return false
	}
	p.packet_type()
	p.match(':')
	api.packet_type = p.literal()
	p.name()
	p.match(':')
	api.name = p.literal()
	p.payload()
	p.match(':')
	api.payload = p.literal()
	p.desc()
	p.match(':')
	api.desc = p.literal()

	p.exprs = append(p.exprs, api)
	return true
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
	log.Println(p.exprs)
}
