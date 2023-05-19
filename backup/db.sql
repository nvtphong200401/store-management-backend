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
    store_id bigint
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
-- Name: products; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.products (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    product_name text,
    category text,
    price numeric,
    stock bigint,
    store_id bigint NOT NULL
);


ALTER TABLE public.products OWNER TO root;

--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.products_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.products_id_seq OWNER TO root;

--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


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
-- Name: products id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: store_models id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.store_models ALTER COLUMN id SET DEFAULT nextval('public.store_models_id_seq'::regclass);


--
-- Data for Name: employees; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.employees (id, created_at, updated_at, deleted_at, username, password, store_id) FROM stdin;
1	2023-05-19 09:20:39.550606+00	2023-05-19 09:20:53.900344+00	\N	phongdz	\\x24326124313024687551746f5976525734784a58465271304b4a6c50656952574974536734632f6c317579486c6265417155305755646c495a69624b	1
2	2023-05-19 09:37:34.759449+00	2023-05-19 09:38:08.908065+00	\N	phongdz2	\\x243261243130246e7154524357554377597065627a36393132534f662e51546a68554751357465537a6b647a73426d586f466e53736b575368455a61	2
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.products (id, created_at, updated_at, deleted_at, product_name, category, price, stock, store_id) FROM stdin;
2	2023-05-19 09:49:19.832094+00	2023-05-19 09:49:19.832094+00	\N	Keo cao su	Do an	5000	0	2
2	2023-05-19 09:49:37.585965+00	2023-05-19 09:49:37.585965+00	\N	Keo cao su	Do an	5000	0	1
\.


--
-- Data for Name: store_models; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.store_models (id, created_at, updated_at, deleted_at, store_name, address) FROM stdin;
1	2023-05-19 09:20:53.872464+00	2023-05-19 09:20:53.872464+00	\N	TapHoaThuyLien	Trung My Tay
2	2023-05-19 09:38:08.850875+00	2023-05-19 09:38:08.850875+00	\N	TapHoaThuyLien	Trung My Tay
\.


--
-- Name: employees_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.employees_id_seq', 2, true);


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.products_id_seq', 1, false);


--
-- Name: store_models_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.store_models_id_seq', 2, true);


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
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (store_id, id);


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

