package main

import (
    "bytes"
    "context"
    "fmt"
    "runtime"
    "strconv"

    "github.com/xsymphony/notify"
)

func getGID() uint64 {
    b := make([]byte, 64)
    b = b[:runtime.Stack(b, false)]
    b = bytes.TrimPrefix(b, []byte("goroutine "))
    b = b[:bytes.IndexByte(b, ' ')]
    n, _ := strconv.ParseUint(string(b), 10, 64)
    return n
}

func main() {
    notify.Watch(context.TODO(), "simple", func(msg interface{}) {
        fmt.Println("[watch]", "id:", getGID(), msg)
    }, notify.WorkerSize(3))

    for i := 0; i < 10; i++ {
        notify.Send(context.TODO(), "simple", i)
    }
}
