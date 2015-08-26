package dfa

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xwb1989/quickdecider/parser"
)

func TestStateBasic(t *testing.T) {
	Convey("create two initial states", t, func() {
		s1 := InitState()
		s2 := InitState()
		Convey("they should be equivalent", func() {
			So(s1.Equals(s2), ShouldBeTrue)
			So(s1.Equals(s1), ShouldBeTrue)
		})
		pred := MakePredicate("")
		Convey("modify first one then they should be different", func() {
			s1.SetNext(pred, s2)
			So(s1.Equals(s2), ShouldBeFalse)
		})
		Convey("modify both thus they have same edges, then they should be equivalent", func() {
			s3 := InitState()
			s1.SetNext(pred, s3)
			s2.SetNext(pred, s3)
			So(s1.Equals(s2), ShouldBeTrue)
			Convey("but if one has action, they should be different", func() {
				act := MakeAction("")
				s2.SetAction(act)
				So(s1.Equals(s2), ShouldBeFalse)
				Convey("even both has action, but if actions are different, they are different", func() {
					s1.SetAction(MakeAction(""))
					So(s1.Equals(s2), ShouldBeFalse)
					Convey("but if they have the same action, then they should be equivalent", func() {
						s1.SetAction(act)
						So(s1.Equals(s2), ShouldBeTrue)
					})
				})
			})
			Convey("now s3 should be confluent", func() {
				So(s3.IsConfluent(), ShouldBeTrue)
				Convey("then we change s2, s3 should not be confluent", func() {
					s2.SetNext(pred, s1)
					So(s3.IsConfluent(), ShouldBeFalse)
				})
			})
		})
	})
}
