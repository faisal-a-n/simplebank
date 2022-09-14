ALTER TABLE "accounts"
DROP CONSTRAINT "user_currency_key";

ALTER TABLE "accounts" DROP user_id;

drop table if exists users;