package main

import (
    "context"
    "fmt"

    "github.com/xsymphony/notify"
)

func main() {
    notify.Watch(context.TODO(), "simple", func(msg interface{}) {
        fmt.Println("[watch]", msg)
    })

    for i := 0; i < 10; i++ {
        notify.Send(context.TODO(), "simple", i)
    }
}
