package xlog

import (
	"testing"
	"time"
)

func TestPrintln(t *testing.T) {

	Println("abc", "abc", "pppppppp")

	AddSign("miner", "0 0/1 * * * *", "200601021504")

	Println("miner", "ccccc", "ccccccccccccccc")

	var i = 0
	for {
		time.Sleep(1 * time.Minute)

		Println("miner", time.Now(), "ccccc", "eeeeeeeeeeee")
		i++
		if i > 10 {
			return
		}
	}

}
