create database Account_service_app;

use Account_service_app;

CREATE TABLE user(
id int primary key auto_increment,
nama varchar(50) not null,
No_telepon varchar(20) not null unique,
password varchar(50) not null,
tanggal_lahir varchar(20) not null,
created_at datetime default current_timestamp,
updated_at datetime,
deleted_at datetime
);

CREATE TABLE top_Up(
id int primary key auto_increment,
user_id int,
Jumlah_Topup decimal,
created_at datetime default current_timestamp,
updated_at datetime,
deleted_at datetime,
constraint fk_user_top_up FOREIGN KEY (user_id) REFERENCES user(id)
);

CREATE TABLE Saldo(
user_id int,
Jumlah_Saldo decimal,
created_at datetime default current_timestamp,
updated_at datetime,
deleted_at datetime,
constraint fk_user_saldo FOREIGN KEY (user_id) REFERENCES user(id)
);

CREATE TABLE Transfer(
id int primary key auto_increment,
user_id int,
nama_penerima varchar(50),
Nomor_telefon Varchar(20),
Jumlah_Transfer decimal,
created_at datetime default current_timestamp,
updated_at datetime,
deleted_at datetime,
constraint fk_user_Transfer FOREIGN KEY (user_id) REFERENCES user(id)
);

CREATE TABLE TransferPenerima(
id int primary key auto_increment,
user_id int,
nama_pengirim varchar(50),
Nomor_telefon Varchar(20),
Jumlah_Transfer decimal,
created_at datetime default current_timestamp,
updated_at datetime,
deleted_at datetime,
constraint fk_user_TransferPenerima FOREIGN KEY (user_id) REFERENCES user(id)
);

select * from Saldo;
select * from Top_Up;
select * from user;
select * from Transfer;
select * from transfer;
SELECT * FROM user WHERE No_telepon;
select * from Saldo; 
