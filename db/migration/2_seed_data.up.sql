INSERT INTO "authors" ("name") VALUES ('Anonymous'), ('Unknown'), ('Author');

INSERT INTO "categories" ("label") VALUES ('funny'), ('inspirational'), ('life'), ('love'), ('philosophy'), ('success'), ('wisdom');

INSERT INTO "quotes" ("content", "author_id", "category_id") VALUES ('My first quote', 1, 1), ('My second quote', 2, 2), ('My third quote', 3, 3);
