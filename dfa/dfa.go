package dfa

import (
	"fmt"
	"github.com/oleiade/lane"
	. "github.com/xwb1989/rulengine/parser"
	"log"
)

type DFA struct {
	start *State
}

func MakeDFA(rules []*Rule) DFA {
	start := createMinimalDFA(rules)
	return DFA{start: start}
}

func (self *DFA) GetAction(data interface{}) *Action {
	prev := self.start
	found := []*Action{}
	i := 0
	for curr := prev.FindNext(data); curr != nil; curr = curr.FindNext(data) {
		if prev.FindNext(data) != curr {
			fmt.Println("prev.FindNext(val)=", prev.FindNext(data), " curr=", curr)
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
	} else {
		return nil
	}
}

func (self *DFA) Size() int {
	stack := lane.NewStack()
	visited := map[*State]bool{}
	stack.Push(self.start)
	visited[self.start] = true
	cnt := 0
	for !stack.Empty() {
		curr := stack.Pop().(*State)
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

func (self *DFA) Debug() {
	stack := lane.NewStack()
	visited := map[*State]bool{}
	stack.Push(self.start)
	visited[self.start] = true
	cnt := 0
	log.Println("DFA Debug:")
	for !stack.Empty() {
		curr := stack.Pop().(*State)
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
func createMinimalDFA(rules []*Rule) *State {
	start := InitState()
	start.IsRegistered = false
	reg := MakeRegister()
	reg.GetOrPut(start)
	//this must be called before state is modified
	unRegister := func(state *State) {
		if state != start && state.IsRegistered {
			if !reg.Remove(state) {
				log.Panic("fail to unregister state:", state)
			} else {
				state.IsRegistered = false
			}
		}
	}
	for _, rule := range rules {
		path := make([]*State, len(rule.Predicates)+1)
		path[0] = start
		prev := start
		for i, pred := range rule.Predicates {
			curr := prev.Next(pred)
			if curr == nil {
				curr = InitState()
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
			in_path := path[i+1]
			in_reg := reg.GetOrPut(in_path)
			if in_path != in_reg { //fail to register
				prev_in_path := path[i]
				unRegister(prev_in_path)
				prev_in_path.SetNext(rule.Predicates[i], in_reg)
			} else {
				in_path.IsRegistered = true
			}
		}
	}
	return start
}
