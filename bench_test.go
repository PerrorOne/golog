package golog

import (
	"testing"
)

func BenchmarkWrite(b *testing.B){
	num:=10
	InitLogger("", 0, false)
	for i:=0;i<b.N;i++{
		Info("aaa: %d", num)
	}
}