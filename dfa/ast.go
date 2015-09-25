package dfa

import (
	"hash/fnv"
)

// PredicateFn a function that represents the logic of a Predicate
type PredicateFn func(interface{}) bool

// ActionFn a function that represents the logic of an Action
type ActionFn func(interface{}) interface{}

type exprBase struct {
	id   uint64
	expr string
}

// Node a node in ast
type Node interface {
	Hash() uint64
	String() string
}

// String() return string representation of this expression
func (expr *exprBase) String() string {
	return expr.expr
}

func (expr *exprBase) Hash() uint64 {
	return expr.id
}

func (expr *exprBase) Equals(other Node) bool {
	if expr == other {
		return true
	} else if expr == nil || other == nil {
		return false
	}
	return expr.Hash() == other.Hash() && expr.String() == other.String()
}

func makeExprBase(expr string) exprBase {
	h := fnv.New64a()
	h.Write([]byte(expr))
	return exprBase{id: h.Sum64(), expr: expr}
}

// Predicate a predicate
type Predicate struct {
	exprBase
	functor PredicateFn
}

// MakePredicate create a predicate
func MakePredicate(expr string, functor PredicateFn) *Predicate {
	return &Predicate{exprBase: makeExprBase(expr), functor: functor}
}

// IsTrue given data, return whether the predicate evaluates to true
func (expr *Predicate) IsTrue(data interface{}) bool {
	return expr.functor(data)
}

// Action an action
type Action struct {
	exprBase
	functor ActionFn
}

// MakeAction create an Action
func MakeAction(expr string, functor ActionFn) *Action {
	return &Action{exprBase: makeExprBase(expr), functor: functor}
}

// Apply apply the action
func (expr *Action) Apply(val interface{}) interface{} {
	return expr.functor(val)
}

// Rule a rule
type Rule struct {
	Predicates []*Predicate
	Action     *Action
}

// MakeRule create a rule
func MakeRule(preds []*Predicate, act *Action) *Rule {
	return &Rule{Predicates: preds, Action: act}
}

/**
some functions
*/

// BoolFn functions that return bool
type BoolFn func(interface{}, interface{}) bool

// LtFn compare a and b
func LtFn(a interface{}, b interface{}) bool {
	return a.(float64) < b.(float64)
}

// GtFn compare a and b
func GtFn(a interface{}, b interface{}) bool {
	return a.(float64) > b.(float64)
}

// LeFn compare a and b
func LeFn(a interface{}, b interface{}) bool {
	return a.(float64) <= b.(float64)
}

// GeFn compare a and b
func GeFn(a interface{}, b interface{}) bool {
	return a.(float64) >= b.(float64)
}

// EqFn compare a and b
func EqFn(a interface{}, b interface{}) bool {
	return a == b
}

// NeFn compare a and b
func NeFn(a interface{}, b interface{}) bool {
	return a != b
}

/**
we have lazy evaluation...
*/
