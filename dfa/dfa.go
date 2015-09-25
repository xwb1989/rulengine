package dfa

import (
	"fmt"
	"log"

	"github.com/oleiade/lane"
)

// DFA the DFA
type DFA struct {
	start *state
}

// MakeDFA create a DFA by given a set of rules
func MakeDFA(rules []*Rule) DFA {
	start := createMinimalDFA(rules)
	return DFA{start: start}
}

// GetAction by given data, traverse the DFA rule space and return an Action
func (d *DFA) GetAction(data interface{}) *Action {
	prev := d.start
	found := []*Action{}
	i := 0
	for curr := prev.findNext(data); curr != nil; curr = curr.findNext(data) {
		if prev.findNext(data) != curr {
			fmt.Println("prev.findNext(val)=", prev.findNext(data), " curr=", curr)
			panic("wrong!!!")
		}
		if curr.Action != nil {
			found = append(found, curr.Action)
		}
		prev = curr
		i++
	}
	if len(found) != 0 {
		return found[len(found)-1]
	}
	return nil
}

// Size size of this DFA
func (d *DFA) Size() int {
	stack := lane.NewStack()
	visited := map[*state]bool{}
	stack.Push(d.start)
	visited[d.start] = true
	cnt := 0
	for !stack.Empty() {
		curr := stack.Pop().(*state)
		cnt++
		for _, v := range curr.edges {
			if _, ok := visited[v]; !ok {
				visited[v] = true
				stack.Push(v)
			}
		}
	}
	return cnt
}

// Debug output debug information about this DFA
func (d *DFA) Debug() {
	stack := lane.NewStack()
	visited := map[*state]bool{}
	stack.Push(d.start)
	visited[d.start] = true
	cnt := 0
	log.Println("DFA Debug:")
	for !stack.Empty() {
		curr := stack.Pop().(*state)
		log.Printf("\n%p:\n%v\n\n", curr, curr)
		cnt++
		for _, v := range curr.edges {
			if _, ok := visited[v]; !ok {
				visited[v] = true
				stack.Push(v)
			}
		}
	}
	log.Println("DFA Debug: size ", cnt)
}

//Implements algorithm in paper: Incremental construction of minimal acyclic finite-state automata
//Predicates in each rule must be sorted in a consistent order
func createMinimalDFA(rules []*Rule) *state {
	start := initState()
	start.IsRegistered = false
	reg := makeRegistry()
	reg.getOrPut(start)
	//this must be called before state is modified
	unRegister := func(state *state) {
		if state != start && state.IsRegistered {
			if !reg.Remove(state) {
				log.Panic("fail to unregister state:", state)
			} else {
				state.IsRegistered = false
			}
		}
	}
	for _, rule := range rules {
		path := make([]*state, len(rule.Predicates)+1)
		path[0] = start
		prev := start
		for i, pred := range rule.Predicates {
			curr := prev.Next(pred)
			if curr == nil {
				curr = initState()
				unRegister(prev)
				prev.SetNext(pred, curr)
			} else if curr.IsConfluent() {
				curr = curr.Clone()
				unRegister(prev)
				prev.SetNext(pred, curr)
			}
			path[i+1] = curr
			prev = curr
		}
		//we can for sure move on
		if prev.Action != nil && prev.Action == rule.Action {
			continue
		} else if prev.Action != nil && prev.Action != rule.Action {
			//TODO: error! Same predications but different action!
			log.Panic("error! same predicates but different action!\n", prev.Action, "\n", rule.Action)
		}

		unRegister(prev) //because we are gonna modify its Action
		prev.Action = rule.Action

		//traverse back the path
		for i := len(rule.Predicates) - 1; i >= 0; i-- {
			inPath := path[i+1]
			inReg := reg.getOrPut(inPath)
			if inPath != inReg { //fail to register
				prevInPath := path[i]
				unRegister(prevInPath)
				prevInPath.SetNext(rule.Predicates[i], inReg)
			} else {
				inPath.IsRegistered = true
			}
		}
	}
	return start
}
