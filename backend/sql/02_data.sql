BEGIN TRANSACTION;

INSERT INTO "money2" (id, pair_id, type_id, user_id, amount) VALUES(1, 1, 1, 1, 0);
-- TODO:sign up時に初期金額を登録するようにすれば、↑はいらないかも

INSERT INTO "types" VALUES(1,'収入');
INSERT INTO "types" VALUES(2,'合計支出');
INSERT INTO "types" VALUES(3,'支出');

COMMIT;