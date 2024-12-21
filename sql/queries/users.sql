-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, hashed_password)
values (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
returning *;

-- name: DeleteUsers :exec
delete from users;

-- name: GetUser :one
select * from users where email = $1;

-- name: UpdateUser :one
update users
    set email = $2,
        hashed_password = $3,
        updated_at = now()
    where id = $1
returning *;

-- name: UpgradeUser :one
update users
    set is_chirpy_red = true
    where id = $1
returning *;
