package test

import (
	"douyin/douyin-api/util"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestGetToken(t *testing.T) {
	// 测试获取 token
	userID := int64(24)
	expiration := time.Hour * 24

	// 生成 token
	token, err := util.GenerateToken(userID, expiration)
	if err != nil {
		t.Errorf("Failed to generate token: %v", err)
	}

	fmt.Println(token)
	// 从 Redis 中获取 uid
	result, err := util.GetUserId(token)
	if err != nil {
		t.Errorf("Failed to get token from Redis: %v", err)
	}

	// 验证获取的 token 是否正确
	if result != strconv.FormatInt(userID, 10) {
		t.Errorf("Expected token %s, but got %s", token, result)
	}
}
