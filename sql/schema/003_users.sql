-- +goose Up
alter table users 
    add column hashed_password text not null default 'unset',
    alter column hashed_password set default '';

-- +goose Down
alter table users drop column hashed_password;