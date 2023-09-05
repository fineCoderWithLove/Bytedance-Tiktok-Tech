package test

import (
	"crypto/sha512"
	"demotest/base-service/global"
	"demotest/base-service/model"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRegisterUser(t *testing.T) {
	username := "郑梓桐"
	oldpassword := "20021211zzt"
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(oldpassword, options)
	newpassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(newpassword)

	var user model.User
	result := global.DB.Table("user").Where("name = ?", username).First(&user)
	fmt.Println("-------")
	fmt.Println(result)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			/* 用户名不存在,就把用户添加进去
			添加一条记录
			生成token
			*/
			fmt.Println("not exit")
			user := model.User{
				Id:       0,
				Name:     username,
				Password: newpassword,
			}
			result := global.DB.Table("user").Create(&user) // 通过数据的指针来创建
			fmt.Println(user)
			if result.Error == nil {
				//创建成功根据id生成token,注册的时候不存token只返回
				// 创建一个新的 JWT
				token := jwt.New(jwt.SigningMethodHS256)

				// 设置 JWT 的声明（Payload）
				claims := token.Claims.(jwt.MapClaims)
				claims["sub"] = user.Id                               // 主题
				claims["name"] = username                             // 名称
				claims["iat"] = time.Now().Unix()                     // 签发时间
				claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 过期时间

				// 设置密钥
				// 注意：在生产环境中，请使用更安全的方法存储和管理密钥
				key := []byte("my-secret-key")

				// 签名并获取完整的 JWT 字符串
				tokenString, err := token.SignedString(key)
				if err != nil {
					fmt.Println("Failed to generate JWT:", err)
					return
				}

				// 输出生成的 JWT
				fmt.Println("Generated JWT:", tokenString)
				err = global.RS.Set(strconv.Itoa(int(user.Id)), tokenString, 0).Err()
				if err != nil {
					fmt.Println("Failed to set key:", err)
					return
				}
			}
		}

	} else {
		// 用户名已存在
		fmt.Println("exit")
	}

}
func TestLogin(t *testing.T) {
	username := "郑梓桐"
	oldpassword := "20021211zzt"
	var user model.User
	result := global.DB.Table("user").Where("name = ?", username).First(&user)
	if result.Error != nil {
		// 处理查询错误
		fmt.Println("something wrong")
	}
	fmt.Println(user.Password)
	options := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(user.Password, "$")
	fmt.Println(passwordInfo)
	check := password.Verify(oldpassword, passwordInfo[2], passwordInfo[3], options)
	fmt.Println(check) // true
	if check {
		fmt.Println("login in success")
		//将token存入到redis中，并且返回到前端
		// 创建一个新的 JWT
		token := jwt.New(jwt.SigningMethodHS256)

		// 设置 JWT 的声明（Payload）
		claims := token.Claims.(jwt.MapClaims)
		claims["sub"] = user.Id                               // 主题
		claims["name"] = username                             // 名称
		claims["iat"] = time.Now().Unix()                     // 签发时间
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 过期时间

		// 设置密钥
		// 注意：在生产环境中，请使用更安全的方法存储和管理密钥
		key := []byte("my-secret-key")

		// 签名并获取完整的 JWT 字符串
		tokenString, err := token.SignedString(key)
		if err != nil {
			fmt.Println("Failed to generate JWT:", err)
			return
		}

		// 输出生成的 JWT
		fmt.Println("Generated JWT:", tokenString)
		err = global.RS.Set(strconv.Itoa(int(user.Id)), tokenString, 0).Err()
		if err != nil {
			fmt.Println("Failed to set key:", err)
			return
		}
	}

}

/*
	查询用户的详情信息
*/
func TestSelectUserDetail(t *testing.T) {
	userId := 7
	var user model.User
	global.DB.Table("user").Where("id = ?", userId).First(&user)
	fmt.Println(user)

}
