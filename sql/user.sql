use touhou_gamejam;

create table user (
    id                int unsigned auto_increment primary key,
    account_id        int unsigned     not null comment '账户id索引',
    uname             varchar(16)      not null comment '用户名',
    avatar_id         int unsigned     not null default 0 comment '头像id',
    bio               varchar(64)      not null default '' comment '个人简介',
    display_tag1      tinyint unsigned not null default 0 comment '展示的标签1',
    display_tag2      tinyint unsigned not null default 0 comment '展示的标签2',
    display_tag3      tinyint unsigned not null default 0 comment '展示的标签3',
    allow_syndication boolean          not null default true comment '是否允许联合发布邀请',
    status            tinyint unsigned not null default 1 comment '用户状态',
    ctime             datetime         not null comment '创建时间',
    mtime             datetime         not null comment '修改时间'
);

alter table user
    auto_increment = 100000;

create unique index idx_user_account ON user (account_id);


create table user_tag_relational (
    id          int unsigned auto_increment primary key,
    user_id     int unsigned     not null comment '用户id',
    user_tag_id tinyint unsigned not null default 0 comment '用户标签索引',
    del         bool             not null default false comment '删除状态',
    ctime       datetime         not null comment '创建时间',
    mtime       datetime         not null comment '修改时间'
);

create index user_tag_relational_user_id_index on user_tag_relational (user_id);


create table user_tag_enum (
    id      tinyint unsigned auto_increment primary key,
    content varchar(8) not null default '' comment '标签内容',
    ctime   datetime   not null comment '创建时间',
    mtime   datetime   not null comment '修改时间'
);

alter table user_tag_enum
    auto_increment = 1;

create unique index idx_user_tag_enum_content ON user_tag_enum (content);


create table user_relationship (
    id            int unsigned auto_increment primary key,
    user_id       int      not null comment '用户id',
    follow_userid int      not null comment '被关注用户id',
    ctime         datetime not null comment '创建时间',
    mtime         datetime not null comment '修改时间'
);

create index user_relationship_user_id_index on user_relationship (user_id);


create table user_follow_info (
    id           int unsigned auto_increment primary key,
    user_id      int      not null comment '用户id',
    fans_count   int      not null default 0 comment '粉丝数量',
    follow_count int      not null default 0 comment '关注数量',
    ctime        datetime not null comment '创建时间',
    mtime        datetime not null comment '修改时间'
);

create index user_follow_info_user_id_index on user_follow_info (user_id);
