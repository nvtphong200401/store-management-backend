CREATE SEQUENCE IF NOT EXISTS products_id_seq

-- Table Definition
CREATE TABLE "public"."products" (
    "id" int4 NOT NULL DEFAULT nextval('products_id_seq'::regclass),
    "product_name" text NOT NULL,
    "category" text NOT NULL,
    "price" numeric,
    "created_at" timestamp NOT NULL DEFAULT '2023-04-09 08:15:00.723576'::timestamp without time zone,
    "updated_at" timestamp NOT NULL,
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

INSERT INTO "public"."products" ("id","product_name","category","price","created_at","updated_at","deleted_at") VALUES ('2','test','cate','11','0001-01-01 00:00:00','2023-04-09 16:49:25.436464',NULL);
INSERT INTO "public"."products" ("id","product_name","category","price","created_at","updated_at","deleted_at") VALUES ('1','test','cate','10','2023-04-09 15:37:36.498656','2023-04-09 16:41:30.982763','2023-04-09 20:16:09.581921');
