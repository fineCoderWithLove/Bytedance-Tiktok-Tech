package db

import (
	"douyin/douyin-api/global"
	"douyin/douyin-api/social-service/proto"
	"github.com/jinzhu/gorm"
)

type IUserRepository interface {
	GetUserById(id int64) (proto.User, error)                        // 获取单个用户
	CreateUser(userName string, password string) (proto.User, error) // 创建新用户
	CheckUser(userName string, password string) (proto.User, error)  // 检查是否存在用户
	QueryUserIdByUserName(name string) (int64, error)                // 根据用户名称查询用户ID
	UpdateWorkCount(userId int64, isAdd bool)                        // 更新作品数量
	UpdateFollowCount(userId int64, isAdd bool)                      // 更新关注数量
	UpdateFollowerCount(userId int64, isAdd bool)                    // 更新粉丝数量
	UpdateFavoriteCount(userId int64, isAdd bool)                    // 更新喜欢数量
	UpdateTotalFavoritedCount(userId int64, isAdd bool)              // 增加获赞总数

}

// UserRepository 定义一个结构体
type UserRepository struct {
}

func (ur UserRepository) UpdateWorkCount(userId int64, isAdd bool) {
	db := global.DB.Table("tb_user").Where("user_id = ?", userId)

	if isAdd {
		db.UpdateColumn("work_count", gorm.Expr("work_count + ?", 1))
	}
	db.UpdateColumn("work_count", gorm.Expr("work_count - ?", 1))
}

func (ur UserRepository) UpdateFollowCount(userId int64, isAdd bool) {
	db := global.DB.Table("tb_user").Where("user_id = ?", userId)

	if isAdd {
		db.UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1))
	}
	db.UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1))
}

func (ur UserRepository) UpdateFollowerCount(userId int64, isAdd bool) {
	db := global.DB.Table("tb_user").Where("user_id = ?", userId)

	if isAdd {
		db.UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1))
	}
	db.UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1))
}

func (ur UserRepository) UpdateFavoriteCount(userId int64, isAdd bool) {
	db := global.DB.Table("tb_user").Where("user_id = ?", userId)

	if isAdd {
		db.UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))
	}
	db.UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1))
}

// todo total_favorited 在数据库中存储的是 varchar 考虑改成 BigInt 在返回的时候转化为 string
func (ur UserRepository) UpdateTotalFavoritedCount(userId int64, isAdd bool) {
	db := global.DB.Table("tb_user").Where("user_id = ?", userId)

	if isAdd {
		db.UpdateColumn("total_favorited", gorm.Expr("total_favorited + ?", 1))
	}
	db.UpdateColumn("total_favorited", gorm.Expr("total_favorited - ?", 1))
}

func (ur UserRepository) QueryUserIdByUserName(name string) (int64, error) {
	var user proto.User
	global.DB.Table("tb_user").Select("user_id").Where("username = ?", name).Scan(&user)
	return user.UserId, nil
}

// CheckUser 检查用户是否存在
func (ur UserRepository) CheckUser(userName string, password string) (proto.User, error) {
	var user proto.User
	err := global.DB.Table("tb_user").Where("username = ? and password = ?", userName, password).Scan(&user).Error
	return user, err
}

// GetUserById 获取单个用户
func (ur UserRepository) GetUserById(id int64) (proto.User, error) {
	var user proto.User
	err := global.DB.Where("user_id = ?", id).First(&user).Error
	return user, err
}

// CreateUser 创建新用户
func (ur UserRepository) CreateUser(userName string, password string) (proto.User, error) {
	var user proto.User

	// 填写默认值
	user.BackgroundImage = "D:/aaaaaaSpace/Go/goTikTok/social-service/utils/download.jpg"
	user.Avatar = "D:/aaaaaaSpace/Go/goTikTok/social-service/utils/download.jpg"
	user.UserName = userName
	user.Password = password
	err := global.DB.Table("tb_user").Create(&user).Error
	return user, err
}

// NewUserRepository UserRepository构造函数
func NewUserRepository() IUserRepository {
	return UserRepository{}
}
