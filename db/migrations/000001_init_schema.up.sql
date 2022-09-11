CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(50) NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar(10) NOT NULL,
  "created_at" bigint NOT NULL
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" bigint NOT NULL
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "from_entry_id" bigint NOT NULL,
  "to_entry_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" bigint NOT NULL
);


CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transactions" ("from_account_id");

CREATE INDEX ON "transactions" ("to_account_id");

CREATE INDEX ON "transactions" ("from_account_id", "to_account_id");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("to_entry_id") REFERENCES "entries" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("from_entry_id") REFERENCES "entries" ("id");


