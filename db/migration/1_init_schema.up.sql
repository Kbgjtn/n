CREATE TABLE IF NOT EXISTS "quotes" (
  "id" bigserial PRIMARY KEY,
  "content" varchar NOT NULL,
  "author_id" bigint NOT NULL,
  "category_id" bigint NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS  "authors" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS "categories" (
  "id" bigserial PRIMARY KEY,
  "label" varchar UNIQUE NOT NULL
);

CREATE INDEX ON "quotes" ("author_id");

CREATE INDEX ON "authors" ("name");

CREATE INDEX ON "categories" ("label");

COMMENT ON TABLE "authors" IS 'The authors can have many quotes';

ALTER TABLE "quotes" ADD FOREIGN KEY ("author_id") REFERENCES "authors" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "quotes" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
