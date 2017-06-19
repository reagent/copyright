package app

import (
	"database/sql"
	"net/http"
	"regexp"
)

type App struct {
	DB     *sql.DB
	Routes Routes
}

type Handler func(ctx *Context)

func New(db *sql.DB) *App {
	return &App{DB: db, Routes: Routes{}}
}

func (a *App) AddRoute(method string, pattern string, handler Handler) {
	re := regexp.MustCompile(pattern)
	a.Routes = append(a.Routes, Route{Pattern: re, Method: method, Handler: handler})
}

func (a *App) GET(pattern string, handler Handler) {
	a.AddRoute("GET", pattern, handler)
}

func (a *App) POST(pattern string, handler Handler) {
	a.AddRoute("POST", pattern, handler)
}

func (a *App) Dispatch(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		App:      a,
		Request:  NewRequest(r),
		Response: NewResponse(w),
	}

	for _, route := range a.Routes {
		matches, ok := route.MatchURL(ctx.URL.Path)

		if ok && ctx.Method == route.Method {
			ctx.params = matches

			route.Handler(ctx)
			return
		}
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Dispatch(w, r)
}
