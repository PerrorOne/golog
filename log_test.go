package golog

import "testing"

func Test_log(t *testing.T) {
	InitLogger("",0,false)
	Access("two")
	Error("one")
}
