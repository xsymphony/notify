package main

import (
    "context"
    "fmt"
    "time"

    "github.com/xsymphony/notify"
)

func main() {
    var fn = func(msg interface{}) {
        time.Sleep(3*time.Second)
        fmt.Println("[watch]", msg)
    }
    var timeoutFn = func(msg interface{}) {
        fmt.Println("[watch]msg timeout", msg)
    }

    notify.Watch(context.TODO(), "timeout", fn, notify.Timeout(1*time.Second), notify.TimeoutFunc(timeoutFn))

    for i := 0; i < 10; i++ {
        notify.Send(context.TODO(), "timeout", i)
    }
}