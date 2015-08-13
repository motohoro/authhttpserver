# authhttpserver

autobasicauth http server

BasicAuthなhtmlやrssなどを、ブラウザ経由せず自動でfirefoxに保存されてるID,PWを取得してGETするためのserver

### ビルド

go build authhttpserver.go

### アクセスするURL

http://localhost:8087/?url=TARGET_URL&buser=BASICAUTH_USERNAME

スタートアップの時カレントディレクトリをexeファイルの場所に変更してからexe起動すること
