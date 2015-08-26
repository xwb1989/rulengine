package dfa

import (
	. "github.com/xwb1989/quickdecider/parser"
)

type State struct {
	isRegistered bool
	edges        map[*Predicate]*State
	action       *Action //if not nil, it's an accept state
	inCount      int
}

func MakeState(isRegistered bool, edges map[*Predicate]*State, action *Action) *State {
	return &State{isRegistered: isRegistered, edges: edges, action: action}
}
func InitState() *State {
	return MakeState(false, make(map[*Predicate]*State), nil)
}

func (state *State) SetAction(action *Action) {
	state.action = action
}

func (state *State) IsRegistered() bool {
	return state.isRegistered
}

func (state *State) Register() {
	state.isRegistered = true
}

func (state *State) UnRegister() {
	state.isRegistered = false
}

func (state *State) Next(val interface{}) *State {
	for predicate, v := range state.edges {
		if predicate.IsTrue(val) {
			return v
		}
	}
	return nil
}

func (state *State) Hit() {
	state.inCount++
}

func (state *State) Unhit() {
	state.inCount--
}

func (state *State) IsConfluent() bool {
	return state.inCount > 1
}

func (state *State) SetNext(pred *Predicate, next *State) {
	if oldNext, ok := state.edges[pred]; ok {
		oldNext.Unhit()
	}
	next.Hit()
	state.edges[pred] = next
}

func (state *State) Equals(other *State) bool {
	if state == other {
		return true
	}
	//need to make sure for same expressions we only creates one Action instance
	if len(state.edges) != len(state.edges) || state.action != other.action {
		return false
	}
	for k, v := range state.edges {
		//need to make sure that for same expressions we only creates one Predicate instance
		if v != other.edges[k] {
			return false
		}
	}
	return true
}

func (state *State) Hash() uint64 {
	var sum uint64 = 0
	if state.action != nil {
		sum = state.action.Hash()
	}
	var i uint64 = 0
	for k, v := range state.edges {
		sum += (k.Hash()*7 + (v.Hash()>>2)*101) * (11 + 2*i)
		i++
	}
	return sum
}
