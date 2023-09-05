package util

import (
	"demotest/social-service/global"
	"demotest/social-service/proto/user"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"strconv"
)

// Video 结构体定义，与数据库表映射
type Video struct {
	VideoID       int64
	FavoriteCount int64
	CommentCount  int64
}

var client *redis.Client

const (
	VideoRedisPrefix = "video:"
	UserRedisPrefix  = "user:"
)

func init() {
	client = global.RS
}

// InitRedisData 初始化操作：查询数据库并将数据存入 Redis
func InitRedisData() {
	var videos []Video
	global.DB.Find(&videos)

	for _, video := range videos {
		key := VideoRedisPrefix + fmt.Sprintf("%d", video.VideoID)
		data := map[string]interface{}{
			"favorite":  strconv.FormatInt(video.FavoriteCount, 10), // 转为字符串存储
			"comment":   strconv.FormatInt(video.CommentCount, 10),  // 转为字符串存储
			"liked":     "false",                                    // 初始化点赞状态为 false
			"commented": "false",                                    // 初始化评论状态为 false
		}
		client.HMSet(key, data)
	}

	var users []pb.User
	global.DB.Table("user").Find(&users)
	for _, user := range users {
		key := UserRedisPrefix + fmt.Sprintf("%d", user.Id)
		data := map[string]interface{}{
			"totalFavorited": strconv.FormatInt(user.TotalFavorited, 10),
			"favoriteCount":  strconv.FormatInt(user.FavoriteCount, 10),
			"followerCount":  strconv.FormatInt(user.FollowerCount, 10),
			"followCount":    strconv.FormatInt(user.FollowCount, 10),
			"statusChanged":  "false",
		}
		client.HMSet(key, data)
	}
}

// WriteRedisVideoToMySQL 将 Redis 数据写入 MySQL 数据库
func WriteRedisVideoToMySQL() {
	keys, err := client.Keys(VideoRedisPrefix + "*").Result()
	if err != nil {
		fmt.Println("Error getting keys from Redis:", err)
		return
	}

	for _, key := range keys {
		videoIDStr := key[len(VideoRedisPrefix):]
		videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
		if err != nil {
			fmt.Println("Error converting videoID to int:", err)
			continue
		}

		data, err := client.HGetAll(key).Result()
		if err != nil {
			fmt.Println("Error getting data from Redis:", err)
			continue
		}

		liked := data["liked"]
		commented := data["commented"]

		// 仅在点赞数或评论数有变化时更新 MySQL 数据库
		if liked == "true" || commented == "true" {
			// 从 Redis 数据中获取点赞数和评论数，转为 int64 类型
			favoriteCount, _ := strconv.ParseInt(data["favorite"], 10, 64)
			commentCount, _ := strconv.ParseInt(data["comment"], 10, 64)
			fmt.Println(favoriteCount, commentCount, videoID)
			updateData := map[string]interface{}{
				"favorite_count": favoriteCount,
				"comment_count":  commentCount,
			}
			global.DB.Model(&Video{}).Where("video_id = ?", videoID).Updates(updateData)
			fmt.Println("更新数据库")
			//将状态设置为false
			if _, err = client.HSet(key, "liked", "false").Result(); err != nil {
				zap.S().Infof("设置点赞状态异常")
			}
			if _, err := client.HSet(key, "commented", "false").Result(); err != nil {
				zap.S().Infof("设置评论状态异常")
			}
		} else {
		}
	}
}

// WriteRedisUserToMySQL 将 Redis 数据写入 MySQL 数据库
func WriteRedisUserToMySQL() {
	keys, err := client.Keys(UserRedisPrefix + "*").Result()
	if err != nil {
		fmt.Println("Error getting keys from Redis:", err)
		return
	}

	for _, key := range keys {
		userIDStr := key[len(UserRedisPrefix):]
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			fmt.Println("Error converting userID to int:", err)
			continue
		}

		data, err := client.HGetAll(key).Result()
		if err != nil {
			fmt.Println("Error getting data from Redis:", err)
			continue
		}
		statusChanged := data["statusChanged"]

		//仅在数据有变化才更新Mysql
		if statusChanged == "true" {
			totalFavorited, favoriteCount, followerCount, followCount, err := GetUserAllUserData(userID)

			updateData := map[string]interface{}{
				"total_favorited": totalFavorited,
				"favorite_count":  favoriteCount,
				"follower_count":  followerCount,
				"follow_count":    followCount,
			}
			if err != nil {
				fmt.Println("Error getting user data:", err)
				continue
			}
			fmt.Println(totalFavorited, favoriteCount, followerCount, followCount, userID)
			global.DB.Table("user").Model(&pb.User{}).Where("id = ?", userID).Updates(updateData)

			if _, err := client.HSet(key, "statusChanged", "false").Result(); err != nil {
				zap.S().Infof("设置用户状态异常")
			}

		}
	}

}

// LikeVideo 方法2：点赞操作
func LikeVideo(videoID int64) error {
	key := VideoRedisPrefix + fmt.Sprintf("%d", videoID)

	// 更新点赞状态
	_, err := client.HSet(key, "liked", "true").Result()
	if err != nil {
		if err != nil {
			//如果redis没有则从数据获取并初始化这条视频
			_, _, err = GetVideoFavoriteAndCommentCount(videoID)
			if err != nil {
				return err
			}
		}
	}

	// 增加点赞数
	newCount, err := client.HIncrBy(key, "favorite", 1).Result()
	if err != nil {
		return err
	}

	fmt.Printf("Liked video %d. New favorite count: %d\n", videoID, newCount)
	return nil
}

// UnlikeVideo 方法3：取消点赞操作
func UnlikeVideo(videoID int64) error {
	key := VideoRedisPrefix + fmt.Sprintf("%d", videoID)

	// 更新点赞状态
	_, err := client.HSet(key, "liked", "true").Result()
	if err != nil {
		return err
	}

	// 减少点赞数
	newCount, err := client.HIncrBy(key, "favorite", -1).Result()
	if err != nil {
		return err
	}

	fmt.Printf("Unliked video %d. New favorite count: %d\n", videoID, newCount)
	return nil
}

// AddComment 方法4：评论数+1操作
func AddComment(videoID int64) error {
	key := VideoRedisPrefix + fmt.Sprintf("%d", videoID)

	// 更新评论状态
	_, err := client.HSet(key, "commented", "true").Result()
	if err != nil {
		return err
	}
	// 增加评论数
	newCount, err := client.HIncrBy(key, "comment", 1).Result()
	if err != nil {
		return err
	}

	fmt.Printf("Added comment to video %d. New comment count: %d\n", videoID, newCount)
	return nil
}

func DelComment(videoID int64) error {
	key := VideoRedisPrefix + fmt.Sprintf("%d", videoID)

	// 更新评论状态
	_, err := client.HSet(key, "commented", "true").Result()
	if err != nil {
		return err
	}

	// 减少评论数
	newCount, err := client.HIncrBy(key, "comment", -1).Result()
	if err != nil {
		return err
	}

	fmt.Printf("Unliked video %d. New favorite count: %d\n", videoID, newCount)
	return nil
}

// GetVideoFavoriteAndCommentCount 获取视频点赞数和评论数，如果 Redis 中不存在则从 MySQL 查询并初始化
func GetVideoFavoriteAndCommentCount(videoID int64) (int64, int64, error) {
	key := VideoRedisPrefix + fmt.Sprintf("%d", videoID)

	// 尝试从 Redis 获取数据
	data, err := client.HGetAll(key).Result()
	if len(data) == 0 || err != nil {
		fmt.Println("视频不存在", videoID)
		// Redis 中不存在，从 MySQL 查询数据
		var video *Video
		result := global.DB.First(video, "video_id = ?", videoID)
		if result.Error != nil && !(result.RowsAffected > 0) {
			return 0, 0, errors.New("视频不存在")
		}
		//初始化 Redis 数据
		fields := map[string]interface{}{
			"favorite":  strconv.FormatInt(video.FavoriteCount, 10),
			"comment":   strconv.FormatInt(video.CommentCount, 10),
			"liked":     "false", // 初始化点赞状态为 false
			"commented": "false", // 初始化评论状态为 false
		}
		_, err := client.HMSet(key, fields).Result()
		if err != nil {
			return 0, 0, err
		}

		return video.FavoriteCount, video.CommentCount, nil
	}
	// 从 Redis 获取点赞数和评论数
	favoriteCount, _ := strconv.ParseInt(data["favorite"], 10, 64)
	commentCount, _ := strconv.ParseInt(data["comment"], 10, 64)
	return favoriteCount, commentCount, nil
}

func GetUserAllUserData(userID int64) (int64, int64, int64, int64, error) {

	key := UserRedisPrefix + fmt.Sprintf("%d", userID)
	data, err := client.HGetAll(key).Result()
	if len(data) == 0 || err != nil {
		// redis不存在从数据库中获取用户数据
		var user *pb.User
		result := global.DB.Table("user").First(&user, "id = ?", userID)
		if result.Error != nil && !(result.RowsAffected > 0) {
			return 0, 0, 0, 0, errors.New("用户不存在")
		}
		// 存入 Redis
		userData := map[string]interface{}{
			"totalFavorited": strconv.FormatInt(user.TotalFavorited, 10),
			"favoriteCount":  strconv.FormatInt(user.FavoriteCount, 10),
			"followerCount":  strconv.FormatInt(user.FollowerCount, 10),
			"followCount":    strconv.FormatInt(user.FollowCount, 10),
			"statusChanged":  "false",
		}
		_, err := client.HMSet(key, userData).Result()
		if err != nil {
			return 0, 0, 0, 0, err
		}
		return user.TotalFavorited, user.FavoriteCount, user.FollowerCount, user.FollowCount, nil

	}

	// 从 Redis 数据中获取用户信息
	totalFavorited, _ := strconv.ParseInt(data["totalFavorited"], 10, 64)
	favoriteCount, _ := strconv.ParseInt(data["favoriteCount"], 10, 64)
	followerCount, _ := strconv.ParseInt(data["followerCount"], 10, 64)
	followCount, _ := strconv.ParseInt(data["followCount"], 10, 64)
	return totalFavorited, favoriteCount, followerCount, followCount, nil
}

// IncreaseTotalFavorited 获赞总数+1操作
func IncreaseTotalFavorited(userID int64) error {
	key := UserRedisPrefix + fmt.Sprintf("%d", userID)

	// 更新 totalFavorited + 1，同时更新 statusChanged 为 true
	_, err := client.HIncrBy(key, "totalFavorited", 1).Result()
	if err != nil {
		return err
	}

	_, err = client.HSet(key, "statusChanged", "true").Result()
	if err != nil {
		return err
	}

	fmt.Printf("Increased total favorited for user %d and changed status to true\n", userID)

	return nil
}

// DecreaseTotalFavorited 获赞总数-1操作
func DecreaseTotalFavorited(userID int64) error {
	key := UserRedisPrefix + fmt.Sprintf("%d", userID)

	// 更新 totalFavorited - 1，同时更新 statusChanged 为 true
	_, err := client.HIncrBy(key, "totalFavorited", -1).Result()
	if err != nil {
		return err
	}

	_, err = client.HSet(key, "statusChanged", "true").Result()
	if err != nil {
		return err
	}

	fmt.Printf("Decreased total favorited for user %d and changed status to true\n", userID)

	return nil
}

// IncreaseFavoriteCount 用户点赞总数+1操作
func IncreaseFavoriteCount(userID int64) error {
	key := UserRedisPrefix + fmt.Sprintf("%d", userID)

	// 更新 favoriteCount + 1，同时更新 statusChanged 为 true
	_, err := client.HIncrBy(key, "favoriteCount", 1).Result()
	if err != nil {
		return err
	}

	_, err = client.HSet(key, "statusChanged", "true").Result()
	if err != nil {
		return err
	}

	fmt.Printf("Increased favorited count for user %d and changed status to true\n", userID)

	return nil
}

// DecreaseFavoriteCount 用户点赞总数-1操作
func DecreaseFavoriteCount(userID int64) error {
	key := UserRedisPrefix + fmt.Sprintf("%d", userID)

	// 更新 totalFavorited - 1，同时更新 statusChanged 为 true
	_, err := client.HIncrBy(key, "favoriteCount", -1).Result()
	if err != nil {
		return err
	}

	_, err = client.HSet(key, "statusChanged", "true").Result()
	if err != nil {
		return err
	}

	fmt.Printf("Decreased favorited count for user %d and changed status to true\n", userID)

	return nil
}

// IncreaseFollowerCount 粉丝总数+1操作
func IncreaseFollowerCount(userID int64) error {
	key := UserRedisPrefix + fmt.Sprintf("%d", userID)

	// 更新 favoriteCount + 1，同时更新 statusChanged 为 true
	_, err := client.HIncrBy(key, "followerCount", 1).Result()
	if err != nil {
		return err
	}

	_, err = client.HSet(key, "statusChanged", "true").Result()
	if err != nil {
		return err
	}

	fmt.Printf("Increased follower count for user %d and changed status to true\n", userID)

	return nil
}

// DecreaseFollowerCount 粉丝总数-1操作
func DecreaseFollowerCount(userID int64) error {
	key := UserRedisPrefix + fmt.Sprintf("%d", userID)

	// 更新 followerCount - 1，同时更新 statusChanged 为 true
	_, err := client.HIncrBy(key, "followerCount", -1).Result()
	if err != nil {
		return err
	}

	_, err = client.HSet(key, "statusChanged", "true").Result()
	if err != nil {
		return err
	}

	fmt.Printf("Decreased follower count for user %d and changed status to true\n", userID)

	return nil
}

// IncreaseFollowCount 用户关注总数+1操作
func IncreaseFollowCount(userID int64) error {
	key := UserRedisPrefix + fmt.Sprintf("%d", userID)

	// 更新 followCount + 1，同时更新 statusChanged 为 true
	_, err := client.HIncrBy(key, "followCount", 1).Result()
	if err != nil {
		return err
	}

	_, err = client.HSet(key, "statusChanged", "true").Result()
	if err != nil {
		return err
	}

	fmt.Printf("Increased follow count for user %d and changed status to true\n", userID)

	return nil
}

// DecreaseFollowCount 用户关注总数-1操作
func DecreaseFollowCount(client *redis.Client, userID int64) error {
	key := UserRedisPrefix + fmt.Sprintf("%d", userID)

	// 更新 followCount - 1，同时更新 statusChanged 为 true
	_, err := client.HIncrBy(key, "followCount", -1).Result()
	if err != nil {
		return err
	}

	_, err = client.HSet(key, "statusChanged", "true").Result()
	if err != nil {
		return err
	}

	fmt.Printf("Decreased follow count for user %d and changed status to true\n", userID)

	return nil
}
