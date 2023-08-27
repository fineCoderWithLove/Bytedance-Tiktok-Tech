create database if not exists mini_douyin
    character set utf8mb4
    collate utf8mb4_general_ci;

use mini_douyin;

create table if not exists mini_douyin.tb_comment
(
    comment_id  bigint auto_increment comment '评论id'
        primary key,
    user_id     bigint                             not null comment '用户id',
    video_id    bigint                             not null comment '视频id',
    content     varchar(255)                       not null comment '评论内容',
    create_date date                               not null comment '评论发布日期',
    create_time datetime default CURRENT_TIMESTAMP null comment '创建时间',
    update_time datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间'
)
    comment '评论表';

create table if not exists mini_douyin.tb_like
(
    like_id     bigint auto_increment comment '点赞id'
        primary key,
    user_id     bigint                             not null comment '用户id',
    video_id    bigint                             not null comment '视频id',
    create_time datetime default CURRENT_TIMESTAMP null comment '创建时间',
    update_time datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间',
    constraint idx_userId_videoId
        unique (user_id, video_id)
)
    comment '点赞表';

create table if not exists mini_douyin.tb_message
(
    message_id  bigint auto_increment comment '消息id'
        primary key,
    user_id     bigint                             not null comment '用户id',
    to_user_id  bigint                             not null comment '对方用户id',
    content     varchar(255)                       not null comment '消息内容',
    action_type char                               not null comment '消息类型',
    create_time datetime default CURRENT_TIMESTAMP null comment '创建时间',
    update_time datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间'
)
    comment '关注表';

create table if not exists mini_douyin.tb_relation
(
    relation_id bigint auto_increment comment '关注id'
        primary key,
    user_id     bigint                             not null comment '用户id',
    to_user_id  bigint                             not null comment '被关注者id',
    create_time datetime default CURRENT_TIMESTAMP null comment '创建时间',
    update_time datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间',
    constraint idx_userId_toUserId
        unique (user_id, to_user_id)
)
    comment '关注表';

create table if not exists mini_douyin.tb_user
(
    user_id          bigint auto_increment comment '用户id'
        primary key,
    username         varchar(32)                             not null comment '用户名',
    password         varchar(32)                             not null comment '密码',
    avatar           varchar(255)                            null comment '用户头像',
    follow_count     bigint       default 0                  null comment '关注总数',
    follower_count   bigint       default 0                  null comment '粉丝总数',
    work_count       bigint       default 0                  null comment '作品数',
    favorite_count   bigint       default 0                  null comment '喜欢数',
    total_favorited  bigint       default 0                  null comment '获赞数量 ',
    signature        varchar(255) default '此人没有填写简介' null comment '个人简介',
    background_image varchar(255) default ''                 null comment '用户个人页顶部大图',
    create_time      datetime     default CURRENT_TIMESTAMP  null comment '创建时间',
    update_time      datetime     default CURRENT_TIMESTAMP  null on update CURRENT_TIMESTAMP comment '更新时间',
    constraint idx_username
        unique (username)
)
    comment '用户表';

create table if not exists mini_douyin.tb_video
(
    video_id       bigint auto_increment comment '视频id'
        primary key,
    user_id        bigint                             not null comment '用户id',
    title          varchar(255)                       not null comment '视频标题',
    play_url       varchar(512)                       not null comment '播放地址',
    cover_url      varchar(512)                       not null comment '视频封面地址',
    comment_count  bigint   default 0                 not null comment '视频的评论总数',
    create_time    datetime default CURRENT_TIMESTAMP null comment '创建时间',
    update_time    datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间',
    favorite_count bigint   default 0                 null comment '点赞数量'
)
    comment '视频表';


