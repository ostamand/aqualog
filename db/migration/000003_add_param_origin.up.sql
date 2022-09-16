DROP INDEX param_types_user_id_name_idx;
CREATE UNIQUE INDEX ON "param_types" ("user_id", "name");

CREATE TABLE "param_origins" (
  "id" bigserial PRIMARY KEY,
  "param_type_name" varchar NOT NULL,
  "name" varchar NOT NULL,
  "created_by" varchar NOT NULL DEFAULT ('admin'),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

INSERT INTO param_origins (param_type_name, "name") VALUES
    ('phosphate', 'Hanna - Low Range Phosphate'),
    ('phosphate', 'Hanna - Ultra Low Range Phosphate'),
    ('phosphate', 'Read Sea - Phosphate Pro'),
    ('phosphate', 'Read Sea - High Accuracy Phosphate Pro'),
    ('phosphate', 'API - Phosphate'),
    ('phosphate', 'Fauna Marin - Aquahometest PO4'), 
    ('nitrate', 'Hanna - Low Range Nitrate'),
    ('nitrate', 'Hanna - High Range Nitrate'),
    ('nitrate', 'Red Sea - Nitrate Pro'),
    ('nitrate', 'Red Sea - High Accuracy Nitrate Pro'),
    ('nitrate', 'API - Nitrate'),
    ('alkalinity', 'Hanna - Alkalinity'),
    ('alkalinity', 'Red Sea - Alkalinity Pro'),
    ('alkalinity', 'Red Sea - High Accuracy Alkalinity Pro'),
    ('magnesium', 'Salifert - Magnesium'), 
    ('magnesium', 'Red Sea - High Accuracy Magnesium Pro');