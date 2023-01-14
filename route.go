package main

import "github.com/lackone/go-web/framework"

func registerRouter(core *framework.Core) {
	core.Get("/test", TestHandler)

	group := core.Group("/aaa")
	{
		group.Get("/ccc", TestHandler)

		group.Get("/ccc/:id", TestHandler)
	}
}

func TestHandler(c *framework.Context) error {
	c.Text(200, "hello")
	return nil
}
