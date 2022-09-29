CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT FALSE,
    "expires_at" bigint NOT NULL, 
  "created_at" bigint NOT NULL
);

ALTER TABLE sessions ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id")