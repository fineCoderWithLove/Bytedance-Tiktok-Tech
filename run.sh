
nohup go run douyin-api/cmd/main.go > api.log 2>&1 &

nohup go run base-service/cmd/main.go > base.log 2>&1 &

nohup go run interaction-service/server/comment/comment.go > comment.log 2>&1 &

nohup go run interaction-service/server/favorite/favorite.go > favorite.log 2>&1 &

nohup go run social-service/cmd/main.go > social.log 2>&1 &

echo "当显示404not found的时候，项目运行成功"
wait