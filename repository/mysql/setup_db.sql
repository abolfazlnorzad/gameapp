CREATE table users(
                      id int primary key AUTO_INCREMENT,
                      name varchar(255) not null,
                      password varchar(255) not null,
                      phone_number varchar(255) not null unique,
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);