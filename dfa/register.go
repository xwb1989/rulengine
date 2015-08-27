package dfa

type Register struct {
	states map[uint64][]*State
	size   int
}

//Hash state, find its equivalent, or selfister it
func (self *Register) GetOrPut(state *State) *State {
	key := state.Hash()
	states, ok := self.states[key]
	if ok {
		for _, other := range states {
			if state.Equals(other) {
				return other
			}
		}
		self.states[key] = append(states, state)
	} else {
		self.states[key] = []*State{state}
	}
	self.size += 1
	return state
}

func (self *Register) Size() int {
	return self.size
}

func (self *Register) Remove(state *State) bool {
	key := state.Hash()
	states, ok := self.states[key]
	if ok {
		for i, other := range states {
			if state.Equals(other) {
				self.states[key] = deleteFromSlice(states, i)
				self.size--
				if len(self.states[key]) == 0 {
					delete(self.states, key)
				}
				return true
			}
		}
	}
	return false
}

func MakeRegister() Register {
	return Register{states: make(map[uint64][]*State), size: 0}
}

func deleteFromSlice(states []*State, i int) []*State {
	states[i] = states[len(states)-1]
	states[len(states)-1] = nil
	states = states[:len(states)-1]
	return states
}
