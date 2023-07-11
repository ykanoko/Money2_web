CREATE TABLE IF NOT EXISTS money2
(
    id                  integer primary key autoincrement,
    pair_id integer,
    type_id             integer,
    user_id             integer,
    amount              integer,
    created_at  text NOT NULL DEFAULT (DATETIME('now', 'localtime'))
    -- updated_at  text NOT NULL DEFAULT (DATETIME('now', 'localtime'))

);

CREATE TABLE IF NOT EXISTS types
(
    id   integer primary key,
    name varchar(50)
);

CREATE TABLE IF NOT EXISTS pairs
(
    id       integer primary key autoincrement,
    password binary(60),
    user1_id integer,
    user2_id integer,
    calculation_user1  decimal,
    created_at  text NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TABLE IF NOT EXISTS users
(
    id       integer primary key autoincrement,
    name     varchar(50),
    balance decimal
);