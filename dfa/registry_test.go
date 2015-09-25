package dfa

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegisterBasic(t *testing.T) {
	states := []*state{}
	for i := 0; i < 10; i++ {
		states = append(states, initState())
	}
	Convey("should be able to create a register", t, func() {
		reg := makeRegistry()
		Convey("should be able to put a new and simple state", func() {
			ret := reg.getOrPut(states[0])
			So(ret, ShouldEqual, ret)
			pred := MakePredicate("", nil)
			act := MakeAction("", nil)
			Convey("should be able to put a different state", func() {
				states[2].SetNext(pred, states[0]) //same as state 2
				ret := reg.getOrPut(states[2])
				So(ret, ShouldEqual, ret)
				Convey("should fail to put if there is already an equivalent one in register", func() {
					states[3].SetNext(pred, states[0]) //same as state 2
					ret := reg.getOrPut(states[3])
					So(ret, ShouldNotEqual, states[3])
					Convey("and the returned should be same as previous one", func() {
						So(ret, ShouldEqual, states[2])
					})
					Convey("and the size should not change", func() {
						So(reg.Size(), ShouldEqual, 2)
					})
					Convey("but we should be able to remove its equivalent by just passing it to remove()", func() {
						So(reg.Remove(states[3]), ShouldBeTrue)
						So(reg.Size(), ShouldEqual, 1)
						Convey("then we should be able to put it into register", func() {
							ret := reg.getOrPut(states[3])
							So(ret, ShouldEqual, states[3])
							So(reg.Size(), ShouldEqual, 2)
						})
					})
					Convey("however removing a state that has not in register and has no equivalent should fail", func() {
						states[4].SetNext(pred, states[0])
						states[4].Action = act //same edge but with different action
						So(reg.Remove(states[4]), ShouldBeFalse)
						So(reg.Size(), ShouldEqual, 2)
						Convey("and we should be able to put it into register", func() {
							ret := reg.getOrPut(states[4])
							So(ret, ShouldEqual, states[4])
							So(reg.Size(), ShouldEqual, 3)
						})
					})
				})
			})
		})
	})
}
