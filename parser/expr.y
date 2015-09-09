%{
package parser

import "bytes"

func SetParseTree(yylex interface{}, rule *Rule) {
  yylex.(*Lexer).ParseTree = rule
}

func ForceEOF(yylex interface{}) {
  yylex.(*Lexer).ForceEOF = true
}


%}

%union {
  empty         struct{}
  Rul           *Rule
  Pred          *Predicate
  Preds         []*Predicate
  Act           *Action
  Str           string
  Number        float64
}
/*
Tokens include: number, &&, ->, identifier, >, <, >=, <=, ==, =
*/
%token <Str> identifier
%token <Number> number 
%token <empty> LEX_ERROR
%token <empty> THEN
%token <empty> LE GE NE AND OR
%token <empty> '=' '<' '>'
%left <empty> '&' '|' 
%left <empty> '+' '-'
%left <empty> '*' '/' '%' 
/*%left <empty> '^' UMINUS*/

%start any_rule
%type <Rul> rule
%type <Pred> predicate
%type <Preds> predicate_list
%type <Act> action


%%

any_rule:
  rule
  {
    SetParseTree(yylex, $1)
  }

rule: 
  predicate_list THEN action
  {
    MakeAction($1, $3)
  }

predicate_list: 
  predicate
  {
    $$ = Predicates{$1}
  }
| predicate_list AND predicate
  {
    $$ = append($$, $3)
  }

predicate:
  value compare value
  {
    $$ = MakePredicate()
  }

action:
  identifier '=' value { $$ = MakeAction()}
    

value:
  identifier
| number

compare:
  '<'
| '>'
| LE
| GE
| NE

