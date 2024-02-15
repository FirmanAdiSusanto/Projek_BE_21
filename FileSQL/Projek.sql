create database project;

use project;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    phone_number VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    balance DOUBLE
);


CREATE TABLE transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    user_id INT,
    phone_number VARCHAR(255),
    amount DOUBLE,
    type VARCHAR(255)
);


select * from transactions;
select * from users;