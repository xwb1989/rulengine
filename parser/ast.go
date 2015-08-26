package parser

type PredicateFunc func(interface{}) bool
type ActionFunc func(interface{}) interface{}

type exprBase struct {
	id   uint64
	expr string
}

func (expr *exprBase) Hash() uint64 {
	return expr.id
}
func (expr *exprBase) String() string {
	return expr.expr
}

type Predicate struct {
	exprBase
	functor PredicateFunc
}

func MakePredicate(expr string) *Predicate {
	//hash expr to get id, parse expr to get functor
	return &Predicate{exprBase: exprBase{id: 0, expr: expr}, functor: nil}
}

func (pred *Predicate) IsTrue(val interface{}) bool {
	return pred.functor(val)
}

type Action struct {
	exprBase
	functor ActionFunc
}

func MakeAction(expr string) *Action {
	//hash expr to get id, parse expr to get functor
	return &Action{exprBase: exprBase{id: 0, expr: expr}, functor: nil}
}

func (pred *Action) Apply(val interface{}) interface{} {
	return pred.functor(val)
}
