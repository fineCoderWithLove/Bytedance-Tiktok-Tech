/*
 Navicat Premium Data Transfer

 Source Server         : zzt_connect
 Source Server Type    : MySQL
 Source Server Version : 80032
 Source Host           : 43.143.80.216:3306
 Source Schema         : douyin

 Target Server Type    : MySQL
 Target Server Version : 80032
 File Encoding         : 65001

 Date: 01/09/2023 12:12:17
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for attention
-- ----------------------------
DROP TABLE IF EXISTS `attention`;
CREATE TABLE `attention`  (
  `user_id` int NOT NULL COMMENT '用户id',
  `to_user_id` int NOT NULL COMMENT '用户关注的人的id'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_vietnamese_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`  (
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci NOT NULL COMMENT '评论内容',
  `create_date` datetime(0) NOT NULL COMMENT '评论发布日期，格式 mm-dd',
  `comment_id` int NOT NULL COMMENT '评论id',
  `video_id` int NOT NULL COMMENT '评论所属于哪一个视频',
  `user_id` int NOT NULL COMMENT '评论属于哪一个用户',
  `deleted_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`comment_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_vietnamese_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for favorites
-- ----------------------------
DROP TABLE IF EXISTS `favorites`;
CREATE TABLE `favorites`  (
  `user_id` int NOT NULL COMMENT '用户的id',
  `video_id` int NOT NULL COMMENT '视频的id，如果都有记录说明该user赞过该视频',
  `create_time` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`, `video_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_vietnamese_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `user_id` int NOT NULL COMMENT '发送者id',
  `msg_id` int NOT NULL AUTO_INCREMENT,
  `to_user_id` int NULL DEFAULT NULL COMMENT '接收者id',
  `create_time` int NULL DEFAULT NULL COMMENT '发送消息时间',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci NULL DEFAULT NULL COMMENT '消息内容',
  PRIMARY KEY (`msg_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_vietnamese_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci NOT NULL COMMENT '用户名称',
  `follow_count` int NULL DEFAULT NULL COMMENT '关注总数',
  `follower_count` int NULL DEFAULT NULL COMMENT '粉丝数',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci NULL DEFAULT NULL COMMENT '用户头像存地址',
  `background_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci NULL DEFAULT NULL COMMENT '用户个人页顶部大图',
  `signature` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci NULL DEFAULT NULL COMMENT '个人简介',
  `total_favorited` int NULL DEFAULT NULL COMMENT '用户获赞数量',
  `work_count` int NULL DEFAULT NULL COMMENT '作品数目',
  `favorite_count` int NULL DEFAULT NULL COMMENT '点赞数量',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci NULL DEFAULT NULL COMMENT '加密后的密码',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_id`(`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 29 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_vietnamese_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos`  (
  `video_id` int NOT NULL AUTO_INCREMENT,
  `play_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci NOT NULL COMMENT '视频播放地址',
  `cover_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci NULL DEFAULT NULL COMMENT '视频封面壁纸地址',
  `favorite_count` int NULL DEFAULT NULL COMMENT '视频点赞数目',
  `comment_count` int NULL DEFAULT NULL COMMENT '视频评论数目',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_vietnamese_ci NULL DEFAULT NULL COMMENT '视频标题',
  `user_id` int NOT NULL COMMENT '视频作者的唯一标识',
  `created_time` int NOT NULL COMMENT '视频发布的时间',
  PRIMARY KEY (`video_id`) USING BTREE,
  INDEX `fk_videos_user`(`user_id`) USING BTREE,
  CONSTRAINT `fk_videos_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 32 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_vietnamese_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
