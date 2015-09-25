package parser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type itemType uint8

const (
	itemError itemType = iota //numbers
	itemEOF
	itemIdentifier
	itemKeyword
	itemThen
	itemAnd
)

type item struct {
	typ itemType
	pos int
	val string
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

const eof = -1

type stateFn func(*Lexer) stateFn

type Lexer struct {
	name    string
	input   string
	state   stateFn
	pos     int
	start   int
	width   int
	lastPos int
	items   chan item
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

func (l *Lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

func (l *Lexer) drain() {
	for range l.items {
	}
}

func lex(name, input string) *Lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

func (l *Lexer) run() {
	for state := stateStart; state != nil; {
		state = state(l)
	}
	close(l.items)
}

// States...
func stateStart(l *Lexer) stateFn {
	return nil
}

// Lex returns the next token form the Tokenizer.
// This function is called by go yacc.
func (self *Lexer) Lex(lval *yySymType) int {
	return 0
}

// Error is called by go yacc if there's a parsing error.
func (self *Lexer) Error(err string) {
}

// Helper functions
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
