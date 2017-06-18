# golang-website-sample
Go言語でWebサイトを作ってみるサンプルです。フレームワークは
Echo https://echo.labstack.com/ 
を使用しています。

## 概要

WebサイトのサーバーサイドをGoで一通り作っていっています。
詳細につきましては以下のQiita記事を参照してください。

### Go言語でWebサイトを作ってみる：

* Hello World的な http://qiita.com/y_ussie/items/ca8dc5e423eec318a436
* リクエストパラメータの扱い http://qiita.com/y_ussie/items/2d397e70bfc38f75ca51
* セッションデータストアを作る http://qiita.com/y_ussie/items/b1db86b0b54ec69bb928
* Cookieを使用したセッション管理 http://qiita.com/y_ussie/items/00e542cb3531b48fd21a
* ひとまずコード整理 http://qiita.com/y_ussie/items/12bb4fd8cefb740581f8
* ユーザー情報をJSONから読み込んで参照してみる http://qiita.com/y_ussie/items/8704ce209704bf191e63

## コード全体構成

```
/
└─webserver
    │  handler.go  リクエストハンドラの定義
    │  server.go   サーバーのメイン処理
    │  static.go   静的ファイルパスの定義
    │  template.go HTMLテンプレートの定義
    ├─data       JSONファイルなど
    │  users.json  ユーザー情報のJSONファイル
    ├─model      データモデルとアクセサ
    │  user.go     ユーザー情報のモデルとアクセサ
    ├─public     静的ファイル
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
            layout.html       共通レイアウト
            user.html         ユーザー情報の表示画面
```
