package tbarray

//import (
//	"rpp.de/BackendLib/base/testbase"
//	"rpp.de/BackendLib/base/testbase/assert"
//	"rpp.de/BackendLib/util/arrayutil"
//	"testing"
//)
//
//func TestExclude(t *testing.T) {
//	testbase.InitSimpleTest(t)
//
//	list := []string{"1", "2", "3", "4"}
//	exclude := []string{"3", "4", "5", "6"}
//
//	result := arrayutil.Exclude(list, exclude)
//	assert.Size(result, 2)
//	assert.Equal(result, []string{"1", "2"})
//
//	assert.Size(list, 4)
//	assert.Size(exclude, 4)
//}
//
//func TestExclude_sameList(t *testing.T) {
//	testbase.InitSimpleTest(t)
//
//	list := []string{"1", "2", "3", "4"}
//
//	result := arrayutil.Exclude(list, list)
//	assert.Size(result, 0)
//}
//
//func TestExclude_emptyExcludeList(t *testing.T) {
//	testbase.InitSimpleTest(t)
//
//	list := []string{"1", "2", "3", "4"}
//
//	result := arrayutil.Exclude(list, []string{})
//	assert.Size(result, 4)
//	assert.Equal(result, list)
//}
//
//func TestExclude_emptyList(t *testing.T) {
//	testbase.InitSimpleTest(t)
//
//	result := arrayutil.Exclude([]string{}, []string{})
//	assert.Size(result, 0)
//	assert.Equal(result, []string{})
//}
