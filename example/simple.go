package main

import (
	"time"

	"github.com/codcodog/progress"
)

func main() {
	p := progress.New(100)
	for i := 0; i < 100; i++ {
		p.Incr()
		time.Sleep(time.Millisecond * 20)
	}
}
