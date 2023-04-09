-- Sequences
CREATE SEQUENCE IF NOT EXISTS users_id_seq

-- Table Definition
CREATE TABLE "public"."users" (
    "id" int8 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "username" text,
    "password" bytea,
    PRIMARY KEY ("id")
);

INSERT INTO "public"."users" ("id","created_at","updated_at","deleted_at","username","password") VALUES ('1','2023-04-09 14:30:48.787399+00','2023-04-09 14:30:48.787399+00',NULL,'phong','\x24326124313024683153305931645148706a32594e666c776c375846655635473952506379314f7337473139682e3938646d615459684551785a6a36');
