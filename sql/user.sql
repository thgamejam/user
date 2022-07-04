create schema touhou_gamejam;

use touhou_gamejam;

create table user (
    id          int unsigned        auto_increment  primary key,

    ctime       datetime            not null                comment '创建时间',
    mtime       datetime            not null                comment '修改时间'
);
