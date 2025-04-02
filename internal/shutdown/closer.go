package shutdown

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

var closer *Closer

type Func func(ctx context.Context) error

type Closer struct {
	mu    sync.Mutex
	funcs []Func
}

func Init() {
	closer = &Closer{}
}

func GetCloser() *Closer {
	return closer
}

func (c *Closer) Add(f Func) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, f)
}

func (c *Closer) Close(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	msgs := make([]string, 0, len(c.funcs))
	complete := make(chan struct{}, 1)

	go func() {
		for _, f := range c.funcs {
			if err := f(ctx); err != nil {
				msgs = append(msgs, fmt.Sprintf("[!] %v", err))
			}
		}

		complete <- struct{}{}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		println("shutdown cancelled by timeout: " + ctx.Err().Error())
	}

	if len(msgs) > 0 {
		println("shutdown finished with error(s):\n" + strings.Join(msgs, "\n"))
	}

	println(`Bye Bye`)
}
