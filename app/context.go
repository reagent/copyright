package app

import (
	"database/sql"
)

type Context struct {
	*Request
	*Response

	App *App

	params map[string]string
}

func (c *Context) Param(key string) (string, bool) {
	v, ok := c.params[key]

	return v, ok
}

func (c *Context) DB() *sql.DB {
	return c.App.DB
}
