package main

import (
    "context"
    "fmt"
    "time"

    "github.com/xsymphony/notify"
)

func main() {
    notify.Watch(context.TODO(), "fanout", func(msg interface{}) {
        fmt.Println("[watch][A]", msg)
    })
    notify.Watch(context.TODO(), "fanout", func(msg interface{}) {
        fmt.Println("[watch][B]", msg)
    })
    notify.Watch(context.TODO(), "fanout", func(msg interface{}) {
        fmt.Println("[watch][C]", msg)
    })

    for i := 0; i < 10; i++ {
        notify.Send(context.TODO(), "fanout", i)
    }

    time.Sleep(2*time.Second)
}