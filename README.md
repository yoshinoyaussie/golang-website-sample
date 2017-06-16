# golang-website-sample
Go言語でWebサイトを作ってみるサンプルです。フレームワークはEchoを使用しています。

## コード全体構成

```
/
└─webserver
    │  handler.go  リクエストハンドラの定義
    │  server.go   サーバーのメイン処理
    │  static.go   静的ファイルパスの定義
    │  template.go HTMLテンプレートの定義
    ├─public     静的ファイル
    │  ├─css       CSSファイル
    │  ├─img       画像ファイル
    │  └─js        JavaScriptファイル
    ├─session    セッション関連の処理
    │      cookie.go          セッションCookie関連
    │      manager.go         セッションデータ管理（公開関数）
    │      manager_local.go   セッションデータ管理（非公開関数）
    ├─setting    設定関連の処理
    │      setting.go         設定データの定義
    └─templates  HTMLテンプレート
            error.html        エラーメッセージ画面
            index.html        index画面
            layout.html       共通レイアウト
```
