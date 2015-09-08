package parser

import (
	"strings"
)

const EOFCHAR = 0x100

type StateFn func(*Lexer) StateFn
type TokenType uint8
type Token struct {
	Type  TokenType
	Value string
}

const (
	TokenNumber     TokenType = iota //numbers
	TokenIdentifier                  //identifiers
	TokenOperator                    //operators, like >, <
	TokenString                      //quoted strings
)

// Tokenizer is the struct used to generate SQL
// tokens for the parser.
type Lexer struct {
	InStream  *strings.Reader
	lastChar  uint16
	Name      string
	Position  int
	Start     int
	Width     int
	ParseTree Rule
}

func stateStart(lex *Lexer) StateFn {
	return nil
}

func (self *Lexer) skipBlank() {
	ch := self.lastChar
	for ch == ' ' || ch == '\n' || ch == '\r' || ch == '\t' {
		self.next()
		ch = self.lastChar
	}
}

func (self *Lexer) next() {
	if ch, err := self.InStream.ReadByte(); err != nil {
		// Only EOF is possible.
		self.lastChar = EOFCHAR
	} else {
		self.lastChar = uint16(ch)
	}
	self.Position++
}

func NewStringLexer(expr string) *Lexer {
	return &Lexer{InStream: strings.NewReader(expr)}
}

func (self *Lexer) Run() {
	for state := stateStart; state != nil; {
		state(self)
	}
}

// Lex returns the next token form the Tokenizer.
// This function is used by go yacc.
func (self *Lexer) Lex(lval *yySymType) int {
	return 0
}

// Error is called by go yacc if there's a parsing error.
func (self *Lexer) Error(err string) {
}
