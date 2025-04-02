create table if not exists users(
    id serial primary key,
    name text not null,
    surname text not null,
    patronymic text not null,
    age int not null,
    gender text not null,
    nationality text not null
);