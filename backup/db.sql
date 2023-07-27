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
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    product_name text,
    category text,
    price numeric,
    stock bigint,
    store_id bigint NOT NULL,
    price_in numeric,
    price_out numeric
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
-- Name: sale_items; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.sale_items (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    sale_id bigint,
    product_id bigint,
    quantity bigint,
    stock bigint
);


ALTER TABLE public.sale_items OWNER TO root;

--
-- Name: sale_items_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.sale_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sale_items_id_seq OWNER TO root;

--
-- Name: sale_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.sale_items_id_seq OWNED BY public.sale_items.id;


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
-- Name: products id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: sale_items id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.sale_items ALTER COLUMN id SET DEFAULT nextval('public.sale_items_id_seq'::regclass);


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
1	2023-05-23 06:24:00.17537+00	2023-05-23 06:24:04.42433+00	\N	phongdz	\\x243261243130247630497748765245756132476e645252617a666c4b6533482e796b496d385943387a464a72794b4d386f58524c37465654684c5132	1	owner
2	2023-05-23 08:48:31.518098+00	2023-05-24 03:19:03.467027+00	\N	phongdz1	\\x243261243130246b39584c34494351544648656a38556e73542e4d592e705a434b48685a413930693057495a436f787a3373477749574d6464665a4b	1	staff
\.


--
-- Data for Name: join_requests; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.join_requests (employee_id, store_id, status) FROM stdin;
2	1	accepted
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.products (id, created_at, updated_at, deleted_at, product_name, category, price, stock, store_id, price_in, price_out) FROM stdin;
2	2023-07-25 10:13:52.188476+00	2023-07-25 10:13:52.188476+00	\N	Mi Hao Hao	Mi an lien	\N	0	1	5000	10000
1	0001-01-01 00:00:00+00	2023-07-25 10:14:12.624275+00	\N	Mi Hao Hao chua cay	Mi an lien	10000	10	1	5000	10000
3	2023-07-26 04:58:21.179383+00	2023-07-26 04:58:21.179383+00	\N	Mi Hao Hao	Mi an lien	\N	5	1	5000	10000
4	2023-07-26 06:16:51.80351+00	2023-07-26 06:16:51.80351+00	\N	Phong test		\N	4	1	1000	10000
5	2023-07-26 09:57:28.685967+00	2023-07-26 09:57:28.685967+00	\N	Mi Hao Hao	Mi an lien	\N	5	1	5000	10000
6	2023-07-26 09:57:47.013506+00	2023-07-26 09:57:47.013506+00	\N	Mi Hao Hao	Mi an lien	\N	5	1	5000	10000
7	2023-07-26 09:57:47.013506+00	2023-07-26 09:57:47.013506+00	\N	Mi Hao Hao	Mi an lien	\N	5	1	5000	10000
8	2023-07-26 15:19:32.723486+00	2023-07-26 15:19:32.723486+00	\N	Test add		\N	10	1	1200	12000
9	2023-07-26 15:33:11.276284+00	2023-07-26 15:33:11.276284+00	\N	test add multi		\N	2	1	100	1000
10	2023-07-26 15:33:11.276284+00	2023-07-26 15:33:11.276284+00	\N	test add multi pro		\N	24	1	100	1000
11	2023-07-26 15:34:31.859157+00	2023-07-26 15:34:31.859157+00	\N	asd		\N	2	1	123	1234
12	2023-07-26 15:37:02.607006+00	2023-07-26 15:37:02.607006+00	\N	asd		\N	12	1	10	1000
14	2023-07-26 15:38:13.450678+00	2023-07-26 15:38:13.450678+00	\N	asd		\N	11	1	100	10000
123	2023-07-26 15:41:09.18206+00	2023-07-26 15:41:09.18206+00	\N	asfsad		\N	1	1	1252	12213
\.


--
-- Data for Name: sale_items; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.sale_items (id, created_at, updated_at, deleted_at, sale_id, product_id, quantity, stock) FROM stdin;
1	2023-05-23 07:16:35.810871+00	2023-05-23 07:16:35.810871+00	\N	3	1	2	\N
2	2023-05-23 07:26:27.58943+00	2023-05-23 07:26:27.58943+00	\N	4	1	2	\N
3	2023-07-26 07:16:40.164527+00	2023-07-26 07:16:40.164527+00	\N	5	1	2	\N
4	2023-07-27 03:41:06.282949+00	2023-07-27 03:41:06.282949+00	\N	8	1	\N	0
5	2023-07-27 03:44:40.45725+00	2023-07-27 03:44:40.45725+00	\N	9	1	\N	0
6	2023-07-27 03:44:54.737257+00	2023-07-27 03:44:54.737257+00	\N	10	1	\N	2
7	2023-07-27 04:01:47.56328+00	2023-07-27 04:01:47.56328+00	\N	16	1	\N	2
8	2023-07-27 04:02:06.314618+00	2023-07-27 04:02:06.314618+00	\N	17	1	\N	10
9	2023-07-27 04:04:00.770363+00	2023-07-27 04:04:00.770363+00	\N	18	1	\N	10
10	2023-07-27 04:04:28.02983+00	2023-07-27 04:04:28.02983+00	\N	19	1	\N	10
11	2023-07-27 04:08:36.20644+00	2023-07-27 04:08:36.20644+00	\N	20	1	\N	10
12	2023-07-27 04:09:08.956261+00	2023-07-27 04:09:08.956261+00	\N	21	1	\N	10
\.


--
-- Data for Name: sale_models; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.sale_models (id, created_at, updated_at, deleted_at, store_id, employee_id, total_price) FROM stdin;
1	2023-05-23 07:14:25.335785+00	2023-05-23 07:14:25.335785+00	\N	1	1	10000
2	2023-05-23 07:14:42.730001+00	2023-05-23 07:14:42.730001+00	\N	1	1	10000
3	2023-05-23 07:16:35.71541+00	2023-05-23 07:16:35.71541+00	\N	1	1	10000
4	2023-05-23 07:26:27.477665+00	2023-05-23 07:26:27.477665+00	\N	1	1	10000
5	2023-07-26 07:16:40.047546+00	2023-07-26 07:16:40.047546+00	\N	1	1	20000
6	2023-07-27 03:33:04.924067+00	2023-07-27 03:33:04.924067+00	\N	1	1	0
7	2023-07-27 03:37:47.003339+00	2023-07-27 03:37:47.003339+00	\N	1	1	0
8	2023-07-27 03:41:06.189521+00	2023-07-27 03:41:06.189521+00	\N	1	1	0
9	2023-07-27 03:44:40.396536+00	2023-07-27 03:44:40.396536+00	\N	1	1	0
10	2023-07-27 03:44:54.676458+00	2023-07-27 03:44:54.676458+00	\N	1	1	20000
11	2023-07-27 03:50:36.678987+00	2023-07-27 03:50:36.678987+00	\N	1	1	0
12	2023-07-27 03:52:14.581493+00	2023-07-27 03:52:14.581493+00	\N	1	1	0
13	2023-07-27 03:54:20.053672+00	2023-07-27 03:54:20.053672+00	\N	1	1	0
14	2023-07-27 03:56:12.931209+00	2023-07-27 03:56:12.931209+00	\N	1	1	0
15	2023-07-27 04:00:47.822944+00	2023-07-27 04:00:47.822944+00	\N	1	1	100000
16	2023-07-27 04:01:47.441051+00	2023-07-27 04:01:47.441051+00	\N	1	1	20000
17	2023-07-27 04:02:06.203778+00	2023-07-27 04:02:06.203778+00	\N	1	1	100000
18	2023-07-27 04:04:00.699347+00	2023-07-27 04:04:00.699347+00	\N	1	1	100000
19	2023-07-27 04:04:27.88161+00	2023-07-27 04:04:27.88161+00	\N	1	1	100000
20	2023-07-27 04:08:36.057664+00	2023-07-27 04:08:36.057664+00	\N	1	1	100000
21	2023-07-27 04:09:08.738048+00	2023-07-27 04:09:08.738048+00	\N	1	1	100000
\.


--
-- Data for Name: store_models; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.store_models (id, created_at, updated_at, deleted_at, store_name, address) FROM stdin;
1	2023-05-23 06:24:04.384954+00	2023-05-23 06:24:04.384954+00	\N	TapHoaThuyLien	Trung My Tay
\.


--
-- Name: employees_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.employees_id_seq', 2, true);


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.products_id_seq', 2, true);


--
-- Name: sale_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.sale_items_id_seq', 12, true);


--
-- Name: sale_models_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.sale_models_id_seq', 21, true);


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
    ADD CONSTRAINT products_pkey PRIMARY KEY (store_id, id);


--
-- Name: sale_items sale_items_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_pkey PRIMARY KEY (id);


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
-- Name: products_name_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX products_name_idx ON public.products USING gin (product_name public.gin_trgm_ops);


--
-- Name: products_name_trigram_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX products_name_trigram_idx ON public.products USING gin (product_name public.gin_trgm_ops);


--
-- Name: products_product_name_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX products_product_name_idx ON public.products USING gin (product_name public.gin_trgm_ops);


--
-- Name: join_requests fk_join_requests_employee; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.join_requests
    ADD CONSTRAINT fk_join_requests_employee FOREIGN KEY (employee_id) REFERENCES public.employees(id);


--
-- Name: join_requests fk_join_requests_store; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.join_requests
    ADD CONSTRAINT fk_join_requests_store FOREIGN KEY (store_id) REFERENCES public.store_models(id);


--
-- Name: sale_items fk_sale_items_sale; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT fk_sale_items_sale FOREIGN KEY (sale_id) REFERENCES public.sale_models(id);


--
-- Name: sale_models fk_sale_models_employee; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.sale_models
    ADD CONSTRAINT fk_sale_models_employee FOREIGN KEY (employee_id) REFERENCES public.employees(id);


--
-- Name: sale_models fk_sale_models_store; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.sale_models
    ADD CONSTRAINT fk_sale_models_store FOREIGN KEY (store_id) REFERENCES public.store_models(id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: root
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;


--
-- PostgreSQL database dump complete
--

