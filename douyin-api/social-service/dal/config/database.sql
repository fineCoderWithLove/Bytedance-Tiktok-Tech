create database if not exists douyin
    character set utf8mb4
    collate utf8mb4_general_ci;

use douyin;


create table if not exists douyin.message
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

create table if not exists douyin.relation
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


