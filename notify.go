package notify

import (
    "context"
    "sync"
    "time"
)

type Notify struct {
    consumer sync.Map
    producer sync.Map
}

func (notify *Notify) Send(ctx context.Context, topic string, msg interface{}) {
    v, ok := notify.producer.LoadOrStore(topic, make(chan interface{}))
    if !ok {
        notify.exchange(ctx, topic)
    }
    v.(chan interface{}) <- msg
}

func (notify *Notify) exchange(ctx context.Context, topic string) {
    v, _ := notify.producer.Load(topic)
    ch := v.(chan interface{})
    go func() {
        for msg := range ch {
            v, ok := notify.consumer.Load(topic)
            if !ok {
                continue
            }
            v.(*consumers).iter(func(h handler) {
                if h.timeout <= 0 {
                    h.ch <- msg
                    return
                }
                select {
                case h.ch <- msg:
                case <-time.After(h.timeout):
                    if h.timeoutFunc != nil {
                        h.timeoutFunc(msg)
                    }
                }
            })
        }
    }()
}

func (notify *Notify) Watch(ctx context.Context, topic string, f func(msg interface{}), options ...Options) (unwatch func()) {
    h := handler{
        quit:       make(chan struct{}, 1),
        ch:         make(chan interface{}),
        fn:         f,
        workerSize: 1,
    }
    for _, option := range options {
        h = option(h)
    }
    go func() {
        // single worker
        if h.workerSize <= 1 {
            for {
                select {
                case msg := <-h.ch:
                    f(msg)
                case <-h.quit:
                    return
                }
            }
        }
        // multi workers
        ctr := make([]chan struct{}, h.workerSize)
        for i := uint32(0); i < h.workerSize; i++ {
            quit := make(chan struct{}, 1)
            ctr[i] = quit
            go func() {
                for {
                    select {
                    case msg := <-h.ch:
                        f(msg)
                    case <-quit:
                       return
                    }
                }
            }()
        }
        <-h.quit
        for _, ch := range ctr {
           ch<- struct{}{}
        }
    }()
    v, ok := notify.consumer.LoadOrStore(topic, &consumers{
        consumers: []handler{h},
    })
    if ok {
        v.(*consumers).append(h)
    }
    return func() {
        v.(*consumers).remove(h)
    }
}

type consumers struct {
    sync.RWMutex
    consumers []handler
}

func (c *consumers) append(p handler) {
    c.Lock()
    defer c.Unlock()
    c.consumers = append(c.consumers, p)
}

func (c *consumers) remove(target handler) {
    c.Lock()
    defer c.Unlock()
    for i, el := range c.consumers {
        if el.ch == target.ch {
            c.consumers = append(c.consumers[:i], c.consumers[i+1:]...)
            el.quit <- struct{}{}
            return
        }
    }
}

func (c *consumers) iter(f func(h handler)) {
    c.RLock()
    defer c.RUnlock()
    for _, consumer := range c.consumers {
        f(consumer)
    }
}

type handler struct {
    fn   func(msg interface{})
    ch   chan interface{}
    quit chan struct{}

    workerSize uint32

    timeout     time.Duration
    timeoutFunc func(msg interface{})
}

type Options func(h handler) handler

func Timeout(d time.Duration) Options {
    return func(h handler) handler {
        h.timeout = d
        return h
    }
}

func TimeoutFunc(f func(msg interface{})) Options {
    return func(h handler) handler {
        h.timeoutFunc = f
        return h
    }
}

func WorkerSize(size uint32) Options {
    return func(h handler) handler {
        h.workerSize = size
        return h
    }
}

var _defaultNotify Notify

func Send(ctx context.Context, topic string, msg interface{}) {
    _defaultNotify.Send(ctx, topic, msg)
}

func Watch(ctx context.Context, topic string, f func(msg interface{}), options ...Options) (unwatch func()) {
    return _defaultNotify.Watch(ctx, topic, f, options...)
}
