package main

import (
	"html/template"
)

// HTMLテンプレートの読み込み
func loadTemplates() {
	var baseTemplate = "templates/layout.html"
	templates = make(map[string]*template.Template)
	// 各HTMLテンプレートに共通レイアウトを適用した結果をmapに保存する
	templates["index"] = template.Must(
		template.ParseFiles(baseTemplate, "templates/index.html"))
	templates["error"] = template.Must(
		template.ParseFiles(baseTemplate, "templates/error.html"))
	templates["user"] = template.Must(
		template.ParseFiles(baseTemplate, "templates/user.html"))
	templates["login"] = template.Must(
		template.ParseFiles(baseTemplate, "templates/login.html"))
	templates["admin"] = template.Must(
		template.ParseFiles(baseTemplate, "templates/admin.html"))
}
