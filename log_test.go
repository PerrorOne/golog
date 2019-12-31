package golog

import (
	"fmt"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	fmt.Println(time.Now())
	InitLogger("../data",0,false)
	Access("two")
	Error("one")
}
