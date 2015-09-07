package parser

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	p := Parser{Buffer: "1+2"}
	p.Init()
	p.Execute()
	fmt.Println(p.Sum)
}
