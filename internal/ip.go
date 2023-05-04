package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prophittcorey/tor"
)

func isTorExitNode(address string) bool {
	res, err := tor.IsExitNode(address)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if res {
		return true
	}
	return false
}

func ipMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if isTorExitNode(c.RealIP()) {
				return MakeError(c, http.StatusForbidden, ts.HTTPErrors.TorNotAllowed)
			}
			return next(c)
		}
	}
}
