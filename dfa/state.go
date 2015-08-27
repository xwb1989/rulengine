package dfa

import (
	. "github.com/xwb1989/quickdecider/parser"
	"math"
)

type State struct {
	IsRegistered bool
	Action       *Action //if not nil, it's an accept self
	edges        map[*Predicate]*State
	inCount      int
}

func MakeState(isRegistered bool, edges map[*Predicate]*State, action *Action) *State {
	return &State{IsRegistered: isRegistered, edges: edges, Action: action}
}
func InitState() *State {
	return MakeState(false, make(map[*Predicate]*State), nil)
}

func (self *State) FindNext(data interface{}) *State {
	var ret *State = nil
	for predicate, v := range self.edges {
		if predicate.IsTrue(data) {
			if ret != nil {
				panic("work on multiple predicate")
			}
			ret = v
		}
	}
	return ret
}

func (self *State) Next(pred *Predicate) *State {
	if self, ok := self.edges[pred]; ok {
		return self
	}
	return nil
}

func (self *State) Hit() {
	self.inCount++
}

func (self *State) Unhit() {
	self.inCount--
}

func (self *State) IsConfluent() bool {
	return self.inCount > 1
}

func (self *State) SetNext(pred *Predicate, next *State) {
	if oldNext, ok := self.edges[pred]; ok {
		oldNext.Unhit()
	}
	next.Hit()
	self.edges[pred] = next
}

func (self *State) Clone() *State {
	cloned := InitState()
	cloned.Action = self.Action
	for k, v := range self.edges {
		cloned.edges[k] = v
		v.Hit()
	}
	return cloned
}

func (self *State) Equals(other *State) bool {
	if self == other {
		return true
	}
	//need to make sure for same expressions we only create one Action instance
	if len(self.edges) != len(self.edges) || self.Action != other.Action {
		return false
	}
	for k, v := range self.edges {
		//need to make sure that for same expressions we only creates one Predicate instance
		if v != other.edges[k] {
			return false
		}
	}
	return true
}

func (self *State) Hash() uint64 {
	var sum uint64 = 0
	if self.Action != nil {
		sum = self.Action.Hash()
	}
	for k, v := range self.edges {
		sum += k.Hash() + v.Hash()
		sum %= math.MaxUint64
	}
	return sum
}
