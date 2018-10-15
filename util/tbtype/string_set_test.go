package tbtype_test

import (
	"testing"

	"github.com/maprost/testbox/should"
	"github.com/maprost/toolbox/util/tbtype"
)

func TestStringSet(t *testing.T) {
	s := tbtype.NewStringSet()

	// add
	s.Add("blob")
	s.Add("drop")

	// contains
	should.BeTrue(t, s.Contains("blob"))
	should.BeTrue(t, s.Contains("drop"))
	should.BeFalse(t, s.Contains("what?"))

	// len
	should.BeEqual(t, s.Len(), 2)

	// add list
	s.AddList([]string{"hello", "world"})

	// contains
	should.BeTrue(t, s.Contains("blob"))
	should.BeTrue(t, s.Contains("drop"))
	should.BeTrue(t, s.Contains("hello"))
	should.BeTrue(t, s.Contains("world"))

	// len
	should.BeEqual(t, s.Len(), 4)

	// keys
	should.BeSimilar(t, s.Keys(), []string{"blob", "drop", "hello", "world"})
}

func TestNewStringSet(t *testing.T) {
	s := tbtype.NewStringSet("joa")

	// contains
	should.BeTrue(t, s.Contains("joa"))

	// len
	should.BeEqual(t, s.Len(), 1)
}

func TestNewStringSet_emptyString(t *testing.T) {
	s := tbtype.NewStringSet("")

	// len
	should.BeEqual(t, s.Len(), 0)
}

func TestStringSet_addEmptyString(t *testing.T) {
	s := tbtype.NewStringSet()

	// add
	s.Add("")

	// len
	should.BeEqual(t, s.Len(), 0)

	// keys
	should.BeSimilar(t, s.Keys(), []string{})
}

func TestStringSet_addStringSet(t *testing.T) {
	s1 := tbtype.NewStringSet("blob")
	s2 := tbtype.NewStringSet("drop")

	s1.AddSet(s2)

	// contains (s1)
	should.BeTrue(t, s1.Contains("blob"))
	should.BeTrue(t, s1.Contains("drop"))

	// contains (s2)
	should.BeTrue(t, s2.Contains("drop"))
	should.BeFalse(t, s2.Contains("blob"))

	// len
	should.BeEqual(t, s1.Len(), 2)
	should.BeEqual(t, s2.Len(), 1)
}
