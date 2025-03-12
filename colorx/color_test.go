package colorx_test

import (
	"fmt"
	"testing"

	"github.com/maprost/testbox/should"
	"github.com/maprost/toolbox/colorx"
)

func TestRemove(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		msg := colorx.Text("hallo", colorx.Green)
		should.BeEqual(t, msg, fmt.Sprintf("%shallo%s", colorx.Green, colorx.ColorOff))

		msg = colorx.Remove(msg)
		should.BeEqual(t, msg, "hallo")
	})

	t.Run("multi", func(t *testing.T) {
		msg1 := colorx.Text("hallo", colorx.Green)
		msg2 := colorx.Text("welt", colorx.Cyan)
		msg := fmt.Sprintf("lol %s%s jo %s", msg1, msg2, msg2)
		should.BeEqual(t, msg, fmt.Sprintf("lol %shallo%s%swelt%s jo %swelt%s",
			colorx.Green, colorx.ColorOff,
			colorx.Cyan, colorx.ColorOff,
			colorx.Cyan, colorx.ColorOff,
		))

		msg = colorx.Remove(msg)
		should.BeEqual(t, msg, "lol hallowelt jo welt")
	})

	t.Run("list", func(t *testing.T) {
		msgs := []string{
			colorx.Text("hallo", colorx.Green),
			colorx.Text("welt", colorx.Cyan),
		}
		msgs = colorx.Removes(msgs)
		should.BeEqual(t, msgs[0], "hallo")
		should.BeEqual(t, msgs[1], "welt")
	})
}
