CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "email" varchar NOT NULL,
  "admin" boolean NOT NULL DEFAULT (false),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "params" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "param_type_id" bigint NOT NULL,
  "value" float NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "param_types" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar,
  "unit" varchar NOT NULL,
  "user_id" bigint NOT NULL,
  "target" float,
  "min" float,
  "max" float,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "params" ("user_id");

CREATE INDEX ON "params" ("user_id", "param_type_id");

CREATE INDEX ON "param_types" ("user_id");

CREATE INDEX ON "param_types" ("user_id", "name");

ALTER TABLE "params" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "params" ADD FOREIGN KEY ("param_type_id") REFERENCES "param_types" ("id");

ALTER TABLE "param_types" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
