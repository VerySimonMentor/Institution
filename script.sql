create user if not exists "institution"@"localhost" identified by "institution";
GRANT ALL PRIVILEGES ON *.* TO 'institution'@'localhost' with grant option;
flush privileges;

create database if not exists institution;

use institution;

create table if not exists country(
    countryId int primary key auto_increment not null,
    countryEngName varchar(255) not null,
    countryChiName varchar(255) not null,
    countryAndSchool JSON
);

create table if not exists school(
    schoolId int primary key auto_increment not null,
    schoolEngName varchar(255) not null,
    schoolChiName varchar(255) not null,
    schoolAbbreviation varchar(255),
    schoolType varchar(255) not null,
    province varchar(255) not null,
    officialWebLink varchar(255) not null,
    schoolRemark varchar(4095),
    schoolAndItem JSON
);

create table if not exists item(
    itemId int primary key auto_increment not null,
    itemName varchar(255) not null,
    externalDisplay JSON,
    itemRemark varchar(4095)
);

create table if not exists user(
    userId int primary key auto_increment not null,
    userAccount varchar(31) not null,
    userPasswd varchar(31),
    userEmail varchar(255),
    userNumber varchar(255) not null,
    userLevel int not null,
    studentCount int
);
