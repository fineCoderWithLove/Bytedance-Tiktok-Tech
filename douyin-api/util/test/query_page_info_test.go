package test

import (
	"testing"
)

func TestPostPublishHandler(t *testing.T) {

	// 初始化测试用的 Gin 引擎

	//userId := int64(1)
	//// 生成 token
	//token, err := util.GenerateToken(userId, time.Hour)
	//if err != nil {
	//	t.Errorf("Failed to generate token: %v", err)
	//}
	//fmt.Println(token)
	//
	//// 创建一个测试请求
	//body := "content=测试内容"
	//req, err := http.NewRequest("POST", "localhost:8080/douyin/favorite/action?video_id=1&action_type=1&token="+token, strings.NewReader(body))
	//if err != nil {
	//	t.Fatalf("Failed to create request: %v", err)
	//}
	//
	//// 创建一个 ResponseRecorder 来记录响应
	//rr := httptest.NewRecorder()
	//
	//// 执行请求
	////engine.ServeHTTP(rr, req)
	//
	//// 检查响应状态码
	//if rr.Code != http.StatusOK {
	//	t.Errorf("Expected status %v, but got %v", http.StatusOK, rr.Code)
	//}
	//
	//// 输出响应内容
	//fmt.Println(rr.Body.String())
}
