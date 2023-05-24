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
-- Name: sale_items; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.sale_items (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    sale_id bigint,
    product_id bigint,
    quantity bigint
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

COPY public.products (id, created_at, updated_at, deleted_at, product_name, category, price, stock, store_id) FROM stdin;
1	2023-05-23 06:31:43.828697+00	2023-05-23 06:31:43.828697+00	\N	Mi Hao Hao	Mi an lien	5000	0	1
\.


--
-- Data for Name: sale_items; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.sale_items (id, created_at, updated_at, deleted_at, sale_id, product_id, quantity) FROM stdin;
1	2023-05-23 07:16:35.810871+00	2023-05-23 07:16:35.810871+00	\N	3	1	2
2	2023-05-23 07:26:27.58943+00	2023-05-23 07:26:27.58943+00	\N	4	1	2
\.


--
-- Data for Name: sale_models; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.sale_models (id, created_at, updated_at, deleted_at, store_id, employee_id, total_price) FROM stdin;
1	2023-05-23 07:14:25.335785+00	2023-05-23 07:14:25.335785+00	\N	1	1	10000
2	2023-05-23 07:14:42.730001+00	2023-05-23 07:14:42.730001+00	\N	1	1	10000
3	2023-05-23 07:16:35.71541+00	2023-05-23 07:16:35.71541+00	\N	1	1	10000
4	2023-05-23 07:26:27.477665+00	2023-05-23 07:26:27.477665+00	\N	1	1	10000
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

SELECT pg_catalog.setval('public.products_id_seq', 1, true);


--
-- Name: sale_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.sale_items_id_seq', 2, true);


--
-- Name: sale_models_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.sale_models_id_seq', 4, true);


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

