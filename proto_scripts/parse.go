package main

import (
	"bufio"
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
}

type Lexer struct {
	reader *bufio.Reader
}

func (lex *Lexer) init(r io.Reader) {
	bts, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println(err)
	}

	re := regexp.MustCompile("(?m:^#(.*)$)")
	bts = re.ReplaceAllLiteral(bts, []byte("\n"))
	lex.reader = bufio.NewReader(bytes.NewBuffer(bts))
}

func (lex *Lexer) keyword() *token {
	var r rune
	var err error
	for {
		r, _, err = lex.reader.ReadRune()
		if err != nil {
			return nil
		}
		if unicode.IsSpace(r) {
			continue
		} else {
			break
		}
	}

	var runes []rune
	for {
		runes = append(runes, r)
		r, _, err = lex.reader.ReadRune()
		if r != ':' {
			continue
		} else {
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
	default:
		t.typ = LITERAL
	}
	return t
}

func (lex *Lexer) r() *token {
	t := &token{}
	r, _, err := lex.reader.ReadRune()
	if err != nil {
		return nil
	}
	t.literal = string(r)
	return t
}

func (lex *Lexer) str() *token {
	var r rune
	var err error
	var runes []rune
	t := &token{}
	for {
		r, _, err = lex.reader.ReadRune()
		if err != nil {
			return nil
		}
		if r != '\r' && r != '\n' {
			runes = append(runes, r)
			continue
		} else {
			lex.reader.UnreadRune()
			break
		}
	}
	t.literal = string(runes)
	return t
}

//////////////////////////////////////////////////////////////
type Parser struct {
	exprs []api_expr
	lexer *Lexer
}

func (p *Parser) init(lex *Lexer) {
	p.lexer = lex
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

func (p *Parser) expr() bool {
	api := api_expr{}
	p.packet_type()
	p.lexer.r()
	api.packet_type = p.lexer.str().literal
	p.name()
	p.lexer.r()
	api.name = p.lexer.str().literal
	p.payload()
	p.lexer.r()
	api.payload = p.lexer.str().literal
	p.desc()
	p.lexer.r()
	api.desc = p.lexer.str().literal

	p.exprs = append(p.exprs, api)
	return false
}

func check(t *token) {
	if t == nil {
		log.Fatal(SYNTAX_ERROR)
	}

	log.Println(t)
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
