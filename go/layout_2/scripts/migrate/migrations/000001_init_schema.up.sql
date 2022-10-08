create table student
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)            null,
    updated_at datetime(3)            null,
    deleted_at datetime(3)            null,
    name       varchar(10) default '' not null comment '姓名',
    gender     tinyint(1)  default 0  not null comment '性別 0:男, 1:女',
    status     tinyint(1)  default 1  not null comment '狀態 0:禁用, 1:啟用'
)
    comment '學生表';

create index idx_name
    on student (name);

create index idx_student_deleted_at
    on student (deleted_at);

