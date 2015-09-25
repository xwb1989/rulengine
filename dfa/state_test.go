package dfa

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TeststateBasic(t *testing.T) {
	Convey("create two initial states", t, func() {
		s1 := initState()
		s2 := initState()
		Convey("they should be equivalent", func() {
			So(s1.Equals(s2), ShouldBeTrue)
			So(s1.Equals(s1), ShouldBeTrue)
		})
		pred := MakePredicate("", nil)
		Convey("modify first one then they should be different", func() {
			s1.SetNext(pred, s2)
			So(s1.Equals(s2), ShouldBeFalse)
		})
		Convey("modify both thus they have same edges, then they should be equivalent", func() {
			s3 := initState()
			s1.SetNext(pred, s3)
			s2.SetNext(pred, s3)
			So(s1.Equals(s2), ShouldBeTrue)
			Convey("but if one has action, they should be different", func() {
				act := MakeAction("", nil)
				s2.Action = act
				So(s1.Equals(s2), ShouldBeFalse)
				Convey("even both has action, but if actions are different, they are different", func() {
					s1.Action = MakeAction("1", nil)
					So(s1.Equals(s2), ShouldBeFalse)
					Convey("but if they have the same action, then they should be equivalent", func() {
						s1.Action = act
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
