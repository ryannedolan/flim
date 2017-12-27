package main

import (
	"fmt"
	"github.com/ryannedolan/flim/fl"
)

func main() {
	foobar := fl.Promise(fl.Success("foo")).Iter().Map(` x + "bar" `).List()
	fmt.Println(foobar)

	nobar := fl.Promise(fl.Failuref("foo")).OrElse(fl.Success("no")).Iter().Map(` x + "bar" `).List()
	fmt.Println(nobar)

	rpc := func() (string, error) {
		return "apples", nil
	}

	f := fl.Future()
	f.CompleteWith(func() (interface{}, error) { return rpc() })
	resp, err := f.Wait()
	fmt.Println("server says:", resp, err)
}
