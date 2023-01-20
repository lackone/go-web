package middlewares

import (
	"context"
	"github.com/lackone/go-web/framework"
	"log"
	"time"
)

func Timeout(d time.Duration) framework.ControllerHandler {
	return func(c *framework.Context) error {
		finishChan := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		ctx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		c.GetRequest().WithContext(ctx)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			//执行具体的业务逻辑
			c.Next()

			finishChan <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			log.Println(p)
			c.SetStatus(500).Json("panic")
		case <-finishChan:
			log.Println("finish")
		case <-ctx.Done():
			c.SetIsTimeout()
			c.SetStatus(500).Json("time out")
		}

		return nil
	}
}
