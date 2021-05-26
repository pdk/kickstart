-- 2021-05-25-create-thing.sql

create table if not exists things (
    id          integer primary key autoincrement,
    name        varchar not null unique,
    value       varchar not null
);
