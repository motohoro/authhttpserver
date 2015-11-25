# authhttpserver

autobasicauth http server

[firefoxdecrypt](https://github.com/motohoro/firefoxdecrypt)と組み合わせて、BasicAuthなhtmlやrssなどを、ブラウザ経由せず自動でfirefoxに保存されてるID,PWを取得してGETするためのserver

### ビルド

go build authhttpserver.go

go 1.4でビルドはexeが強制終了になることあり、go1.5でビルドは今のところ安定

### 起動

カレントディレクトリをexeファイルの場所に変更してからexe起動すること

### アクセスするURL

http://localhost:8087/?url=TARGET_URL&buser=BASICAUTH_USERNAME

