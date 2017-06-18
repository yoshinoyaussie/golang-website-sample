package main

import (
	"net/http"

	"./model"

	"github.com/labstack/echo"
)

// ルーティングに対応するハンドラを設定します。
func setRoute(e *echo.Echo) {
	e.GET("/", handleIndexGet)
	e.GET("/users/:user_id", handleUsersGet)
}

// GET:/
func handleIndexGet(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "world")
}

// GET:/users/:id
func handleUsersGet(c echo.Context) error {
	users, err := userDA.FindByUserID(c.Param("user_id"), model.FindFirst)
	if err != nil {
		return c.Render(http.StatusOK, "error", err)
	}
	user := users[0]
	return c.Render(http.StatusOK, "user", user)
}
