package object

import "testing"

const (
	errDiffHashKeys = "strings with the same content have different hash keys"
	errSameHashKeys = "strings with different content have same hash keys"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Error(errDiffHashKeys)
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Error(errDiffHashKeys)
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Error(errSameHashKeys)
	}
}
