package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var mySecret = []byte("jwt")

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}

// MyClaims 自定义声明结构体并内嵌 jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段，若需要额外记录其他字段，就可以自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中

type MyClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

func GenToken(userID int64) (string, error) {
	// 创建一个我们自己的声明的数据
	c := MyClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(10000) * time.Hour).Unix(), // 过期时间
			Issuer: "easyTikTok",
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, keyFunc)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println(token)
	// 校验token
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
