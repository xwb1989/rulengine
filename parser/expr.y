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
  Rul           *Rule
  Pred          *Predicate
  Preds         []*Predicate
  Act           *Action
  Str           string
  Number        float64
  BooleanFn     BoolFn
  PredFunc      PredicateFunc
  ActFunc       ActionFunc
  Value         func(data interface{}) interface{} 
}
/*
Tokens include: number, &&, ->, identifier, >, <, >=, <=, ==, =
*/
%token <Str> identifier
%token <Number> number 
%token LEX_ERROR
%token THEN AND
%token <BooleanFn> LE GE NE EQ LT GT '<' '>'
%token '=' 
%left '&' '|' 
%left '+' '-'
%left '*' '/' '%' 
/*%left <empty> '^' UMINUS*/

%start any_rule
%type <Rul> rule
%type <Pred> predicate
%type <Preds> predicate_list
%type <Act> action
%type <BooleanFn> compare
%type <PredFunc> predicateFunc
%type <ActFunc> actionFunc
%type <Value> value

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
  predicateFunc 
  {
    $$ = MakePredicate(yytext, $1)
  }
predicateFunc:
  value compare value
  {
    $$ = func(data interface{}) bool { return $2($1, $3) }
  }

action:
  actionFunc
  { 
    $$ = MakeAction(yytext, $1)
  }

actionFunc:
  identifier '=' value
  {
    $$ = func(data interface{}) interface{} { setValue(data, $1, $3); return data }
  }

value:
  identifier 
  {
    $$ = func(data interface{}) interface{} { return getValue(data, $1) }
  }
| number
  {
    $$ = func(data interface{}) interface{} { return $1 }
  }

compare:
  '<' 
  { 
    $$ = ltFn 
  }
| '>' 
  { 
    $$ = gtFn 
  }
| LE  
  { 
    $$ = leFn 
  } 
| GE  
  { 
    $$ = geFn
  }
| NE
  {
    $$ = neFn 
  }
| EQ
  {
    $$ = eqFn
  }
