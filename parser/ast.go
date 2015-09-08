package parser

import (
	"hash/fnv"
)

type PredicateFunc func(interface{}) bool
type ActionFunc func(interface{}) interface{}

type exprBase struct {
	id   uint64
	expr string
}

type Node interface {
	Hash() uint64
	String() string
}

func (self *exprBase) String() string {
	return self.expr
}

func (self *exprBase) Hash() uint64 {
	return self.id
}

func (self *exprBase) Equals(other Node) bool {
	if self == other {
		return true
	} else if self == nil || other == nil {
		return false
	}
	return self.Hash() == other.Hash() && self.String() == other.String()
}

func MakeExprBase(expr string) exprBase {
	h := fnv.New64a()
	h.Write([]byte(expr))
	return exprBase{id: h.Sum64(), expr: expr}
}

type Predicate struct {
	exprBase
	functor PredicateFunc
}

func MakePredicate(expr string, functor PredicateFunc) *Predicate {
	return &Predicate{exprBase: MakeExprBase(expr), functor: functor}
}

func (self *Predicate) IsTrue(data interface{}) bool {
	return self.functor(data)
}

type Action struct {
	exprBase
	functor ActionFunc
}

func MakeAction(expr string, functor ActionFunc) *Action {
	return &Action{exprBase: MakeExprBase(expr), functor: functor}
}

func (self *Action) Apply(val interface{}) interface{} {
	return self.functor(val)
}

type Rule struct {
	Predicates []*Predicate
	Action     *Action
}

func MakeRule(preds []*Predicate, act *Action) *Rule {
	return &Rule{Predicates: preds, Action: act}
}

func lt(a int, b int) bool {
	return a < b
}

func le(a int, b int) bool {
	return a <= b
}

func eq(a interface{}, b interface{}) bool {
	return a == b
}

/**
we have lazy evaluation...
*/
type BoolFn func() bool

func and(a BoolFn, b BoolFn) bool {
	return a() && b()
}

func or(a BoolFn, b BoolFn) bool {
	return a() || b()
}
