package handler

import (
	"context"
	"crypto/md5"
	"crypto/sha512"
	"douyin/base-service/global"
	"douyin/base-service/model"
	pb "douyin/base-service/proto"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
)

// 定义md5加密
func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

type UserServe struct {
	pb.UnimplementedUserServiceServer
}

/*
	要实现三个接口的实例
	UserRegister(context.Context, *RegisterReq) (*RegisterOrLoginInfoResp, error)
	UserLogin(context.Context, *LoginReq) (*RegisterOrLoginInfoResp, error)
	UserDetail(context.Context, *DetailRep) (*UserDetailResp,	 error)
*/
//密码加密的函数
func EncodingPassword(oldpassword string) string {

	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(oldpassword, options)
	newpassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	return newpassword
}

/*
1.实现用户注册的接口
*/
func (s *UserServe) UserRegister(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterOrLoginInfoResp, error) {
	//取出信息
	fmt.Println("服务端接收的值")
	username := req.Username
	fmt.Println(username)
	password := req.Password
	fmt.Println(password)
	//加密后的密码
	newpassword := EncodingPassword(password)
	//查询用户名是否存在
	var user model.User
	result := global.DB.Table("user").Where("name = ?", username).First(&user)
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
				claims["sub"] = user.Id                                    // 主题
				claims["name"] = username                                  // 名称
				claims["iat"] = time.Now().Unix()                          // 签发时间
				claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix() // 过期时间

				// 设置密钥
				// 注意：在生产环境中，请使用更安全的方法存储和管理密钥
				key := []byte("my-secret-key")

				// 签名并获取完整的 JWT 字符串
				tokenString, err := token.SignedString(key)
				if err != nil {
					fmt.Println("Failed to generate JWT:", err)
					return nil, nil
				}

				// 输出生成的 JWT
				fmt.Println("Generated JWT:", tokenString)
				//并且把生成的token返回
				resp := &pb.RegisterOrLoginInfoResp{
					StatusCode: 200,
					StatusMsg:  "注册成功",
					UserId:     user.Id,
					Token:      tokenString,
				}
				return resp, nil
			}
		}

	} else {
		// 用户名已存在
		fmt.Println("exit")
		resp := &pb.RegisterOrLoginInfoResp{
			StatusCode: 500,
			StatusMsg:  "用户名已经存在",
			UserId:     0,
			Token:      "",
		}
		return resp, nil
	}

	return nil, nil
}

/*
2.实现用户登录的接口
*/
func (s *UserServe) UserLogin(ctx context.Context, req *pb.LoginReq) (*pb.RegisterOrLoginInfoResp, error) {
	/*
		1.先通过用户名查询加密后的密码
		2.然后利用函数解密
		3.通过传来的密码进行比对
	*/
	username := req.Username
	oldpassword := req.Password
	var user model.User
	fmt.Println("用户姓名")
	fmt.Println(username)

	result := global.DB.Table("user").Where("name = ?", username).First(&user)
	if result.Error != nil {
		// 处理查询错误
		zap.S().Error("用户尚未注册")
		resp := &pb.RegisterOrLoginInfoResp{
			StatusCode: 500,
			StatusMsg:  "用户尚未注册",
			UserId:     0,
			Token:      "",
		}
		return resp, nil
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
		}

		// 输出生成的 JWT
		fmt.Println("Generated JWT:", tokenString)
		//存入redis中
		err = global.RS.Set(tokenString, strconv.Itoa(int(user.Id)), 0).Err()
		if err != nil {
			fmt.Println("Failed to set key:", err)
		}
		zap.S().Info("登录成功")
		resp := &pb.RegisterOrLoginInfoResp{
			StatusCode: 0,
			StatusMsg:  "登录成功",
			UserId:     user.Id,
			Token:      tokenString,
		}
		return resp, nil
	} else {
		zap.S().Error("密码不正确")
		resp := &pb.RegisterOrLoginInfoResp{
			StatusCode: 500,
			StatusMsg:  "请检查用户名或者密码",
			UserId:     0,
			Token:      "",
		}
		return resp, nil
	}

}

/*
3.实现查询用户详情的接口
*/
func (s *UserServe) UserDetail(ctx context.Context, req *pb.DetailRep) (*pb.UserDetailResp, error) {
	//需要有token才能查询
	fmt.Println("进入了UserDetail")
	if req.Token != "" {
		userId := req.UserId
		fmt.Println(userId)
		var user model.User
		result := global.DB.Table("user").Where("id = ?", userId).First(&user)
		if result.Error != nil {
			zap.S().Error("未注册")
			resp := &pb.UserDetailResp{
				StatusCode: 500,
				StatusMsg:  "您尚未注册",
				User:       nil,
			}
			return resp, nil
		} else {
			pbUser := &pb.User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				//TODO 需要另外进行查询
				IsFollow:        false,
				Avatar:          user.Avatar,
				BackgroundImage: user.BackgroundImage,
				Signature:       user.Signature,
				TotalFavorited:  user.TotalFavorited,
				WorkCount:       user.WorkCount,
				FavoriteCount:   user.FavoriteCount,
			}
			resp := &pb.UserDetailResp{
				StatusCode: 0,
				StatusMsg:  "成功",
				User:       pbUser,
			}
			return resp, nil
		}

	} else {
		resp := &pb.UserDetailResp{
			StatusCode: 500,
			StatusMsg:  "登录已经过期,请重试",
			User:       nil,
		}
		return resp, nil
	}
}
