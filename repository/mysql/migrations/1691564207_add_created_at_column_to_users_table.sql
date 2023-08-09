-- +migrate Up
alter table users
    add created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;


-- +migrate Down
alter table users drop column created_at;