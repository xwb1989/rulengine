package dfa

// registry a customed map thus holding states that in the minimal DFA
type registry struct {
	states map[uint64][]*state
	size   int
}

func (r *registry) getOrPut(s *state) *state {
	key := s.Hash()
	states, ok := r.states[key]
	if ok {
		for _, other := range states {
			if s.Equals(other) {
				return other
			}
		}
		r.states[key] = append(states, s)
	} else {
		r.states[key] = []*state{s}
	}
	r.size++
	return s
}

//
func (r *registry) Size() int {
	return r.size
}

func (r *registry) Remove(state *state) bool {
	key := state.Hash()
	states, ok := r.states[key]
	if ok {
		for i, other := range states {
			if state.Equals(other) {
				r.states[key] = deleteFromSlice(states, i)
				r.size--
				if len(r.states[key]) == 0 {
					delete(r.states, key)
				}
				return true
			}
		}
	}
	return false
}

func makeRegistry() registry {
	return registry{states: make(map[uint64][]*state), size: 0}
}

func deleteFromSlice(states []*state, i int) []*state {
	states[i] = states[len(states)-1]
	states[len(states)-1] = nil
	states = states[:len(states)-1]
	return states
}
