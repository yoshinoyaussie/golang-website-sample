package main

import (
	"context"
	"html/template"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"./session"
	"./setting"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

// レイアウト適用済のテンプレートを保存するmap
var templates map[string]*template.Template

// セッション管理のインスタンス
var sessionManager *session.Manager

// Template はHTMLテンプレートを利用するためのRenderer Interfaceです。
type Template struct {
}

// Render はHTMLテンプレートにデータを埋め込んだ結果をWriterに書き込みます。
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return templates[name].ExecuteTemplate(w, "layout.html", data)
}

func main() {
	// Echoのインスタンスを生成
	e := echo.New()

	// ログの出力レベルを設定
	//	e.Logger.SetLevel(log.INFO)
	e.Logger.SetLevel(log.DEBUG)

	// テンプレートを利用するためのRendererの設定
	t := &Template{}
	e.Renderer = t

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 静的ファイルを配置するルーティングを設定
	setStaticRoute(e)

	// 各ルーティングに対するハンドラを設定
	setRoute(e)

	// セッション管理を開始
	sessionManager = &session.Manager{}
	sessionManager.Start(e)

	// サーバーを開始
	go func() {
		if err := e.Start(setting.Server.Port); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// 中断を検知したらリクエストの完了を10秒まで待ってサーバーを終了する
	// (Graceful Shutdown)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Info(err)
		e.Close()
	}

	// セッション管理を停止
	sessionManager.Stop()

	// 終了ログが出るまで少し待つ
	time.Sleep(1 * time.Second)
}

// 初期化を行う
func init() {
	// 設定の読み込み
	setting.Load()
	// HTMLテンプレートの読み込み
	loadTemplates()
}
