package db

import (
	"douyin/douyin-api/global"
	"douyin/douyin-api/social-service/proto"
	"log"
)

type IRelationRepository interface {
	AddFollow(userId int64, toUserId int64) error       // 添加关注记录
	GetFollow(userId int64) ([]proto.Relation, error)   // 获取关注者
	GetFollower(userId int64) ([]proto.Relation, error) // 获取关注者
	CheckIsFollow(userId int64, toUserId int64) bool    // 检查是否关注
	RemoveFollow(userId int64, toUserId int64) error    //取消关注
}

type RelationRepository struct {
}

func (r RelationRepository) RemoveFollow(userId int64, toUserId int64) error {
	var relation proto.Relation
	err := global.DB.Table("relation").Where("user_id = ? and to_user_id = ?", userId, toUserId).Delete(&relation).Error
	return err
}

// CheckIsFollow 检查是否被关注
func (r RelationRepository) CheckIsFollow(userId int64, toUserId int64) bool {
	var count int64
	err := global.DB.Table("relation").Where("user_id = ? and to_user_id = ?", userId, toUserId).Count(&count).Error
	if err != nil {
		log.Printf("CheckIsFollow|数据库获取数量错误|%v", err)
		return false
	}
	return count == 1
}

func (r RelationRepository) GetFollower(userId int64) ([]proto.Relation, error) {
	var relations []proto.Relation
	err := global.DB.Table("relation").Where("to_user_id = ?", userId).Find(&relations).Error
	return relations, err
}

func (r RelationRepository) GetFollow(userId int64) ([]proto.Relation, error) {
	var relations []proto.Relation
	err := global.DB.Table("relation").Where("user_id = ?", userId).Find(&relations).Error
	return relations, err
}

func (r RelationRepository) AddFollow(userId int64, toUserId int64) error {
	relationDomain := proto.Relation{ToUserId: toUserId, UserId: userId}
	err := global.DB.Table("relation").Create(&relationDomain).Error
	return err
}

func NewRelationRepository() IRelationRepository {
	relationRepository := RelationRepository{}
	return relationRepository
}
