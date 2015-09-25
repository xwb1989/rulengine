package dfa

import (
	"math"
)

type state struct {
	IsRegistered bool
	Action       *Action //if not nil, it's an accept s
	edges        map[*Predicate]*state
	inCount      int
}

func makestate(isRegistered bool, edges map[*Predicate]*state, action *Action) *state {
	return &state{IsRegistered: isRegistered, edges: edges, Action: action}
}
func initState() *state {
	return makestate(false, make(map[*Predicate]*state), nil)
}

func (s *state) findNext(data interface{}) *state {
	var ret *state
	for predicate, v := range s.edges {
		if predicate.IsTrue(data) {
			if ret != nil {
				panic("work on multiple predicate")
			}
			ret = v
		}
	}
	return ret
}

func (s *state) Next(pred *Predicate) *state {
	if s, ok := s.edges[pred]; ok {
		return s
	}
	return nil
}

func (s *state) Hit() {
	s.inCount++
}

func (s *state) Unhit() {
	s.inCount--
}

func (s *state) IsConfluent() bool {
	return s.inCount > 1
}

func (s *state) SetNext(pred *Predicate, next *state) {
	if oldNext, ok := s.edges[pred]; ok {
		oldNext.Unhit()
	}
	next.Hit()
	s.edges[pred] = next
}

func (s *state) Clone() *state {
	cloned := initState()
	cloned.Action = s.Action
	for k, v := range s.edges {
		cloned.edges[k] = v
		v.Hit()
	}
	return cloned
}

func (s *state) Equals(other *state) bool {
	if s == other {
		return true
	}
	//need to make sure for same expressions we only create one Action instance
	if len(s.edges) != len(s.edges) || s.Action != other.Action {
		return false
	}
	for k, v := range s.edges {
		//need to make sure that for same expressions we only creates one Predicate instance
		if v != other.edges[k] {
			return false
		}
	}
	return true
}

func (s *state) Hash() uint64 {
	var sum uint64
	if s.Action != nil {
		sum = s.Action.Hash()
	}
	for k, v := range s.edges {
		sum += k.Hash() + v.Hash()
		sum %= math.MaxUint64
	}
	return sum
}
