package main

import (
	"net/http"

	"./model"

	"github.com/labstack/echo"
)

// ルーティングに対応するハンドラを設定します。
func setRoute(e *echo.Echo) {
	e.GET("/", handleIndexGet)
	e.GET("/login", handleLoginGet)
	e.POST("/login", handleLoginPost)
	e.POST("/logout", handleLogoutPost)
	e.GET("/users/:user_id", handleUsers)
	e.POST("/users/:user_id", handleUsers)
}

// GET:/
func handleIndexGet(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "world")
}

// GET:/users/:user_id
// POST:/users/:user_id
func handleUsers(c echo.Context) error {
	userID := c.Param("user_id")
	err := CheckUserID(c, userID)
	if err != nil {
		c.Echo().Logger.Debugf("User Page[%s] Roll Error. [%s]", userID, err)
		msg := "ログインしていません。"
		return c.Render(http.StatusOK, "error", msg)
	}
	users, err := userDA.FindByUserID(c.Param("user_id"), model.FindFirst)
	if err != nil {
		return c.Render(http.StatusOK, "error", err)
	}
	user := users[0]
	return c.Render(http.StatusOK, "user", user)
}

// GET:/login
func handleLoginGet(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

// POST:/login
func handleLoginPost(c echo.Context) error {
	userID := c.FormValue("userid")
	password := c.FormValue("password")
	err := UserLogin(c, userID, password)
	if err != nil {
		c.Echo().Logger.Debugf("User[%s] Login Error. [%s]", userID, err)
		msg := "ユーザーIDまたはパスワードが誤っています。"
		data := map[string]string{"user_id": userID, "password": "", "msg": msg}
		return c.Render(http.StatusOK, "login", data)
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/users/"+userID)
}

// POST:/logout
func handleLogoutPost(c echo.Context) error {
	err := UserLogout(c)
	if err != nil {
		c.Echo().Logger.Debugf("User Logout Error. [%s]", err)
		return c.Render(http.StatusOK, "login", nil)
	}
	msg := "ログアウトしました。"
	data := map[string]string{"user_id": "", "password": "", "msg": msg}
	return c.Render(http.StatusOK, "login", data)
}
