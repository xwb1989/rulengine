package peg

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	expr := "(1 - -3) / 3 + 2 *(3 + -4) +3 % 2 ^ 2"
	calc := &Calculator{Buffer: expr}
	calc.Init()
	calc.Expression.Init(expr)
	if err := calc.Parse(); err != nil {
		fmt.Errorf("Error: unable to parse expression %v", expr)
	}
	calc.Execute()
	fmt.Printf("= %v\n", calc.Evaluate())
}

func BenchmarkParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		expr := "(1 - -3) / 3 + 2 *(3 + -4) +3 % 2 ^ 2"
		calc := &Calculator{Buffer: expr}
		calc.Init()
		calc.Expression.Init(expr)
		if err := calc.Parse(); err != nil {
			fmt.Errorf("Error: unable to parse expression %v", expr)
		}
		calc.Execute()
		calc.Evaluate()
	}
}
