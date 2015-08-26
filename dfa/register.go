package dfa

type Register struct {
	states map[uint64][]*State
	size   int
}

//Hash state, find its equivalent, or register it
func (reg *Register) GetOrPut(state *State) *State {
	key := state.Hash()
	states, ok := reg.states[key]
	if ok {
		for _, other := range states {
			if state.Equals(other) {
				return other
			}
		}
		reg.states[key] = append(states, state)
	} else {
		reg.states[key] = []*State{state}
	}
	state.Register()
	reg.size += 1
	return state
}

func (reg *Register) Size() int {
	return reg.size
}

func (reg *Register) Remove(state *State) bool {
	key := state.Hash()
	states, ok := reg.states[key]
	if ok {
		for i, other := range states {
			if state.Equals(other) {
				reg.states[key] = deleteFromSlice(states, i)
				reg.size--
				state.UnRegister()
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
