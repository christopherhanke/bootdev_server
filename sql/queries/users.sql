-- name: CreateUser :one
insert into users (id, created_at, updated_at, email)
values (
    gen_random_uuid(),
    now(),
    now(),
    $1
)
returning *;

-- name: DeleteUsers :exec
delete from users;

-- name: GetUser :one
select * from users where email = $1;