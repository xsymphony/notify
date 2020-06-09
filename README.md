# notify
notify包用于单应用内的消息通信，提供`发布/订阅`、`取消订阅`、`超时`、`超时处理`、`扇出`、`工作池`等功能。

## 基础功能
```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/xsymphony/notify"
)

func main() {
    // consume from the "simple" topic
    notify.Watch(context.TODO(), "simple", func(msg interface{}) {
        fmt.Println("[watch]", msg)
    })
    
    // producer msg to the "simple" topic
    for i := 0; i < 10; i++ {
        notify.Send(context.TODO(), "simple", i)
    }
    time.Sleep(1*time.Second)
}

// stdout:
// [watch]1
// [watch]2
// [watch]3
//   ...
// [watch]9
```

## 消费超时设置
```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/xsymphony/notify"
)

func main() {
    var fn = func(msg interface{}) {
        time.Sleep(1*time.Second)
        fmt.Println("[watch]", msg)
    }
    var timeoutFn = func(msg interface{}) {
        fmt.Println("[watch]msg timeout ", msg)
    }

    notify.Watch(context.TODO(), "timeout", fn, notify.Timeout(500*time.Millisecond), notify.TimeoutFunc(timeoutFn))

    notify.Send(context.TODO(), "timeout", 1)
    
    time.Sleep(1*time.Second)
}

// stdout:
// [watch]msg timeout 1
```

## 扇出
```go
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

    notify.Send(context.TODO(), "fanout", 1)

    time.Sleep(1*time.Second)
}

// stdout:
// [watch][A]1  
// [watch][C]1 
// [watch][B]1 
```

## 工作池
```go
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

// stdout:
// [watch] id: 49 0
// [watch] id: 50 1
// [watch] id: 49 3
// [watch] id: 51 2
// [watch] id: 50 4
// [watch] id: 51 6
// [watch] id: 49 5
// [watch] id: 50 7
// [watch] id: 49 9
// [watch] id: 51 8
```