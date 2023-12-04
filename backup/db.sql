--
-- PostgreSQL database dump
--

-- Dumped from database version 15.3
-- Dumped by pg_dump version 15.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: root
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO root;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: root
--

COMMENT ON SCHEMA public IS '';


--
-- Name: pg_trgm; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;


--
-- Name: EXTENSION pg_trgm; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';


--
-- Name: unaccent; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS unaccent WITH SCHEMA public;


--
-- Name: EXTENSION unaccent; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION unaccent IS 'text search dictionary that removes accents';


--
-- Name: fr; Type: TEXT SEARCH CONFIGURATION; Schema: public; Owner: root
--

CREATE TEXT SEARCH CONFIGURATION public.fr (
    PARSER = pg_catalog."default" );

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR asciiword WITH french_stem;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR word WITH public.unaccent, french_stem;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR numword WITH simple;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR email WITH simple;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR url WITH simple;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR host WITH simple;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR sfloat WITH simple;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR version WITH simple;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR hword_numpart WITH simple;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR hword_part WITH public.unaccent, french_stem;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR hword_asciipart WITH french_stem;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR numhword WITH simple;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR asciihword WITH french_stem;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR hword WITH public.unaccent, french_stem;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR url_path WITH simple;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR file WITH simple;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR "float" WITH simple;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR "int" WITH simple;

ALTER TEXT SEARCH CONFIGURATION public.fr
    ADD MAPPING FOR uint WITH simple;


ALTER TEXT SEARCH CONFIGURATION public.fr OWNER TO root;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: employees; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.employees (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    username text,
    password bytea,
    store_id bigint,
    "position" text DEFAULT 'unknown'::text
);


ALTER TABLE public.employees OWNER TO root;

--
-- Name: employees_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.employees_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.employees_id_seq OWNER TO root;

--
-- Name: employees_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.employees_id_seq OWNED BY public.employees.id;


--
-- Name: join_requests; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.join_requests (
    employee_id bigint NOT NULL,
    store_id bigint NOT NULL,
    status text DEFAULT 'pending'::text
);


ALTER TABLE public.join_requests OWNER TO root;

--
-- Name: products; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.products (
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    product_name text,
    category text,
    price_in numeric,
    price_out numeric,
    stock bigint,
    store_id bigint NOT NULL
);


ALTER TABLE public.products OWNER TO root;

--
-- Name: sale_items; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.sale_items (
    sale_id bigint NOT NULL,
    product_id text NOT NULL,
    stock bigint,
    deleted_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.sale_items OWNER TO root;

--
-- Name: sale_models; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.sale_models (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    store_id bigint,
    employee_id bigint,
    total_price numeric
);


ALTER TABLE public.sale_models OWNER TO root;

--
-- Name: sale_models_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.sale_models_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sale_models_id_seq OWNER TO root;

--
-- Name: sale_models_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.sale_models_id_seq OWNED BY public.sale_models.id;


--
-- Name: store_models; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.store_models (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    store_name text,
    address text
);


ALTER TABLE public.store_models OWNER TO root;

--
-- Name: store_models_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.store_models_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.store_models_id_seq OWNER TO root;

--
-- Name: store_models_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.store_models_id_seq OWNED BY public.store_models.id;


--
-- Name: employees id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.employees ALTER COLUMN id SET DEFAULT nextval('public.employees_id_seq'::regclass);


--
-- Name: sale_models id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.sale_models ALTER COLUMN id SET DEFAULT nextval('public.sale_models_id_seq'::regclass);


--
-- Name: store_models id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.store_models ALTER COLUMN id SET DEFAULT nextval('public.store_models_id_seq'::regclass);


--
-- Data for Name: employees; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.employees (id, created_at, updated_at, deleted_at, username, password, store_id, "position") FROM stdin;
1	2023-08-03 08:33:01.771677+00	2023-08-03 08:51:56.237636+00	\N	phongdz	\\x2432612431302439564a5144625031635078705746506c6c4b6a7a326564784e594861676f70503779416931775058354f723252475848504547344b	1	owner
\.


--
-- Data for Name: join_requests; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.join_requests (employee_id, store_id, status) FROM stdin;
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.products (id, created_at, updated_at, deleted_at, product_name, category, price_in, price_out, stock, store_id) FROM stdin;
	2023-09-14 14:53:47.981066+00	2023-09-14 14:53:47.981066+00	\N	djj		0	11000	20	1
123	2023-09-15 13:30:44.406313+00	2023-09-15 13:30:44.406313+00	\N	phongdz		0	12000	20	1
8936136164151	2023-09-15 15:32:28.193831+00	2023-09-15 15:33:51.43683+00	\N	chinsudayne		0	300000	20	1
147	2023-09-15 15:34:54.496094+00	2023-09-15 15:34:54.558769+00	\N	trsst		0	0.8	20	1
\.


--
-- Data for Name: sale_items; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.sale_items (sale_id, product_id, stock, deleted_at, created_at, updated_at) FROM stdin;
1	123	20	\N	2023-09-15 14:29:48.452257+00	2023-09-15 14:29:48.452257+00
2	8936136164151	20	\N	2023-09-15 15:32:46.403485+00	2023-09-15 15:32:46.403485+00
\.


--
-- Data for Name: sale_models; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.sale_models (id, created_at, updated_at, deleted_at, store_id, employee_id, total_price) FROM stdin;
1	2023-09-15 14:29:48.444782+00	2023-09-15 14:29:48.444782+00	\N	1	1	240000
2	2023-09-15 15:32:46.401748+00	2023-09-15 15:32:46.401748+00	\N	1	1	6000000
\.


--
-- Data for Name: store_models; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.store_models (id, created_at, updated_at, deleted_at, store_name, address) FROM stdin;
1	2023-08-03 08:51:56.19731+00	2023-08-03 08:51:56.19731+00	\N	TapHoaThuyLien	Trung My Tay
\.


--
-- Name: employees_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.employees_id_seq', 1, true);


--
-- Name: sale_models_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.sale_models_id_seq', 2, true);


--
-- Name: store_models_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.store_models_id_seq', 1, true);


--
-- Name: employees employees_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_pkey PRIMARY KEY (id);


--
-- Name: employees employees_username_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_username_key UNIQUE (username);


--
-- Name: join_requests join_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.join_requests
    ADD CONSTRAINT join_requests_pkey PRIMARY KEY (employee_id, store_id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id, store_id);


--
-- Name: sale_items sale_items_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_pkey PRIMARY KEY (sale_id, product_id);


--
-- Name: sale_models sale_models_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.sale_models
    ADD CONSTRAINT sale_models_pkey PRIMARY KEY (id);


--
-- Name: store_models store_models_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.store_models
    ADD CONSTRAINT store_models_pkey PRIMARY KEY (id);


--
-- Name: idx_employees_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_employees_deleted_at ON public.employees USING btree (deleted_at);


--
-- Name: idx_products_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_products_deleted_at ON public.products USING btree (deleted_at);


--
-- Name: idx_sale_items_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_sale_items_deleted_at ON public.sale_items USING btree (deleted_at);


--
-- Name: idx_sale_models_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_sale_models_deleted_at ON public.sale_models USING btree (deleted_at);


--
-- Name: idx_store_models_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_store_models_deleted_at ON public.store_models USING btree (deleted_at);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: root
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;


--
-- PostgreSQL database dump complete
--

