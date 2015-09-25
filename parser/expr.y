%{
package parser

import "github.com/xwb1989/rulengine/dfa"

func SetParseTree(yylex interface{}, rule *dfa.Rule) {
  yylex.(*lexer).ParseTree = rule
}

func ForceEOF(yylex interface{}) {
  yylex.(*lexer).ForceEOF = true
}


%}

%union {
  Rul                   *dfa.Rule
  Predicate             *dfa.Predicate
  Predicates            []*dfa.Predicate
  Action                *dfa.Action
  Str                   string
  Number                float64
  BoolFn                dfa.BoolFn
  PredFn                dfa.PredicateFn
  ActFn                 dfa.ActionFn
  Value                 func(data interface{}) interface{} 
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
%type <Predicate> predicate
%type <Predicates> predicate_list
%type <Action> action
%type <BoolFn> compare
%type <PredFn> predicateFunc
%type <ActFn> actionFunc
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
    dfa.MakeAction($1, $3)
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
    $$ = dfa.MakePredicate(yytext, $1)
  }
predicateFunc:
  value compare value
  {
    $$ = func(data interface{}) bool { return $2($1, $3) }
  }

action:
  actionFunc
  { 
    $$ = dfa.MakeAction(yytext, $1)
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
