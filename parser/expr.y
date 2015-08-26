%{
package parser
%}
%union {
    empty       struct{}
    boolExpr    BoolExpr
    valExpr     ValExpr
    bytes       []byte
}

%token <bytes> ID STRING NUMBER VALUE_ARG LIST_ARG COMMENT
%token <empty> '(' ')'
