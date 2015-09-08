%{
package parser
import "bytes"

func SetParseTree(yylex interface{}, stmt Statement) {
  yylex.(*Lexer).ParseTree = stmt
}

func SetAllowComments(yylex interface{}, allow bool) {
  yylex.(*Lexer).AllowComments = allow
}

func ForceEOF(yylex interface{}) {
  yylex.(*Lexer).ForceEOF = true
}


%}

%union {
  empty         struct{}
  rule          Rule
  pred          Predicate
  act           Action
  str           string
  number        float64
}
/*
Tokens include: number, &&, ->, identifier, >, <, >=, <=, ==, =, +, -, *, /, %, ^
*/
%token LEX_ERROR
%token <empty> THEN
%token <empty> LE GE NE AND OR
%token <empty> '(' '=' '<' '>' '~'
%left <empty> '&' '|' 
%left <empty> '+' '-'
%left <empty> '*' '/' '%' 
%left <empty> '^' UMINUS
%token <str> identifier
%%
/*rule: */
    /*predicate_list THEN action_list*/
    /*{*/
    /*}*/

/*predicate_list:*/
      /*predicate */
    /*| predicate_list AND predicate*/

/*predicate:*/
      /*value compare value */
    /*| '(' predicate OR predicate ')'*/

/*action_list:*/
      /*action*/
    /*| action_list ',' action*/

/*action:*/
    /*| identifier '=' value*/

/*value:*/
    /*identifier*/

/*compare:*/
  /*'='*/
  /*{*/
  /*}*/
/*| '<'*/
  /*{*/
  /*}*/
/*| '>'*/
  /*{*/
  /*}*/
/*| LE*/
  /*{*/
  /*}*/
/*| GE*/
  /*{*/
  /*}*/
/*| NE*/
  /*{*/
  /*}*/

/*%union {*/
    /*i int*/
/*}*/

/*%token<i> INT*/
/*%token<i> expr*/

/*%%*/
/*expr: term { $$ = $1 }*/
 /*| expr '+' term { $$ = $1 + $3 }*/
 /*| expr '-' term { $$ = $1 - $3 }*/
 /*| '-' expr { $$ = -$1 }*/

/*term: factor { $$ = $1 }*/
    /*| term '*' factor {$$ = $1 * $3}*/
    /*| term '/' factor {$$ = $1 / $3}*/
    /*| term '%' factor {$$ = $1 % $3}*/

/*factor: INT {$$ = $1}*/



// Copyright 2011 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// based off of Appendix A from http://dinosaur.compilertools.net/yacc/

%{

package parser

import (
    "bufio"
        "fmt"
            "os"
                "unicode"
                )

var regs = make([]int, 26)
                var base int

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
    val int
    }

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <val> expr number

// same for terminals
%token <val> DIGIT LETTER

%left '|'
%left '&'
%left '+'  '-'
%left '*'  '/'  '%'
%left UMINUS      /*  supplies  precedence  for  unary  minus  */

%%

list    : /* empty */
        | list stat '\n'
                ;

stat    :    expr
        {
                            fmt.Printf( "%d\n", $1 );
                                    }
                                        |    LETTER '=' expr
                                                {
                                                            regs[$1]  =  $3
                                                                    }
                                                                        ;

expr    :    '(' expr ')'
        { $$  =  $2 }
                    |    expr '+' expr
                            { $$  =  $1 + $3 }
                                |    expr '-' expr
                                        { $$  =  $1 - $3 }
                                            |    expr '*' expr
                                                    { $$  =  $1 * $3 }
                                                        |    expr '/' expr
                                                                { $$  =  $1 / $3 }
                                                                    |    expr '%' expr
                                                                            { $$  =  $1 % $3 }
                                                                                |    expr '&' expr
                                                                                        { $$  =  $1 & $3 }
                                                                                            |    expr '|' expr
                                                                                                    { $$  =  $1 | $3 }
                                                                                                        |    '-'  expr        %prec  UMINUS
                                                                                                                { $$  = -$2  }
                                                                                                                    |    LETTER
                                                                                                                            { $$  = regs[$1] }
                                                                                                                                |    number
                                                                                                                                    ;

number  :    DIGIT
        {
                            $$ = $1;
                                        if $1==0 {
                                                        base = 8
                                                                    } else {
                                                                                    base = 10
                                                                                                }
                                                                                                        }
                                                                                                            |    number DIGIT
                                                                                                                    { $$ = base * $1 + $2 }
                                                                                                                        ;

%%      /*  start  of  programs  */

type CalcLex struct {
    s string
        pos int
        }


func (l *CalcLex) Lex(lval *CalcSymType) int {
    var c rune = ' '
        for c == ' ' {
                if l.pos == len(l.s) {
                            return 0
                                    }
                                            c = rune(l.s[l.pos])
                                                    l.pos += 1
                                                        }

if unicode.IsDigit(c) {
                                                                    lval.val = int(c) - '0'
                                                                            return DIGIT
                                                                                } else if unicode.IsLower(c) {
                                                                                        lval.val = int(c) - 'a'
                                                                                                return LETTER
                                                                                                    }
                                                                                                        return int(c)
                                                                                                        }

func (l *CalcLex) Error(s string) {
    fmt.Printf("syntax error: %s\n", s)
    }

func main() {
    fi := bufio.NewReader(os.NewFile(0, "stdin"))

for {
                var eqn string
                        var ok bool

fmt.Printf("equation: ")
                                        if eqn, ok = readline(fi); ok {
                                                    CalcParse(&CalcLex{s: eqn})
                                                            } else {
                                                                        break
                                                                                }
                                                                                    }
                                                                                    }

func readline(fi *bufio.Reader) (string, bool) {
    s, err := fi.ReadString('\n')
        if err != nil {
                return "", false
                    }
                        return s, true
                        }

