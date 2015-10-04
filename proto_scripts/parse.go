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
	"strconv"
	"unicode"
)

const (
	PACKET_TYPE = iota
	NAME
	PAYLOAD
	DESC
	LITERAL
	NUM
)

var (
	SYNTAX_ERROR = errors.New("syntax error")
)

type token struct {
	typ     int
	literal string
	num     int
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

func (lex *Lexer) next_rune() rune {
	var r rune
	var err error
	for {
		r, _, err = lex.reader.ReadRune()
		if err != nil {
			return ' '
		}

		if unicode.IsSpace(r) {
			continue
		} else {
			break
		}
	}
	return r
}

func (lex *Lexer) next() *token {
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
		if r != ':' && !unicode.IsSpace(r) {
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
		if num, err := strconv.Atoi(t.literal); err == nil {
			t.num = num
			t.typ = NUM
		} else {
			t.typ = LITERAL
		}
	}

	return t
}

type Parser struct {
	lexer *Lexer
}

func (p *Parser) init(lex *Lexer) {
	p.lexer = lex
}

func (p *Parser) match(r rune) {
	n := p.lexer.next_rune()
	if r != n {
		log.Fatal(SYNTAX_ERROR)
	}
}

func (p *Parser) packet_type() {
	t := p.lexer.next()
	if t.typ != PACKET_TYPE {
		log.Fatal(SYNTAX_ERROR)
	}
}

func (p *Parser) name() {
	t := p.lexer.next()
	if t.typ != NAME {
		log.Fatal(SYNTAX_ERROR)
	}
}

func (p *Parser) payload() {
	t := p.lexer.next()
	if t.typ != PAYLOAD {
		log.Fatal(SYNTAX_ERROR)
	}
}

func (p *Parser) desc() {
	t := p.lexer.next()
	if t.typ != DESC {
		log.Fatal(SYNTAX_ERROR)
	}
}

func (p *Parser) literal() string {
	t := p.lexer.next()
	return t.literal
}

func (p *Parser) number() int {
	t := p.lexer.next()
	if t.typ != LITERAL {
		log.Fatal(SYNTAX_ERROR)
	}
	return t.num
}

func (p *Parser) expr() {
	p.packet_type()
	p.match(':')
	println(p.literal())
	p.name()
	p.match(':')
	println(p.literal())

	p.payload()
	p.match(':')
	println(p.literal())

	p.desc()
	p.match(':')
	println(p.literal())
}

func main() {
	lexer := Lexer{}
	lexer.init(os.Stdin)
	p := Parser{}
	p.init(&lexer)
	p.expr()
}
