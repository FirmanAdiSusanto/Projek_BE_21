create database project;

use project;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    phone_number VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    balance DOUBLE
);

CREATE TABLE transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    phone_number VARCHAR(2055),
    amount DOUBLE,
    type VARCHAR(255)
);


select * from transactions;
select * from users;