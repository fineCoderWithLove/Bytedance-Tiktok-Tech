package handler

import (
	"context"
	"crypto/md5"
	"crypto/sha512"
	"demotest/base-service/global"
	"demotest/base-service/model"
	pb "demotest/base-service/proto"
	"demotest/base-service/util"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"math/rand"
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
1.实现用户注册的接口，注册的时候会自动登录所以也要存token
*/
func (s *UserServe) UserRegister(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterOrLoginInfoResp, error) {
	//取出信息
	username := req.Username
	password := req.Password
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
			// 头像地址列表
			avatarURLs := []string{
				"https://th.bing.com/th?id=OIP.sYQ-V3EviYCtpoSDBEuj4gAAAA&w=250&h=250&c=8&rs=1&qlt=90&o=6&dpr=1.5&pid=3.1&rm=2",
				"https://th.bing.com/th?id=OIP.VGCWVwVU92WcwhgZCbJNhAAAAA&w=250&h=250&c=8&rs=1&qlt=90&o=6&dpr=1.5&pid=3.1&rm=2",
				"https://th.bing.com/th?id=OIP.ZtrvDWpDrJNiyRQj6DwPvwAAAA&w=250&h=250&c=8&rs=1&qlt=90&o=6&dpr=1.5&pid=3.1&rm=2",
				"https://th.bing.com/th?id=OIP.NfRGPXQCHEqoRQsBRqqHngAAAA&w=250&h=250&c=8&rs=1&qlt=90&o=6&dpr=1.5&pid=3.1&rm=2",
				"https://th.bing.com/th/id/OIP.hY5BzKbBo4xGW8XGUi9EagAAAA?w=180&h=180&c=7&r=0&o=5&dpr=1.5&pid=1.7",
				"https://th.bing.com/th/id/OIP.LYRZ_nETF4W4VchbeGIP5gAAAA?w=210&h=209&c=7&r=0&o=5&dpr=1.5&pid=1.7",
				"https://th.bing.com/th/id/OIP.jPmwChP95v4N0kaBQJuAQQAAAA?w=209&h=209&c=7&r=0&o=5&dpr=1.5&pid=1.7",
				// 添加更多头像地址
			}

			// 背景地址列表
			backgroundURLs := []string{
				"https://pic1.zhimg.com/80/v2-73b8307b2db44c617f4e8515ce67dd39_1440w.webp?source=1940ef5c",
				"https://picx.zhimg.com/80/v2-e5427c1e9ad8aaad99d643e7bd7e927b_1440w.webp?source=1940ef5c",
				"https://picx.zhimg.com/80/v2-d024c6ad6851b266e8509d1aa0948ceb_1440w.webp?source=1940ef5c",
				"https://picx.zhimg.com/80/v2-904505bcf0c424788f6028b8952aa2e7_1440w.webp?source=1940ef5c",
				"https://pica.zhimg.com/80/v2-6665ba00a7f4bf46cf76d9169be1c9e1_1440w.webp?source=1940ef5c",
				// 添加更多背景地址
			}

			// 签名列表
			signatureList := []string{
				"抖音走出中国面向世界",
				"姐抽的是烟，它伤肺。但不伤心",
				"你若折断她半边翅膀，我便毁了你整个天堂",
				"你的酒窝没有酒，我却醉得像条狗",
				"心里有座坟，埋着未亡人",
				// 添加更多签名
			}

			// 生成随机数种子
			rand.Seed(time.Now().UnixNano())

			// 生成随机索引
			avatarIndex := rand.Intn(len(avatarURLs))
			backgroundIndex := rand.Intn(len(backgroundURLs))
			signatureIndex := rand.Intn(len(signatureList))

			// 获取对应的地址和签名
			avatarURL := avatarURLs[avatarIndex]
			backgroundURL := backgroundURLs[backgroundIndex]
			signature := signatureList[signatureIndex]
			user := model.User{
				Id:              0,
				Name:            username,
				Password:        newpassword,
				Avatar:          avatarURL,
				BackgroundImage: backgroundURL,
				Signature:       signature,
			}
			result := global.DB.Table("user").Create(&user) // 通过数据的指针来创建
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
					return nil, nil
				}
				err = global.RS.Set(tokenString, strconv.Itoa(int(user.Id)), 0).Err()
				if err != nil {
				}
				//并且把生成的token返回
				resp := &pb.RegisterOrLoginInfoResp{
					StatusCode: 0,
					StatusMsg:  "注册成功",
					UserId:     user.Id,
					Token:      tokenString,
				}
				return resp, nil
			}
		}

	} else {
		// 用户名已存在
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
	options := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(user.Password, "$")
	check := password.Verify(oldpassword, passwordInfo[2], passwordInfo[3], options)
	if check {
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
	uid := global.RS.Get(req.Token).Val()
	Intuid, _ := strconv.Atoi(uid)
	if req.Token != "" {
		userId := req.UserId
		var user model.User
		result := global.DB.Table("user").Where("id = ?", userId).First(&user)
		if result.Error != nil {
			resp := &pb.UserDetailResp{
				StatusCode: 403,
				StatusMsg:  "您尚未注册",
				User:       nil,
			}
			return resp, nil
		} else {
			var Isfollow bool
			var attentions []Attention
			result := global.DB.Table("attention").Where("user_id = ? AND to_user_id = ?", Intuid, req.UserId).Find(&attentions)
			if result.RowsAffected != 0 {
				Isfollow = true
			}
			//查询视频流的时候要从redis中查询出用户的所有信息
			TotalFavorited, FavoriteCount, FollowerCount, FollowCount, _ := util.GetUserAllUserData(int64(Intuid))
			pbUser := &pb.User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   FollowCount,
				FollowerCount: FollowerCount,
				//TODO 需要另外进行查询
				IsFollow:        Isfollow,
				Avatar:          user.Avatar,
				BackgroundImage: user.BackgroundImage,
				Signature:       user.Signature,
				TotalFavorited:  TotalFavorited,
				WorkCount:       user.WorkCount,
				FavoriteCount:   FavoriteCount,
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
