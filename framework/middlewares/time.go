package middlewares

import (
	"github.com/lackone/go-web/framework"
	"log"
	"time"
)

func Time() framework.ControllerHandler {
	return func(c *framework.Context) error {
		start := time.Now()

		c.Next()

		end := time.Since(start)

		log.Printf("%s 用时 %v\n", c.GetRequest().RequestURI, end)

		return nil
	}
}
