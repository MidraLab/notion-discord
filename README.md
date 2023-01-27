# 定例の通知の自動化
* Notionの定例ページのタイトルの自動変更
* Notionの定例ページをdiscordに通知

# setup
1. install docker and docker desktop
2. clone this repository
3. cd `mtg-notification`
4. `docker-compose up -d --build` (docker imageのビルドと起動)
5. `docker-compose exec mtg-notification bash` (dockerの仮想環境内に入る)
6. `go run test.go` (テスト用のコードを実行)

# note
* You don't need to install go on your pc because you are creating a go environment in docker.