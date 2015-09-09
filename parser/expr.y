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
  rule          *Rule
  pred          *Predicate
  preds         []pred
  act           *Action
  str           string
  f             float64
}
/*
Tokens include: number, &&, ->, identifier, >, <, >=, <=, ==, =
*/
%token LEX_ERROR
%token THEN
%token LE GE NE AND OR
%token '=' '<' '>'
%left '&' '|' 
%left '+' '-'
%left '*' '/' '%' 
/*%left <empty> '^' UMINUS*/

%token <str> identifier
%token <f> number 
%token<rule> rule
%token<pred> predicate
%token<preds> predicate_list
%token<act> action
%%
rule: 
  predicate_list THEN action
  {
    $$ = MakeRule($1, $3)
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
  }
    

value:
  identifier
| number

compare:
  '<'
  {
  }
| '>'
  {
  }
| LE
  {
  }
| GE
  {
  }
| NE
  {
  }
