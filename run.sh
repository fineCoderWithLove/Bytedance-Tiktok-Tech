
nohup go run base-service/cmd/main.go   &

nohup go run douyin-api/cmd/main.go  &

nohup go run interaction-service/server/comment/comment.go  &

nohup go run interaction-service/server/favorite/favorite.go  &

nohup go run social-service/cmd/main.go  &

echo "项目已成功启动！"
wait