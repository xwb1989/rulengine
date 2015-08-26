package quickdecider

type Decider struct {
}

func MakeDecider(conf map[string]interface{}) (Decider, error) {
	return Decider{}, nil
}

func (decider *Decider) GetAction(data map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
