CREATE TABLE IF NOT EXISTS money2
(
    id                  integer primary key autoincrement,
    type_id             integer,
    user_id             integer,
    amount              integer,
    created_at  text NOT NULL DEFAULT (DATETIME('now', 'localtime'))
    -- updated_at  text NOT NULL DEFAULT (DATETIME('now', 'localtime'))

);

CREATE TABLE IF NOT EXISTS users
(
    id       integer primary key autoincrement,
    name     varchar(50),
    password binary(60),
    balance integer,
    calculation  integer

);

CREATE TABLE IF NOT EXISTS types
(
    id   integer primary key,
    name varchar(50)
);