-- name: CreateUser :one
insert into users (id, created_at, updated_at, email)
values (
    gen_randum_uuid(),
    now(),
    now(),
    $1
)
returning *;
