CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "values" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "value_type_id" int NOT NULL,
  "value" float NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "value_types" (
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

CREATE INDEX ON "values" ("user_id");

CREATE INDEX ON "values" ("user_id", "value_type_id");

CREATE INDEX ON "value_types" ("user_id");

CREATE INDEX ON "value_types" ("user_id", "name");

ALTER TABLE "values" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "values" ADD FOREIGN KEY ("value_type_id") REFERENCES "value_types" ("id");

ALTER TABLE "value_types" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
