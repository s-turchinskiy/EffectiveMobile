package closer

import (
	"context"
	"fmt"
	"github.com/s-turchinskiy/EffectiveMobile/internal/common/reflect"
	"github.com/s-turchinskiy/EffectiveMobile/internal/middleware/logger"
	"strings"
	"sync"
)

type Closer struct {
	mu    sync.Mutex
	funcs []Func
}

func (c *Closer) Add(f Func) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, f)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		msgs     = make([]string, 0, len(c.funcs))
		complete = make(chan struct{}, 1)
	)

	go func() {
		for _, f := range c.funcs {
			logger.Log.Debugw("stopping " + reflect.GetFunctionName(f))
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
		return fmt.Errorf("shutdown cancelled: %v", ctx.Err())
	}

	if len(msgs) > 0 {
		return fmt.Errorf(
			"shutdown finished with error(s): \n%s",
			strings.Join(msgs, "\n"),
		)
	}

	return nil
}

type Func func(ctx context.Context) error
