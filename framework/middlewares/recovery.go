package middlewares

import "github.com/lackone/go-web/framework"

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.SetStatus(500).Json(err)
			}
		}()

		c.Next()

		return nil
	}
}
