package setting

import (
	"time"
)

// Server はサーバーの動作に関する設定です。
var Server = server{}

type server struct {
	Port string
}

// Session はセッションに関する設定です。
var Session = session{}

type session struct {
	CookieName   string
	CookieExpire time.Duration
}

// Load は設定を読み込みます。
func Load() {
	// ポート番号
	Server.Port = ":3000"
	// セッションのCookie名
	Session.CookieName = "gowebserver_session_id"
	// セッションのCookie有効期限
	Session.CookieExpire = (1 * time.Hour)
}
