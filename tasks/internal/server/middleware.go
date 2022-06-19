package server

import "net/http"

type middlewareFunc func(http.HandlerFunc) http.HandlerFunc

func nest(h http.HandlerFunc, middlewareFuncs ...middlewareFunc) http.HandlerFunc {
	for i := len(middlewareFuncs) - 1; i >= 0; i-- {
		h = middlewareFuncs[i](h)
	}

	return h
}
