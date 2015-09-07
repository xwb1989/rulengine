package dfa

import (
	"testing"

	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xwb1989/rulengine/parser"
	"math/rand"
)

func TestDFA(t *testing.T) {
	Convey("let's have some fun...", t, func() {
		Convey("make sure no false negative", func() {
			words := squareTable(1000, 10)
			rules := WordsToRules(words, false)
			dfa := MakeDFA(rules)
			for _, w := range words {
				act := dfa.GetAction(w)
				So(act, ShouldNotBeNil)
				res := act.Apply(w).(string)
				So(res, ShouldEqual, "accept")
			}
			//rand is deterministic thus we can always expect same size
			So(dfa.Size(), ShouldEqual, 6055)
		})
		Convey("make sure we can generate minimal DFA...", func() {
			Convey("we have luka...", func() {
				words := []string{"luka"}
				rules := WordsToRules(words, true)
				dfa := MakeDFA(rules)
				So(dfa.Size(), ShouldEqual, 5)
				Convey("then ryan...", func() {
					words := append(words, "ryan")
					rules := WordsToRules(words, true)
					dfa := MakeDFA(rules)
					So(dfa.Size(), ShouldEqual, 8)
					Convey("then grady", func() {
						words := append(words, "grady")
						rules := WordsToRules(words, true)
						dfa := MakeDFA(rules)
						So(dfa.Size(), ShouldEqual, 12)
						Convey("then ge...", func() {
							words := append(words, "ge")
							rules := WordsToRules(words, true)
							dfa := MakeDFA(rules)
							So(dfa.Size(), ShouldEqual, 12)
							Convey("then ben...", func() {
								words := append(words, "ben")
								rules := WordsToRules(words, true)
								dfa := MakeDFA(rules)
								So(dfa.Size(), ShouldEqual, 13)
								Convey("finally zach...", func() {
									words := append(words, "zach")
									rules := WordsToRules(words, true)
									dfa := MakeDFA(rules)
									So(dfa.Size(), ShouldEqual, 16)

								})
							})
						})
					})
				})
			})
		})
	})
}

func BenchmarkMakeDFA(b *testing.B) {
	words := squareTable(10000, 10)
	rules := WordsToRules(words, false)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		MakeDFA(rules)
	}
}

func BenchmarkGetAction(b *testing.B) {
	words := squareTable(10000, 10)
	rules := WordsToRules(words, false)
	dfa := MakeDFA(rules)
	word := words[rand.Int31n(10000)]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dfa.GetAction(word)
	}
}

func WordsToRules(words []string, isMin bool) []*Rule {
	rules := make([]*Rule, len(words))
	for i, word := range words {
		rules[i] = wordToRule(word, isMin)
	}
	return rules
}

var pred_map map[string]*Predicate = make(map[string]*Predicate)
var act_map map[string]*Action = make(map[string]*Action)

//a tiny parser
func wordToRule(w string, isMin bool) *Rule {
	preds := []*Predicate{}
	for i := 0; i < len(w); i++ {
		var expr string
		if isMin {
			expr = fmt.Sprintf("d[i]==%v", string(w[i]))
		} else {
			expr = fmt.Sprintf("d[%v]==%v", i, string(w[i]))
		}
		pred, ok := pred_map[expr]

		if !ok {
			j := i
			predFunc := func(data interface{}) bool {
				word := data.(string)
				if j >= len(word) {
					return false
				}
				return word[j] == w[j]
			}
			pred = MakePredicate(expr, predFunc)
			pred_map[expr] = pred
		}
		preds = append(preds, pred)
	}
	actFunc := func(data interface{}) interface{} {
		return "accept"
	}
	act, ok := act_map["accept"]
	if !ok {
		act = MakeAction("accept", actFunc)
		act_map["accept"] = act
	}
	rule := MakeRule(preds, act)
	return rule
}

func squareTable(l int, w int) []string {
	ret := []string{}
	for i := 0; i < l; i++ {
		ret = append(ret, randSeq(w))
	}
	return ret
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
