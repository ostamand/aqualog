DROP TABLE IF EXISTS param_origins;

DROP INDEX param_types_user_id_name_idx;
CREATE INDEX ON "param_types" ("user_id", "name");