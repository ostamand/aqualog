DROP INDEX param_types_user_id_name_idx;
CREATE UNIQUE INDEX ON "param_types" ("user_id", "name");

CREATE TABLE "param_origins" (
  "id" bigserial PRIMARY KEY,
  "param_type_name" varchar NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "param_origins" ADD FOREIGN KEY ("param_type_name") REFERENCES "param_types" ("name");