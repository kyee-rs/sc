package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/prophittcorey/tor"
	"net/http"
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

func torMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if isTorExitNode(c.IP()) {
			return c.SendStatus(http.StatusTeapot)
		}

		return c.Next()
	}
}
