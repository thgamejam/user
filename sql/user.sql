use touhou_gamejam;

create table user
(
    id                int unsigned auto_increment primary key,
    account_id        int unsigned     not null comment '账户id索引',
    name              varchar(16)      not null default '' comment '名称',
    avatar_id         int unsigned     not null default 0 comment '头像id',
    bio               varchar(64)      not null default '' comment '个人简介',
    display_tag1      tinyint unsigned not null default 0 comment '展示的标签1',
    display_tag2      tinyint unsigned not null default 0 comment '展示的标签2',
    display_tag3      tinyint unsigned not null default 0 comment '展示的标签3',
    allow_syndication boolean          not null default true comment '是否允许联合发布邀请',
    ctime             datetime         not null comment '创建时间',
    mtime             datetime         not null comment '修改时间'
);

create unique index idx_user_account ON user (account_id);

create table user_tag_relational
(
    id          int unsigned auto_increment primary key,
    user_id     int unsigned     not null comment '用户id',
    user_tag_id tinyint unsigned not null default 0 comment '用户标签索引',
    status      tinyint unsigned not null default 0 comment '标签状态',
    ctime       datetime         not null comment '创建时间',
    mtime       datetime         not null comment '修改时间'
);

create index user_tag_relational_user_id_index on user_tag_relational (user_id);

create table user_tag_enum
(
    id      tinyint unsigned auto_increment primary key,
    content varchar(8) not null default '' comment '标签内容',
    ctime   datetime   not null comment '创建时间',
    mtime   datetime   not null comment '修改时间'
);
