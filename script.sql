create user if not exists 'schedule'@'localhost' identified by 'schedule';
GRANT ALL PRIVILEGES ON *.* TO 'schedule'@'localhost' with grant option;
flush privileges;

create database if not exists schedule;

use schedule;

create table if not exists students (
    studentId int primary key auto_increment not null,
    studentName varchar(255) not null,
    spareTime JSON,
    status bool
);

create table if not exists teachers (
    teacherId int primary key auto_increment not null,
    teacherName varchar(255) not null,
    spareTime JSON,
    holidayNum int,
    status bool
);

create table if not exists classes (
    classId int primary key auto_increment not null,
    className varchar(255) not null,
    classMates JSON,
    status bool
);

create table if not exists lessons (
    lessonId int primary key auto_increment not null,
    lessonName varchar(255) not null,
    teacherName varchar(255),
    studyName varchar(255),
    studentNum int
);
