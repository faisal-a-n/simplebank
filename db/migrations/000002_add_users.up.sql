CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar NOT NULL UNIQUE,
  "password_changed_at" bigint NOT NULL,
  "created_at" bigint NOT NULL
);

ALTER TABLE "accounts" ADD user_id bigint NOT NULL;

CREATE INDEX ON "accounts" ("user_id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

-- CREATE UNIQUE INDEX ON "accounts" ("user_id", "currency")
ALTER TABLE "accounts" ADD CONSTRAINT "user_currency_key" UNIQUE ("user_id", "currency");
