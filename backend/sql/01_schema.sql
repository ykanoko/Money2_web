CREATE TABLE IF NOT EXISTS money2
(
    id                  integer primary key autoincrement,
    date                text NOT NULL DEFAULT (DATETIME('now', 'localtime')),
    type_id             integer,
    user_id             integer,
    amount              integer,
    money_user1         integer,
    money_user2         integer,
    calculation_user1   integer
);

CREATE TABLE IF NOT EXISTS users
(
    id       integer primary key autoincrement,
    name     varchar(50),
    password binary(60)
);

CREATE TABLE IF NOT EXISTS types
(
    id   integer primary key,
    name varchar(50)
);

BEGIN TRANSACTION;

INSERT INTO "money2" VALUES(0,DATETIME('now', 'localtime'), 0, 0, 0, 0, 0, 0);
-- TODO:sign up時に初期金額を登録するようにすれば、↑はいらないかも

INSERT INTO "types" VALUES(1,'収入');
INSERT INTO "types" VALUES(2,'合計支出');
INSERT INTO "types" VALUES(3,'支出');

COMMIT;