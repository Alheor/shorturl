// Package shutdown - сервис механизма graceful shutdown.
//
// # Описание
//
// Package shutdown принимает и сохраняет код завернутый в функцию, которая будет выполнена перед тем, как сервис завершит работу.
package shutdown

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

var closer *Closer

// Func - функция в которую заворачивается код, для выполнен перед остановкой сервиса.
type Func func(ctx context.Context) error

// Closer - структура хранения действий, выполняемых после подачи команды на завершение работы.
type Closer struct {
	mu    sync.Mutex
	funcs []Func
}

// Init Подготовка graceful shutdown.
func Init() {
	closer = &Closer{}
}

// GetCloser - получение экземпляра.
func GetCloser() *Closer {
	return closer
}

// Add Добавление произвольного кода, завернутого в функцию Func.
func (c *Closer) Add(f Func) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, f)
}

// Close Выполнение кода.
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
