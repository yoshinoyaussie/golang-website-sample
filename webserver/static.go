package main

import (
	"github.com/labstack/echo"
)

// 静的ファイルを配置するルーティングを設定
func setStaticRoute(e *echo.Echo) {
	e.Static("/public/css/", "./public/css")
	e.Static("/public/js/", "./public/js/")
	e.Static("/public/img/", "./public/img/")
}
