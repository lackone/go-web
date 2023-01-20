# go-web
基于go标准库实现的简单web框架

### 简单使用
```go
package main

import (
	"context"
	"github.com/lackone/go-web/framework"
	"github.com/lackone/go-web/framework/middlewares"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func registerRouter(core *framework.Core) {
	core.Get("/test", TestHandler)

	group := core.Group("/aaa")
	group.Use(middlewares.Recovery(), middlewares.Time())
	{
		group.Get("/bbb", TestHandler)

		group.Get("/ccc/:id", TestHandler)
	}
}

func TestHandler(c *framework.Context) error {
	c.Json("test")
	return nil
}

func main() {
	core := framework.NewCore()

	registerRouter(core)

	server := &http.Server{
		Addr:    ":8080",
		Handler: core,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("server listen err:", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	if err := server.Shutdown(timeout); err != nil {
		log.Fatalln(err)
	}

	log.Println("server shutdown")
}
```