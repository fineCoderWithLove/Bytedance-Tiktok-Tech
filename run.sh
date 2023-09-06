
nohup go run app/gateway/cmd/main.go > gateway.log 2>&1 &

nohup go run app/user/cmd/main.go > user.log 2>&1 &

nohup go run app/relation/cmd/main.go > relation.log 2>&1 &

nohup go run app/video/cmd/main.go > video.log 2>&1 &

echo "项目已成功启动！"
wait
