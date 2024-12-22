-- name: CreateChirp :one
insert into chirps (id, created_at, updated_at, body, user_id)
values (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
returning *;

-- name: GetChirps :many
select * from chirps;

-- name: GetChirp :one
select * from chirps
    where id = $1;

-- name: DeleteChirp :exec
delete from chirps
    where id = $1;

-- name: GetChirpsAuthor :many
select * from chirps
    where user_id = $1;