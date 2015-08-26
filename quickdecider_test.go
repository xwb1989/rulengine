package quickdecider

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBasic(t *testing.T) {
	convey.Convey("should has some base functionality...", t, func() {
		convey.Convey("should be able to create new Decider", func() {
			decider, err := MakeDecider(map[string]interface{}{})
			convey.So(err, convey.ShouldBeNil)
			convey.Convey("should be able to get new action", func() {
				_, err := decider.GetAction(map[string]interface{}{})
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}
