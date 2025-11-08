--
-- PostgreSQL database dump
--

\restrict HZJt83l80Aozg7JzSKKGhIVIixGpi0DHKDFyqox4tcR2hm8cKaxiPbGPAYROK4n

-- Dumped from database version 18.0 (Ubuntu 18.0-1.pgdg24.04+3)
-- Dumped by pg_dump version 18.0 (Ubuntu 18.0-1.pgdg24.04+3)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: category; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.category AS ENUM (
    'fruits',
    'vegetables',
    'grains',
    'livestock',
    'dairy',
    'others'
);


ALTER TYPE public.category OWNER TO postgres;

--
-- Name: order_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.order_status AS ENUM (
    'Packaging',
    'Shipping',
    'Recieved'
);


ALTER TYPE public.order_status OWNER TO postgres;

--
-- Name: temp_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.temp_status AS ENUM (
    'Packaging',
    'Shipping',
    'Recieved'
);


ALTER TYPE public.temp_status OWNER TO postgres;

--
-- Name: trade_bid_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.trade_bid_status AS ENUM (
    'pending',
    'accepted',
    'rejected'
);


ALTER TYPE public.trade_bid_status OWNER TO postgres;

--
-- Name: trade_listing_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.trade_listing_status AS ENUM (
    'open',
    'completed',
    'cancelled'
);


ALTER TYPE public.trade_listing_status OWNER TO postgres;

--
-- Name: transaction_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.transaction_status AS ENUM (
    'pending',
    'processing',
    'completed',
    'failed',
    'cancelled'
);


ALTER TYPE public.transaction_status OWNER TO postgres;

--
-- Name: transaction_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.transaction_type AS ENUM (
    'deposit',
    'withdrawal',
    'order_payment',
    'refund',
    'transfer'
);


ALTER TYPE public.transaction_type OWNER TO postgres;

--
-- Name: user_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.user_type AS ENUM (
    'Consumer',
    'Farmer'
);


ALTER TYPE public.user_type OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: Cart; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Cart" (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    grand_total numeric(10,2),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public."Cart" OWNER TO postgres;

--
-- Name: CartItem; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."CartItem" (
    id uuid NOT NULL,
    cart_id uuid NOT NULL,
    item_id uuid NOT NULL,
    quantity integer NOT NULL
);


ALTER TABLE public."CartItem" OWNER TO postgres;

--
-- Name: Comment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Comment" (
    id uuid CONSTRAINT comment_id_not_null NOT NULL,
    user_id uuid CONSTRAINT comment_user_id_not_null NOT NULL,
    post_id uuid CONSTRAINT comment_post_id_not_null NOT NULL,
    content text CONSTRAINT comment_content_not_null NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public."Comment" OWNER TO postgres;

--
-- Name: Farm; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Farm" (
    id uuid CONSTRAINT farm_id_not_null NOT NULL,
    farmer_id uuid CONSTRAINT farm_farmer_id_not_null NOT NULL,
    name character varying(100) CONSTRAINT farm_name_not_null NOT NULL,
    address text CONSTRAINT farm_address_not_null NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public."Farm" OWNER TO postgres;

--
-- Name: Farmer; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Farmer" (
    id uuid NOT NULL,
    user_id uuid NOT NULL
);


ALTER TABLE public."Farmer" OWNER TO postgres;

--
-- Name: Item; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Item" (
    id uuid NOT NULL,
    name character varying(255) NOT NULL,
    description text NOT NULL,
    amount integer DEFAULT 0 NOT NULL,
    costpkilo double precision DEFAULT 0.0 NOT NULL,
    category public.category NOT NULL,
    img_path character varying(255) DEFAULT 'resources/images/no-image.jpg'::character varying NOT NULL
);


ALTER TABLE public."Item" OWNER TO postgres;

--
-- Name: Likes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Likes" (
    id uuid CONSTRAINT likes_id_not_null NOT NULL,
    user_id uuid CONSTRAINT likes_user_id_not_null NOT NULL,
    post_id uuid CONSTRAINT likes_post_id_not_null NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public."Likes" OWNER TO postgres;

--
-- Name: Order; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Order" (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    order_date timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    total_price numeric(10,2) NOT NULL,
    status public.order_status NOT NULL
);


ALTER TABLE public."Order" OWNER TO postgres;

--
-- Name: OrderItem; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."OrderItem" (
    id uuid NOT NULL,
    order_id uuid NOT NULL,
    item_id uuid NOT NULL,
    quantity integer NOT NULL,
    price_at_purchase numeric(10,2) NOT NULL
);


ALTER TABLE public."OrderItem" OWNER TO postgres;

--
-- Name: PaymentTransaction; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."PaymentTransaction" (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    order_id uuid,
    transaction_type public.transaction_type NOT NULL,
    amount numeric(10,2) NOT NULL,
    balance_before numeric(10,2) NOT NULL,
    balance_after numeric(10,2) NOT NULL,
    status public.transaction_status DEFAULT 'pending'::public.transaction_status NOT NULL,
    payment_method character varying(50),
    reference_number character varying(100),
    description text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    completed_at timestamp with time zone,
    CONSTRAINT amount_positive CHECK ((amount > (0)::numeric))
);


ALTER TABLE public."PaymentTransaction" OWNER TO postgres;

--
-- Name: Posts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Posts" (
    id uuid CONSTRAINT posts_id_not_null NOT NULL,
    farmer_id uuid CONSTRAINT posts_farmer_id_not_null NOT NULL,
    farm_id uuid,
    content text CONSTRAINT posts_content_not_null NOT NULL,
    image_url character varying(255) DEFAULT 'resources/images/no-image.jpg'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public."Posts" OWNER TO postgres;

--
-- Name: Review; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Review" (
    id uuid NOT NULL,
    userid uuid NOT NULL,
    itemid uuid NOT NULL,
    content text NOT NULL,
    rating double precision DEFAULT 0.0 NOT NULL
);


ALTER TABLE public."Review" OWNER TO postgres;

--
-- Name: TradeBid; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."TradeBid" (
    id uuid NOT NULL,
    trade_listing_id uuid NOT NULL,
    bidding_farmer_id uuid NOT NULL,
    bid_item_id uuid NOT NULL,
    bid_item_quantity numeric NOT NULL,
    status public.trade_bid_status NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public."TradeBid" OWNER TO postgres;

--
-- Name: TradeListing; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."TradeListing" (
    id uuid NOT NULL,
    offering_farmer_id uuid NOT NULL,
    offered_item_id uuid NOT NULL,
    offered_item_quantity numeric NOT NULL,
    desired_items text,
    status public.trade_listing_status NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamp with time zone,
    image_url character varying(255) DEFAULT 'resources/images/no-image.jpg'::character varying NOT NULL
);


ALTER TABLE public."TradeListing" OWNER TO postgres;

--
-- Name: User; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."User" (
    id uuid NOT NULL,
    username character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    profile_pic_url text DEFAULT '/resources/images/blank_profile.svg'::text NOT NULL
);


ALTER TABLE public."User" OWNER TO postgres;

--
-- Name: UserBalance; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."UserBalance" (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    balance numeric(10,2) DEFAULT 0.00 NOT NULL,
    currency character varying(3) DEFAULT 'PHP'::character varying,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT balance_non_negative CHECK ((balance >= (0)::numeric))
);


ALTER TABLE public."UserBalance" OWNER TO postgres;

--
-- Name: UserDetail; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."UserDetail" (
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    contact_number character varying(255) NOT NULL,
    created_date date DEFAULT CURRENT_DATE NOT NULL,
    address character varying(255) NOT NULL,
    id uuid NOT NULL
);


ALTER TABLE public."UserDetail" OWNER TO postgres;

--
-- Name: UserUserDetail; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."UserUserDetail" (
    user_id uuid NOT NULL,
    detail_id uuid NOT NULL
);


ALTER TABLE public."UserUserDetail" OWNER TO postgres;

--
-- Name: UsersItem; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."UsersItem" (
    user_id uuid NOT NULL,
    item_id uuid NOT NULL
);


ALTER TABLE public."UsersItem" OWNER TO postgres;

--
-- Data for Name: Cart; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Cart" (id, user_id, grand_total, created_at) FROM stdin;
31111111-1111-1111-1111-111111111111	c6554794-849f-4338-87c5-6db2e2f76514	0.00	2025-10-30 14:51:36.304859+08
7846a8ed-d79b-47f2-ba4f-c3ec750e9d07	d30869ec-fb97-46d8-85a3-82608c01f803	0.00	2025-10-30 17:36:58.19938+08
\.


--
-- Data for Name: CartItem; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."CartItem" (id, cart_id, item_id, quantity) FROM stdin;
2d0cc5a3-9362-4ee7-a9af-ad146641dcac	7846a8ed-d79b-47f2-ba4f-c3ec750e9d07	a3e1b9f2-7d94-4d3a-9b4a-111111111111	3
\.


--
-- Data for Name: Comment; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Comment" (id, user_id, post_id, content, created_at) FROM stdin;
c0000001-0001-0001-0001-000000000001	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	Great farming technique! Looking forward to trying this on my own farm.	2025-10-30 15:00:00
c0000001-0001-0001-0001-000000000002	c6554794-849f-4338-87c5-6db2e2f76514	11111111-1111-1111-1111-111111111111	Thanks for sharing this! Very informative.	2025-10-30 16:30:00
c0000001-0001-0001-0001-000000000003	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	11111111-1111-1111-1111-111111111111	This is exactly what I needed. Appreciate the detailed explanation!	2025-10-30 18:15:00
c0000002-0002-0002-0002-000000000001	c6554794-849f-4338-87c5-6db2e2f76514	22222222-2222-2222-2222-222222222222	I've been using crop rotation for years. Works wonders!	2025-10-30 15:45:00
c0000002-0002-0002-0002-000000000002	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	22222222-2222-2222-2222-222222222222	Which crops do you recommend rotating together?	2025-10-30 17:20:00
c0000003-0003-0003-0003-000000000001	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	33333333-3333-3333-3333-333333333333	Wow! That's an amazing harvest. Congratulations!	2025-10-31 09:00:00
c0000003-0003-0003-0003-000000000002	c6554794-849f-4338-87c5-6db2e2f76514	33333333-3333-3333-3333-333333333333	Your hard work really paid off. Well done!	2025-10-31 10:30:00
c0000003-0003-0003-0003-000000000003	d30869ec-fb97-46d8-85a3-82608c01f803	33333333-3333-3333-3333-333333333333	Thank you all for the kind words!	2025-10-31 12:00:00
c0000004-0004-0004-0004-000000000001	d30869ec-fb97-46d8-85a3-82608c01f803	44444444-4444-4444-4444-444444444444	Organic farming is the future! Keep it up!	2025-10-31 14:00:00
c0000004-0004-0004-0004-000000000002	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	44444444-4444-4444-4444-444444444444	Where do you get your organic fertilizer from?	2025-10-31 15:30:00
c0000005-0005-0005-0005-000000000001	d30869ec-fb97-46d8-85a3-82608c01f803	55555555-5555-5555-5555-555555555555	That irrigation system looks very efficient!	2025-11-01 08:30:00
c0000005-0005-0005-0005-000000000002	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	55555555-5555-5555-5555-555555555555	How much did the setup cost? I'm interested in installing something similar.	2025-11-01 10:00:00
c0000005-0005-0005-0005-000000000003	c6554794-849f-4338-87c5-6db2e2f76514	55555555-5555-5555-5555-555555555555	I can help with the installation if you need assistance!	2025-11-01 11:45:00
c0000006-0006-0006-0006-000000000001	d30869ec-fb97-46d8-85a3-82608c01f803	66666666-6666-6666-6666-666666666666	Natural pest control is always better than chemicals!	2025-11-01 13:00:00
c0000006-0006-0006-0006-000000000002	c6554794-849f-4338-87c5-6db2e2f76514	66666666-6666-6666-6666-666666666666	What plants do you use to repel pests naturally?	2025-11-01 14:30:00
\.


--
-- Data for Name: Farm; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Farm" (id, farmer_id, name, address, created_at) FROM stdin;
11111111-aaaa-aaaa-aaaa-111111111111	d30869ec-fb97-46d8-85a3-82608c01f803	Sunny Fields	123 Sunshine Road, Springfield	2025-10-30 08:00:00
22222222-bbbb-bbbb-bbbb-222222222222	c6554794-849f-4338-87c5-6db2e2f76514	Green Valley	456 Greenway Blvd, Greenville	2025-10-31 09:00:00
33333333-cccc-cccc-cccc-333333333333	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	Harvest Haven	789 Harvest Lane, Farmtown	2025-11-01 10:00:00
\.


--
-- Data for Name: Farmer; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Farmer" (id, user_id) FROM stdin;
8c8c73e8-0a16-4d3a-826d-75d50d7a758f	d30869ec-fb97-46d8-85a3-82608c01f803
9c6a8e18-4a0d-4a3a-8c6d-75d50d7a758f	c6554794-849f-4338-87c5-6db2e2f76514
ac7b9f28-5b1e-4b4a-9d7e-85d50d7a758f	6a24dd2b-d441-4b39-ab85-8fa2bd61065e
\.


--
-- Data for Name: Item; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Item" (id, name, description, amount, costpkilo, category, img_path) FROM stdin;
a3e1b9f2-7d94-4d3a-9b4a-111111111111	Banana	Ripe Cavendish bananas, sweet and soft	120	0.8	fruits	resources/images/no-image.jpg
b7f2c6d4-1aeb-4f5b-9c2b-222222222222	Tomato	Vine-ripened red tomatoes, juicy	80	1.2	vegetables	resources/images/no-image.jpg
c9d3e8a1-55b2-4f66-a123-333333333333	Rice (White)	Long-grain white rice, cleaned and polished	500	0.6	grains	resources/images/no-image.jpg
d4a6b7f9-3c11-46e8-b890-444444444444	Chicken (Live)	Local broiler chicken, small farm raised	25	3.5	livestock	resources/images/no-image.jpg
e2b4c5d6-9f88-4caa-b321-555555555555	Fresh Milk	Raw fresh milk, chilled, 1L packs	200	0.9	dairy	resources/images/no-image.jpg
f6c8d9e0-0b77-4d55-a789-666666666666	Handicraft Basket	Woven basket made from nipa and rattan	15	5	others	resources/images/no-image.jpg
d6767055-1150-43e2-86e6-8417eb7b52b7	Test Item	This is a test item	50	10	dairy	resources/images/Test Item_pic.png
1884b874-a59b-43cb-9575-0876ee862911	Test Item	This is a test item	50	10	dairy	resources/images/Test Item_pic.png
4bb15d40-af3d-4626-9b42-f6cca4aa2262	Test Item	This is a test item	50	10	dairy	resources/images/Test Item_pic.png
80d10778-4bf8-4827-95dd-25bc063ca3fe	Test Item	This is a test item	50	10	dairy	resources/images/Test Item_pic.png
b59f6931-97c6-46cd-b6e7-3b1f39ece849	Test Item	This is a test item	50	10	dairy	resources/images/Test Item_pic.png
dbe42e07-f53b-40a7-b3fe-90e03106a2ae	Test Item	This is a test item	50	10	dairy	resources/images/Test Item_pic.png
369a3beb-8e4d-4bb3-8e44-971866913977	Test Item	This is a test item	50	10	dairy	resources/images/Test Item_pic.png
7e478553-598a-49bf-9488-2289b258f830	Test Item	This is a test item	50	10	dairy	resources/images/Test Item_pic.png
af8b4666-d419-4dc0-8f85-c753e86f4cea	Test Item	This is a test item	50	10	dairy	resources/images/Test Item_pic.png
e5ed9386-cccf-4645-bbfd-7de3567fb27c	Test Item	This is a test item	50	10	dairy	resources/images/Test Item_pic.png
\.


--
-- Data for Name: Likes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Likes" (id, user_id, post_id, created_at) FROM stdin;
\.


--
-- Data for Name: Order; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Order" (id, user_id, order_date, total_price, status) FROM stdin;
11111111-1111-1111-1111-111111111111	d30869ec-fb97-46d8-85a3-82608c01f803	2025-10-30 14:51:36.300878+08	4.00	Packaging
b33a55e0-7c70-4854-8755-15fd039437e6	d30869ec-fb97-46d8-85a3-82608c01f803	2025-10-30 14:57:19.45874+08	5.00	Shipping
efbb5ef3-04a0-4256-8ca6-c73788fa4b80	d30869ec-fb97-46d8-85a3-82608c01f803	2025-10-30 14:58:17.105012+08	5.00	Shipping
c0d045a0-9068-4634-8a4c-b47eb2713e09	d30869ec-fb97-46d8-85a3-82608c01f803	2025-10-30 15:02:24.856009+08	5.00	Shipping
87d88c9e-7536-4155-9849-82f7d9f340fc	d30869ec-fb97-46d8-85a3-82608c01f803	2025-10-30 15:10:16.362426+08	5.00	Shipping
f97c3426-d2fa-441f-a934-5260b5f1499b	d30869ec-fb97-46d8-85a3-82608c01f803	2025-10-30 16:59:31.597732+08	5.00	Shipping
1ba41b22-cc20-4d55-a001-d018c5623fa2	d30869ec-fb97-46d8-85a3-82608c01f803	2025-10-30 17:01:16.026041+08	5.00	Shipping
5dd48b67-d43c-4fac-911c-2f2b0e205d15	d30869ec-fb97-46d8-85a3-82608c01f803	2025-10-30 17:38:57.490107+08	5.00	Shipping
0b38e06b-1db0-429c-8926-9ef89c7616db	d30869ec-fb97-46d8-85a3-82608c01f803	2025-10-30 17:41:11.453904+08	5.00	Shipping
a854ad9d-5b27-4399-b71e-6a92c2cae961	d30869ec-fb97-46d8-85a3-82608c01f803	2025-10-30 18:21:14.795816+08	5.00	Shipping
9b386052-49a8-4003-8f6e-b4862ddf70b5	d30869ec-fb97-46d8-85a3-82608c01f803	2025-10-31 09:44:06.368392+08	5.00	Shipping
6d533161-b7eb-4eed-9365-3911e88674b1	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-05 20:54:16.369795+08	5.00	Shipping
983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-07 16:17:39.928068+08	5.00	Shipping
deb53ec7-c5f2-4805-995f-e5c62e8d825d	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-07 16:19:28.872747+08	5.00	Shipping
d7aed1a8-2f22-46a0-a01c-a63ca2ae1497	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-08 15:09:24.380392+08	1234.56	Packaging
\.


--
-- Data for Name: OrderItem; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."OrderItem" (id, order_id, item_id, quantity, price_at_purchase) FROM stdin;
21111111-1111-1111-1111-111111111111	11111111-1111-1111-1111-111111111111	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	0.80
21111111-1111-1111-1111-222222222222	11111111-1111-1111-1111-111111111111	b7f2c6d4-1aeb-4f5b-9c2b-222222222222	2	1.20
607635f4-f251-414a-8f28-870f1f44dca1	b33a55e0-7c70-4854-8755-15fd039437e6	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
546d94a8-6db5-41f1-b22c-32c6dca89c16	b33a55e0-7c70-4854-8755-15fd039437e6	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
cfb4c7f2-dc20-4546-9766-e47013b128d7	efbb5ef3-04a0-4256-8ca6-c73788fa4b80	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
79c54f7c-1c18-4cfd-8c81-f690977d0581	efbb5ef3-04a0-4256-8ca6-c73788fa4b80	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
7ef09cff-05b2-4ed9-9168-87a3caa06246	c0d045a0-9068-4634-8a4c-b47eb2713e09	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
5b40ba22-ef7c-4770-b534-182a97cb202b	c0d045a0-9068-4634-8a4c-b47eb2713e09	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
037fffe8-6567-4de5-8ed9-29e988f4a07f	87d88c9e-7536-4155-9849-82f7d9f340fc	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
14193c38-912d-4af9-95da-27472221a031	87d88c9e-7536-4155-9849-82f7d9f340fc	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
e61f7969-3ae4-4f30-b2d7-753d6ea02e40	f97c3426-d2fa-441f-a934-5260b5f1499b	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
e42bc091-c824-408b-a7ea-3491e87d97fa	f97c3426-d2fa-441f-a934-5260b5f1499b	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
3018ced0-6d1e-4f7d-bd6d-17857e06cd03	1ba41b22-cc20-4d55-a001-d018c5623fa2	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
4c58c3e3-976c-4366-b154-9dc6856dcada	1ba41b22-cc20-4d55-a001-d018c5623fa2	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
698e7d14-88f4-431a-a1d9-980b7bad4880	5dd48b67-d43c-4fac-911c-2f2b0e205d15	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
7c7bc88f-44f3-4386-8066-fa4971b55c83	5dd48b67-d43c-4fac-911c-2f2b0e205d15	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
4313a826-fe81-4751-8df7-ce3322dd9eda	0b38e06b-1db0-429c-8926-9ef89c7616db	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
cc652090-4d7e-4fec-8dda-ad76f36d99d1	0b38e06b-1db0-429c-8926-9ef89c7616db	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
eeeae196-4c38-4992-b5fe-afe432524748	a854ad9d-5b27-4399-b71e-6a92c2cae961	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
cce6c6e8-bf80-420f-a7bc-f9e4b01f9fba	a854ad9d-5b27-4399-b71e-6a92c2cae961	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
950428d1-440b-4c0c-947a-7b0e87993137	9b386052-49a8-4003-8f6e-b4862ddf70b5	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
8748c4fc-b48f-4cd0-8d90-404a646f11f1	9b386052-49a8-4003-8f6e-b4862ddf70b5	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
35e177b1-58a1-4e11-b6a4-43c72d06f339	6d533161-b7eb-4eed-9365-3911e88674b1	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
d6adbca9-46df-4528-a19e-e76620c9cce2	6d533161-b7eb-4eed-9365-3911e88674b1	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
4c2a90b2-05b6-4786-a9c6-1d2a847759c2	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
264d0527-a041-4bdb-b5f6-0f5a3a017171	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
931946c9-7a3a-4326-8ef8-f09b58dfe02f	deb53ec7-c5f2-4805-995f-e5c62e8d825d	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
4b168118-c1a7-462d-9c25-e603f14be000	deb53ec7-c5f2-4805-995f-e5c62e8d825d	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
\.


--
-- Data for Name: PaymentTransaction; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."PaymentTransaction" (id, user_id, order_id, transaction_type, amount, balance_before, balance_after, status, payment_method, reference_number, description, metadata, created_at, updated_at, completed_at) FROM stdin;
20000001-0000-0000-0000-000000000001	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	5000.00	0.00	5000.00	completed	gcash	GC2025110800001	Initial deposit via GCash	\N	2025-10-30 08:30:00+08	2025-10-30 08:30:00+08	2025-10-30 08:30:00+08
20000002-0000-0000-0000-000000000002	c6554794-849f-4338-87c5-6db2e2f76514	\N	deposit	4000.00	0.00	4000.00	completed	maya	MAYA2025110800001	Deposit via Maya	\N	2025-10-31 09:15:00+08	2025-10-31 09:15:00+08	2025-10-31 09:15:00+08
20000003-0000-0000-0000-000000000003	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	\N	deposit	10000.00	0.00	10000.00	completed	bank_transfer	BNK2025110800001	Bank transfer deposit	\N	2025-11-01 10:00:00+08	2025-11-01 10:00:00+08	2025-11-01 10:00:00+08
20000004-0000-0000-0000-000000000004	c6554794-849f-4338-87c5-6db2e2f76514	\N	order_payment	500.00	4000.00	3500.00	completed	balance	\N	Payment for order - Fresh vegetables	\N	2025-11-02 11:30:00+08	2025-11-02 11:30:00+08	2025-11-02 11:30:00+08
20000005-0000-0000-0000-000000000005	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	\N	order_payment	1500.00	10000.00	8500.00	completed	balance	\N	Payment for order - Rice and grains	\N	2025-11-03 14:00:00+08	2025-11-03 14:00:00+08	2025-11-03 14:00:00+08
20000006-0000-0000-0000-000000000006	c6554794-849f-4338-87c5-6db2e2f76514	\N	deposit	500.00	3500.00	4000.00	completed	gcash	GC2025110800002	Additional deposit	\N	2025-11-04 10:00:00+08	2025-11-04 10:00:00+08	2025-11-04 10:00:00+08
20000007-0000-0000-0000-000000000007	c6554794-849f-4338-87c5-6db2e2f76514	\N	withdrawal	500.00	4000.00	3500.00	completed	gcash	GC2025110800003	Withdrawal to GCash	\N	2025-11-05 15:30:00+08	2025-11-05 15:30:00+08	2025-11-05 15:30:00+08
20000008-0000-0000-0000-000000000008	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	\N	refund	250.00	8500.00	8750.00	completed	balance	\N	Refund for cancelled item	\N	2025-11-06 09:00:00+08	2025-11-06 09:00:00+08	2025-11-06 09:00:00+08
20000009-0000-0000-0000-000000000009	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	5000.00	6000.00	completed	maya	MAYA2025110800002	Top-up via Maya	\N	2025-11-07 08:00:00+08	2025-11-07 08:00:00+08	2025-11-07 08:00:00+08
20000010-0000-0000-0000-000000000010	d30869ec-fb97-46d8-85a3-82608c01f803	\N	order_payment	1000.00	6000.00	5000.00	completed	balance	\N	Payment for order - Fresh fruits	\N	2025-11-07 14:30:00+08	2025-11-07 14:30:00+08	2025-11-07 14:30:00+08
20000011-0000-0000-0000-000000000011	c6554794-849f-4338-87c5-6db2e2f76514	\N	deposit	500.00	3500.00	4000.00	pending	gcash	GC2025110800004	Pending deposit	\N	2025-11-08 10:00:00+08	2025-11-08 10:00:00+08	\N
20000012-0000-0000-0000-000000000012	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	\N	withdrawal	1000.00	8750.00	7750.00	failed	bank_transfer	BNK2025110800002	Failed withdrawal - insufficient bank details	\N	2025-11-08 11:00:00+08	2025-11-08 11:00:00+08	\N
20000013-0000-0000-0000-000000000013	c6554794-849f-4338-87c5-6db2e2f76514	\N	deposit	200.00	4000.00	4200.00	processing	gcash	GC2025110800005	Processing deposit	\N	2025-11-08 11:30:00+08	2025-11-08 11:30:00+08	\N
dc5bbf52-e9a9-4cb2-8c3f-9cb6bf83c7d1	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	5000.00	5500.00	completed	gcash	GC6a56feec	Test deposit	\N	2025-11-08 14:47:14.136405+08	2025-11-08 14:47:14.136405+08	2025-11-08 14:47:14.136405+08
4e355838-1daf-4630-a796-fb145e328512	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	500.00	5500.00	5000.00	completed	test	cleanup	Test cleanup	\N	2025-11-08 14:47:14.145736+08	2025-11-08 14:47:14.145736+08	2025-11-08 14:47:14.145736+08
16e2ffbb-cf36-4154-a358-c5731a74dcad	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	5000.00	6000.00	completed	test	TESTc2a0c61a	Setup funds	\N	2025-11-08 14:47:14.153728+08	2025-11-08 14:47:14.153728+08	2025-11-08 14:47:14.153728+08
eb563dfe-fdcd-4577-befb-5cc6c3210c48	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	6000.00	6500.00	completed	test	TESTbfb91dce	Setup funds	\N	2025-11-08 14:47:14.178756+08	2025-11-08 14:47:14.178756+08	2025-11-08 14:47:14.178756+08
2ef5c1b0-6dd6-493a-bab5-a94513c5e415	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	200.00	6500.00	6300.00	completed	bank_transfer	BT00614bb0	Test withdrawal	\N	2025-11-08 14:47:14.183112+08	2025-11-08 14:47:14.183112+08	2025-11-08 14:47:14.183112+08
89e36e72-5bd3-47cd-ad28-6dee6d4500ff	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	300.00	6300.00	6600.00	completed	test	TEST2d5edfeb	Test deposit	\N	2025-11-08 14:47:14.198441+08	2025-11-08 14:47:14.198441+08	2025-11-08 14:47:14.198441+08
30d72958-c876-4ecd-81d1-e81158b945f2	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	6600.00	6700.00	completed	test	TEST1	Test deposit 1	\N	2025-11-08 14:47:14.214886+08	2025-11-08 14:47:14.214886+08	2025-11-08 14:47:14.214886+08
5a2b9381-e65a-434f-a145-b8e5522ee268	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	6700.00	6900.00	completed	test	TEST2	Test deposit 2	\N	2025-11-08 14:47:14.218132+08	2025-11-08 14:47:14.218132+08	2025-11-08 14:47:14.218132+08
1720da22-45e2-4932-87e7-9b3e68d588d2	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	50.00	6900.00	6950.00	completed	test	d2aa8f7b	Test deposit	\N	2025-11-08 14:47:14.227617+08	2025-11-08 14:47:14.227617+08	2025-11-08 14:47:14.227617+08
dd9e99b7-19d6-4e80-a072-5a4edc455781	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	6950.00	7050.00	completed	test	8ce5fbea	Test deposit	\N	2025-11-08 14:47:14.230762+08	2025-11-08 14:47:14.230762+08	2025-11-08 14:47:14.230762+08
b8e6658e-5ef8-4695-94ea-65b571f2c760	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	150.00	7050.00	7200.00	completed	test	11a701bd	Test deposit	\N	2025-11-08 14:47:14.233515+08	2025-11-08 14:47:14.233515+08	2025-11-08 14:47:14.233515+08
0ac8ce0a-0a63-44f7-bafb-3df88a41da33	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	7200.00	7400.00	completed	test	2691c50e	Test deposit	\N	2025-11-08 14:47:14.235634+08	2025-11-08 14:47:14.235634+08	2025-11-08 14:47:14.235634+08
549caadb-6042-4b09-8a62-26fda377e377	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	250.00	7400.00	7650.00	completed	test	547519ca	Test deposit	\N	2025-11-08 14:47:14.238478+08	2025-11-08 14:47:14.238478+08	2025-11-08 14:47:14.238478+08
fd3a7712-3b33-4f27-b24d-3257adeb241a	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	7650.00	8650.00	completed	gcash	GC12345678	Initial deposit	\N	2025-11-08 14:47:14.249078+08	2025-11-08 14:47:14.249078+08	2025-11-08 14:47:14.249078+08
ae5ef38b-cf42-46a4-bb05-27f2a1739f35	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	8650.00	9650.00	completed	test	TESTbf9c9f2d	Setup funds	\N	2025-11-08 14:49:06.259148+08	2025-11-08 14:49:06.259148+08	2025-11-08 14:49:06.259148+08
e259d5c7-13d3-4ffb-a9e9-ef02bad99a98	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	9650.00	10650.00	completed	test	TEST04a4c7f7	Setup funds	\N	2025-11-08 15:04:35.088239+08	2025-11-08 15:04:35.088239+08	2025-11-08 15:04:35.088239+08
404a9856-7545-4cb3-b65c-ff364e5ca056	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	10650.00	11650.00	completed	test	TEST4ba284fe	Setup funds	\N	2025-11-08 15:09:38.672608+08	2025-11-08 15:09:38.672608+08	2025-11-08 15:09:38.672608+08
960eb06e-d6be-4894-ac84-d41a42326bb3	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	11650.00	12650.00	completed	test	TESTded9afc2	Setup funds	\N	2025-11-08 15:14:57.819374+08	2025-11-08 15:14:57.819374+08	2025-11-08 15:14:57.819374+08
94837dc6-b537-4583-be48-503e9c3ed4f9	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	12650.00	13650.00	completed	test	TEST4cf5f045	Setup funds	\N	2025-11-08 15:23:17.638051+08	2025-11-08 15:23:17.638051+08	2025-11-08 15:23:17.638051+08
8bf398d8-9750-4456-b38d-d91bd55cc1d6	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	order_payment	250.00	13650.00	13400.00	completed	balance	\N	Payment for test order	\N	2025-11-08 15:23:17.645566+08	2025-11-08 15:23:17.645566+08	2025-11-08 15:23:17.645566+08
ad7e7f1d-dcad-49fe-8600-ff749744d833	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	refund	250.00	13400.00	13650.00	completed	balance	\N	Test cleanup	\N	2025-11-08 15:23:17.648842+08	2025-11-08 15:23:17.648842+08	2025-11-08 15:23:17.648842+08
787d9abf-7b0b-4c89-b0f1-615c38e4d9f8	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	13650.00	14650.00	completed	test	TEST6fd9a1cc	Setup funds	\N	2025-11-08 15:23:41.813863+08	2025-11-08 15:23:41.813863+08	2025-11-08 15:23:41.813863+08
d04387b7-b350-43e7-b749-45334a31c2d4	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	order_payment	250.00	14650.00	14400.00	completed	balance	\N	Payment for test order	\N	2025-11-08 15:23:41.817978+08	2025-11-08 15:23:41.817978+08	2025-11-08 15:23:41.817978+08
41586273-59bb-4ddf-ac1f-3cb2a95595bd	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	refund	250.00	14400.00	14650.00	completed	balance	\N	Test cleanup	\N	2025-11-08 15:23:41.820864+08	2025-11-08 15:23:41.820864+08	2025-11-08 15:23:41.820864+08
7a92cc9b-fe04-431c-8f3b-5140ed39edad	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	150.00	14650.00	14800.00	completed	balance	\N	Test refund for cancelled order	\N	2025-11-08 15:28:06.680446+08	2025-11-08 15:28:06.680446+08	2025-11-08 15:28:06.680446+08
bcb16a67-dfdd-46ea-a18d-45ae24cb730b	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	14800.00	15300.00	completed	gcash	GCb608a9e6	Test deposit	\N	2025-11-08 15:37:56.935611+08	2025-11-08 15:37:56.935611+08	2025-11-08 15:37:56.935611+08
fb1a5fe9-d700-4824-9859-805329ac329b	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	15300.00	16300.00	completed	test	TEST95bf1f70	Setup funds	\N	2025-11-08 15:37:56.950303+08	2025-11-08 15:37:56.950303+08	2025-11-08 15:37:56.950303+08
d0b9050b-84a3-4d36-baf6-16ff5d68ac68	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	order_payment	250.00	16300.00	16050.00	completed	balance	\N	Payment for test order	\N	2025-11-08 15:37:56.955338+08	2025-11-08 15:37:56.955338+08	2025-11-08 15:37:56.955338+08
a66c3855-a46a-4f6e-9b86-3e1979f51302	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	refund	250.00	16050.00	16300.00	completed	balance	\N	Test cleanup	\N	2025-11-08 15:37:56.959069+08	2025-11-08 15:37:56.959069+08	2025-11-08 15:37:56.959069+08
68def63f-a063-49f3-b5f4-7b4ceb08a8ed	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	150.00	16300.00	16450.00	completed	balance	\N	Test refund for cancelled order	\N	2025-11-08 15:37:56.980901+08	2025-11-08 15:37:56.980901+08	2025-11-08 15:37:56.980901+08
f17234a7-ea69-48d6-b332-53b9e456bd6c	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	16450.00	16950.00	completed	test	TEST4380baca	Setup funds	\N	2025-11-08 15:37:56.992719+08	2025-11-08 15:37:56.992719+08	2025-11-08 15:37:56.992719+08
58c3271e-0ae5-4eb1-9fa2-ae17f7b0c589	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	200.00	16950.00	16750.00	completed	bank_transfer	BT75bd2fa2	Test withdrawal	\N	2025-11-08 15:37:56.996799+08	2025-11-08 15:37:56.996799+08	2025-11-08 15:37:56.996799+08
5a3d9255-5ce1-47dd-91d9-68c82ca07f8e	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	300.00	16750.00	17050.00	completed	test	TEST59d9bbf2	Test deposit	\N	2025-11-08 15:37:57.016466+08	2025-11-08 15:37:57.016466+08	2025-11-08 15:37:57.016466+08
00695acd-1995-4555-ba52-80516d4ae00c	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	17050.00	17150.00	completed	test	43095157-e8ea-494c-b4a7-d52fbca3b901	Test deposit 1	\N	2025-11-08 15:37:57.035039+08	2025-11-08 15:37:57.035039+08	2025-11-08 15:37:57.035039+08
a26e8f26-9a1b-45d2-b1c6-f3a935f3c192	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	17150.00	17350.00	completed	test	8f9cac5e-c8b6-4fa8-9ebd-cae835dd9d81	Test deposit 2	\N	2025-11-08 15:37:57.038676+08	2025-11-08 15:37:57.038676+08	2025-11-08 15:37:57.038676+08
69aa9fd9-75c6-47db-8d1f-6e54ac72c99b	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	50.00	17350.00	17400.00	completed	test	fcca1bc8	Test deposit	\N	2025-11-08 15:37:57.04923+08	2025-11-08 15:37:57.04923+08	2025-11-08 15:37:57.04923+08
a4a19d9b-0714-4414-a6ee-5e9a93011378	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	17400.00	17500.00	completed	test	283e057a	Test deposit	\N	2025-11-08 15:37:57.053031+08	2025-11-08 15:37:57.053031+08	2025-11-08 15:37:57.053031+08
7385427b-2cd7-46ca-b879-e63e7c041f68	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	150.00	17500.00	17650.00	completed	test	99638547	Test deposit	\N	2025-11-08 15:37:57.055376+08	2025-11-08 15:37:57.055376+08	2025-11-08 15:37:57.055376+08
514b073a-30f6-4930-b713-ca5853c24f15	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	17650.00	17850.00	completed	test	da9542b4	Test deposit	\N	2025-11-08 15:37:57.057666+08	2025-11-08 15:37:57.057666+08	2025-11-08 15:37:57.057666+08
14656678-18fd-44de-bab0-e5f016d3fb53	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	250.00	17850.00	18100.00	completed	test	e7303f1a	Test deposit	\N	2025-11-08 15:37:57.059772+08	2025-11-08 15:37:57.059772+08	2025-11-08 15:37:57.059772+08
9fd77b24-d041-4474-8501-652f29bce644	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	18100.00	19100.00	completed	gcash	68278051-c551-4c8e-90dd-00ed6b5bf5c1	Initial deposit	\N	2025-11-08 15:38:32.054913+08	2025-11-08 15:38:32.054913+08	2025-11-08 15:38:32.054913+08
e1ed4391-7d6f-482b-b778-bca4b4640571	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	19100.00	20100.00	completed	gcash	924bdfcd-0978-4c36-af74-a8dc7cc18b34	Initial deposit	\N	2025-11-08 15:41:40.059315+08	2025-11-08 15:41:40.059315+08	2025-11-08 15:41:40.059315+08
3ff65431-45c4-457c-a9eb-691bea58c73e	d30869ec-fb97-46d8-85a3-82608c01f803	d7aed1a8-2f22-46a0-a01c-a63ca2ae1497	order_payment	400.00	20100.00	19700.00	completed	balance	\N	Order payment	\N	2025-11-08 15:41:40.063005+08	2025-11-08 15:41:40.063005+08	2025-11-08 15:41:40.063005+08
a143fe58-6e9f-4625-abb4-d0df6fd71b36	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	19700.00	20700.00	completed	gcash	5ba462d1-7673-4cb4-9bde-f8a495d07427	Initial deposit	\N	2025-11-08 15:42:30.236147+08	2025-11-08 15:42:30.236147+08	2025-11-08 15:42:30.236147+08
f71f3b69-9389-4844-9609-d46a79a2b7f1	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	order_payment	400.00	20700.00	20300.00	completed	balance	\N	Order payment	\N	2025-11-08 15:42:30.240975+08	2025-11-08 15:42:30.240975+08	2025-11-08 15:42:30.240975+08
3bde917f-630c-4c9b-a020-c11d75ae12fd	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	400.00	20300.00	20700.00	completed	balance	\N	Order cancelled	\N	2025-11-08 15:42:30.244103+08	2025-11-08 15:42:30.244103+08	2025-11-08 15:42:30.244103+08
18ab6628-fb75-4390-89a0-8f2e23201680	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	500.00	20700.00	20200.00	completed	bank	BT87654321	Withdrawal to bank	\N	2025-11-08 15:42:30.246216+08	2025-11-08 15:42:30.246216+08	2025-11-08 15:42:30.246216+08
75e8374f-fe7c-4a62-9ea2-794dc2d8bfa0	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	20200.00	20700.00	completed	gcash	GCd5c3eba3	Test deposit	\N	2025-11-08 15:42:34.76829+08	2025-11-08 15:42:34.76829+08	2025-11-08 15:42:34.76829+08
2ddf8158-9c5f-4a05-8b04-6d29a281bf5e	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	20700.00	21700.00	completed	test	TEST4bbd6405	Setup funds	\N	2025-11-08 15:42:34.777694+08	2025-11-08 15:42:34.777694+08	2025-11-08 15:42:34.777694+08
2d0e79ad-738f-41d0-86e9-280c7ec95307	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	order_payment	250.00	21700.00	21450.00	completed	balance	\N	Payment for test order	\N	2025-11-08 15:42:34.781497+08	2025-11-08 15:42:34.781497+08	2025-11-08 15:42:34.781497+08
eab22e9d-c3ca-4bcb-833e-bc9a87ccd43d	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	refund	250.00	21450.00	21700.00	completed	balance	\N	Test cleanup	\N	2025-11-08 15:42:34.784411+08	2025-11-08 15:42:34.784411+08	2025-11-08 15:42:34.784411+08
2e6f523b-e779-4767-a851-b6a7a4b6b73c	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	150.00	21700.00	21850.00	completed	balance	\N	Test refund for cancelled order	\N	2025-11-08 15:42:34.801939+08	2025-11-08 15:42:34.801939+08	2025-11-08 15:42:34.801939+08
186b8245-e211-49b2-bd80-0f51bc62aaf3	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	21850.00	22350.00	completed	test	TEST80dd604f	Setup funds	\N	2025-11-08 15:42:34.812527+08	2025-11-08 15:42:34.812527+08	2025-11-08 15:42:34.812527+08
ab7cdaeb-43fc-4a7c-9a0b-45997ef1c671	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	200.00	22350.00	22150.00	completed	bank_transfer	BT658cce33	Test withdrawal	\N	2025-11-08 15:42:34.816162+08	2025-11-08 15:42:34.816162+08	2025-11-08 15:42:34.816162+08
f9148766-62fc-4978-ab6b-4eef1a0a5919	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	300.00	22150.00	22450.00	completed	test	TEST09a8b154	Test deposit	\N	2025-11-08 15:42:34.833298+08	2025-11-08 15:42:34.833298+08	2025-11-08 15:42:34.833298+08
2fd99a11-29ed-466d-8019-957fbc6b15a0	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	22450.00	22550.00	completed	test	beaae3cb-db6b-4fe5-8658-b70cc41a83fd	Test deposit 1	\N	2025-11-08 15:42:34.849445+08	2025-11-08 15:42:34.849445+08	2025-11-08 15:42:34.849445+08
95f970c9-a24d-485c-b70b-9e5b39edbce0	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	22550.00	22750.00	completed	test	753a01d5-0710-428f-bd05-b714f82c0590	Test deposit 2	\N	2025-11-08 15:42:34.852976+08	2025-11-08 15:42:34.852976+08	2025-11-08 15:42:34.852976+08
122d39f0-e7fa-4af6-9f7b-15d58bc583c1	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	50.00	22750.00	22800.00	completed	test	6a1ad526	Test deposit	\N	2025-11-08 15:42:34.862136+08	2025-11-08 15:42:34.862136+08	2025-11-08 15:42:34.862136+08
f07ec690-f7df-450c-8af0-93744de57c97	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	22800.00	22900.00	completed	test	8ab2208e	Test deposit	\N	2025-11-08 15:42:34.866002+08	2025-11-08 15:42:34.866002+08	2025-11-08 15:42:34.866002+08
83780b2a-c306-440c-9734-0e92d3bf545c	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	150.00	22900.00	23050.00	completed	test	ed63d685	Test deposit	\N	2025-11-08 15:42:34.86795+08	2025-11-08 15:42:34.86795+08	2025-11-08 15:42:34.86795+08
7616f29a-5956-4d6a-9cb2-97c0fd14ccf4	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	23050.00	23250.00	completed	test	07463478	Test deposit	\N	2025-11-08 15:42:34.87029+08	2025-11-08 15:42:34.87029+08	2025-11-08 15:42:34.87029+08
1197620f-e06c-479b-9f38-f638be63a35a	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	250.00	23250.00	23500.00	completed	test	ffab13ca	Test deposit	\N	2025-11-08 15:42:34.87331+08	2025-11-08 15:42:34.87331+08	2025-11-08 15:42:34.87331+08
e36ad44e-82f6-487f-93e7-38cc8443e795	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	23500.00	24500.00	completed	gcash	bdeebae6-46b9-4aa8-8f96-506d10d29161	Initial deposit	\N	2025-11-08 15:42:34.881371+08	2025-11-08 15:42:34.881371+08	2025-11-08 15:42:34.881371+08
c0080da7-16f7-4dd3-8c7d-f4a7655ca08c	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	order_payment	400.00	24500.00	24100.00	completed	balance	\N	Order payment	\N	2025-11-08 15:42:34.884344+08	2025-11-08 15:42:34.884344+08	2025-11-08 15:42:34.884344+08
581a1bfe-f088-4798-a07e-e0f03543bf0a	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	400.00	24100.00	24500.00	completed	balance	\N	Order cancelled	\N	2025-11-08 15:42:34.886587+08	2025-11-08 15:42:34.886587+08	2025-11-08 15:42:34.886587+08
2f28b7d7-5758-4508-9948-d9ae03cc38ae	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	24500.00	25500.00	completed	gcash	a5c81de2-a013-4317-bcff-76d8b65c747c	Initial deposit	\N	2025-11-08 15:42:56.715679+08	2025-11-08 15:42:56.715679+08	2025-11-08 15:42:56.715679+08
1274611a-bfc4-40cb-a104-866accb04319	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	order_payment	400.00	25500.00	25100.00	completed	balance	\N	Order payment	\N	2025-11-08 15:42:56.718585+08	2025-11-08 15:42:56.718585+08	2025-11-08 15:42:56.718585+08
87cc2acf-f205-4755-bf5f-b207d66c1d9d	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	400.00	25100.00	25500.00	completed	balance	\N	Order cancelled	\N	2025-11-08 15:42:56.721076+08	2025-11-08 15:42:56.721076+08	2025-11-08 15:42:56.721076+08
36d0325c-1500-49d4-95de-882868406bea	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	500.00	25500.00	25000.00	completed	bank	c919b473-f120-41fb-bf72-244d41c2fd8a	Withdrawal to bank	\N	2025-11-08 15:42:56.723375+08	2025-11-08 15:42:56.723375+08	2025-11-08 15:42:56.723375+08
050d1ecd-b1f7-472e-a0a1-69966f23b9eb	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	25000.00	25500.00	completed	gcash	GC094c3fad	Test deposit	\N	2025-11-08 15:43:06.249835+08	2025-11-08 15:43:06.249835+08	2025-11-08 15:43:06.249835+08
9620aeff-b568-4390-94da-8aa7bcd6c00f	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	25500.00	26500.00	completed	test	TEST462e91c3	Setup funds	\N	2025-11-08 15:43:06.260309+08	2025-11-08 15:43:06.260309+08	2025-11-08 15:43:06.260309+08
a10fc7a3-21c2-4970-a9ae-1d16b1d37962	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	order_payment	250.00	26500.00	26250.00	completed	balance	\N	Payment for test order	\N	2025-11-08 15:43:06.264351+08	2025-11-08 15:43:06.264351+08	2025-11-08 15:43:06.264351+08
d1374ef9-f16a-4ce6-ab3a-2146e92ac1a9	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	refund	250.00	26250.00	26500.00	completed	balance	\N	Test cleanup	\N	2025-11-08 15:43:06.267892+08	2025-11-08 15:43:06.267892+08	2025-11-08 15:43:06.267892+08
81ae4d08-a1b8-4d8c-b7ed-0c958467ea7c	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	150.00	26500.00	26650.00	completed	balance	\N	Test refund for cancelled order	\N	2025-11-08 15:43:06.285255+08	2025-11-08 15:43:06.285255+08	2025-11-08 15:43:06.285255+08
8a861106-6f97-4280-bb4f-d82376fb1085	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	26650.00	27150.00	completed	test	TEST85526b8e	Setup funds	\N	2025-11-08 15:43:06.295265+08	2025-11-08 15:43:06.295265+08	2025-11-08 15:43:06.295265+08
1e07c9e7-815c-49d8-bae1-5e15f08c5793	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	200.00	27150.00	26950.00	completed	bank_transfer	BT733c8594	Test withdrawal	\N	2025-11-08 15:43:06.299735+08	2025-11-08 15:43:06.299735+08	2025-11-08 15:43:06.299735+08
baacbe9b-5106-44a3-af95-c2f1566e92a5	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	300.00	26950.00	27250.00	completed	test	TESTd947f4ac	Test deposit	\N	2025-11-08 15:43:06.318176+08	2025-11-08 15:43:06.318176+08	2025-11-08 15:43:06.318176+08
1502acb3-acc9-4768-9fc6-6759f57a4e9a	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	27250.00	27350.00	completed	test	56077b3f-148f-44ef-af1e-b4e4f9da2777	Test deposit 1	\N	2025-11-08 15:43:06.335107+08	2025-11-08 15:43:06.335107+08	2025-11-08 15:43:06.335107+08
03db4ffe-9fcd-40be-93a6-e8b8e15629c7	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	27350.00	27550.00	completed	test	b1b90106-0474-4675-a624-e3c9d96b78e6	Test deposit 2	\N	2025-11-08 15:43:06.338701+08	2025-11-08 15:43:06.338701+08	2025-11-08 15:43:06.338701+08
a60f638a-8498-47ab-9752-854d2754cd48	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	50.00	27550.00	27600.00	completed	test	244573da	Test deposit	\N	2025-11-08 15:43:06.348128+08	2025-11-08 15:43:06.348128+08	2025-11-08 15:43:06.348128+08
23a164d5-27d5-466b-871a-b2e1abf9b66c	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	27600.00	27700.00	completed	test	e7fe0978	Test deposit	\N	2025-11-08 15:43:06.352032+08	2025-11-08 15:43:06.352032+08	2025-11-08 15:43:06.352032+08
5f810e2f-61b2-4b90-b8ca-2d803c6686e4	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	150.00	27700.00	27850.00	completed	test	f971209e	Test deposit	\N	2025-11-08 15:43:06.354109+08	2025-11-08 15:43:06.354109+08	2025-11-08 15:43:06.354109+08
26e2caf4-83a7-442b-953f-0955fb9df443	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	27850.00	28050.00	completed	test	762deeff	Test deposit	\N	2025-11-08 15:43:06.355921+08	2025-11-08 15:43:06.355921+08	2025-11-08 15:43:06.355921+08
e429205c-a24c-434b-b7bd-f2a96d83308a	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	250.00	28050.00	28300.00	completed	test	29f47e11	Test deposit	\N	2025-11-08 15:43:06.357995+08	2025-11-08 15:43:06.357995+08	2025-11-08 15:43:06.357995+08
637e587e-5bcf-45e3-9d0e-b646be9fb711	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	28300.00	29300.00	completed	gcash	4c95b2e3-595d-4fcd-ad06-99d3e1fbdceb	Initial deposit	\N	2025-11-08 15:43:06.367224+08	2025-11-08 15:43:06.367224+08	2025-11-08 15:43:06.367224+08
e69712d3-c91e-4809-8cbc-735ac298895c	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	order_payment	400.00	29300.00	28900.00	completed	balance	\N	Order payment	\N	2025-11-08 15:43:06.370644+08	2025-11-08 15:43:06.370644+08	2025-11-08 15:43:06.370644+08
d169b478-800c-4834-ab25-76e7ab1285b7	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	400.00	28900.00	29300.00	completed	balance	\N	Order cancelled	\N	2025-11-08 15:43:06.37277+08	2025-11-08 15:43:06.37277+08	2025-11-08 15:43:06.37277+08
eb40f6c1-8a54-43ad-bfab-dfb8b99ecb45	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	500.00	29300.00	28800.00	completed	bank	a82288fd-47e0-4302-90de-e7fbe36388c0	Withdrawal to bank	\N	2025-11-08 15:43:06.375586+08	2025-11-08 15:43:06.375586+08	2025-11-08 15:43:06.375586+08
\.


--
-- Data for Name: Posts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Posts" (id, farmer_id, farm_id, content, image_url, created_at) FROM stdin;
22222222-2222-2222-2222-222222222222	c6554794-849f-4338-87c5-6db2e2f76514	22222222-bbbb-bbbb-bbbb-222222222222	Green Valley is now offering fresh organic produce!	resources/images/no-image.jpg	2025-10-31 10:00:00
33333333-3333-3333-3333-333333333333	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	33333333-cccc-cccc-cccc-333333333333	Harvest Haven is expanding! Stay tuned for updates.	resources/images/no-image.jpg	2025-11-01 09:30:00
44444444-4444-4444-4444-444444444444	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	Sunny Fields is hosting a community farming workshop!	resources/images/no-image.jpg	2025-11-02 10:00:00
55555555-5555-5555-5555-555555555555	c6554794-849f-4338-87c5-6db2e2f76514	22222222-bbbb-bbbb-bbbb-222222222222	Green Valley introduces a new line of organic herbs.	resources/images/no-image.jpg	2025-11-02 12:00:00
66666666-6666-6666-6666-666666666666	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	33333333-cccc-cccc-cccc-333333333333	Harvest Haven is now certified organic!	resources/images/no-image.jpg	2025-11-03 08:00:00
9af3ee81-346e-4cb5-9eba-2390fedb16dd	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	This is a test post from JohnDoe.	resources/images/9af3ee81-346e-4cb5-9eba-2390fedb16dd_post.png	2025-11-03 17:28:26.212629
f66ea4c3-f169-46e0-af01-4b9666d745bc	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	Final updated content	resources/images/f66ea4c3-f169-46e0-af01-4b9666d745bc_post.png	2025-11-03 17:49:54.264288
11111111-1111-1111-1111-111111111111	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	Updated content	resources/images/updated-image-url.jpg	2025-10-30 14:51:36
4ef9a372-7078-4171-94ff-249a553fad6d	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	This is a test post from JohnDoe.	resources/images/4ef9a372-7078-4171-94ff-249a553fad6d_post.png	2025-11-07 16:19:33.484823
4683f7a1-72d4-40b7-a725-19387799bcb9	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	This is a test post from JohnDoe.	resources/images/4683f7a1-72d4-40b7-a725-19387799bcb9_post.png	2025-11-04 17:18:56.113634
721b120c-a7ba-4f8e-bd90-db879713f475	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	This is a test post from JohnDoe.	resources/images/721b120c-a7ba-4f8e-bd90-db879713f475_post.png	2025-11-04 17:31:06.697544
d37e8c39-6329-426c-b1e5-8b26a7826a89	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	This is a test post from JohnDoe.	resources/images/d37e8c39-6329-426c-b1e5-8b26a7826a89_post.png	2025-11-04 17:58:12.951623
\.


--
-- Data for Name: Review; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Review" (id, userid, itemid, content, rating) FROM stdin;
894169e9-c907-4d89-84c4-3f1542488c9a	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	a3e1b9f2-7d94-4d3a-9b4a-111111111111	Banana very fluffy and has a premium after taste	4.9
8346bc9c-5efc-41bc-b7ea-68182849c473	d30869ec-fb97-46d8-85a3-82608c01f803	a3e1b9f2-7d94-4d3a-9b4a-111111111111	Very good Banana, cheap and very bananaey	4.6
f82ed553-05b3-41e3-817e-8e5ec4855ec0	c6554794-849f-4338-87c5-6db2e2f76514	a3e1b9f2-7d94-4d3a-9b4a-111111111111	Banana good for dishes and very yellow not very black	4.7
73222104-7f54-46a2-9b09-c4e4aba8d2bf	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	b7f2c6d4-1aeb-4f5b-9c2b-222222222222	Tomato very red. Tomatoes is delicious	4.5
\.


--
-- Data for Name: TradeBid; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."TradeBid" (id, trade_listing_id, bidding_farmer_id, bid_item_id, bid_item_quantity, status, created_at) FROM stdin;
6f7a8b9c-0d1e-2f3a-4b5c-6d7e8f9a0b1c	2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e	ac7b9f28-5b1e-4b4a-9d7e-85d50d7a758f	a3e1b9f2-7d94-4d3a-9b4a-111111111111	30	accepted	2025-10-30 12:04:04.204444+08
7a8b9c0d-1e2f-3a4b-5c6d-7e8f9a0b1c2d	2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	a3e1b9f2-7d94-4d3a-9b4a-111111111111	25	rejected	2025-10-29 12:04:04.204444+08
8ef8a66b-a13c-4e2a-aabc-f7f6981bfddc	1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	ac7b9f28-5b1e-4b4a-9d7e-85d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	18	pending	2025-11-07 14:06:15.579066+08
b1427605-04a7-4d47-a66c-acb94e18f32c	80832a67-4220-4300-a9f0-dcfa3b4c24ad	9c6a8e18-4a0d-4a3a-8c6d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	20	pending	2025-11-07 15:28:26.509664+08
b0e92c3c-ea9e-4533-837e-8fc380b2b05e	6792084e-0365-4f97-9c29-f8eabcb2240d	9c6a8e18-4a0d-4a3a-8c6d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	20	pending	2025-11-07 16:27:04.808898+08
924984e6-4f92-4b77-9d12-feb0323c509e	1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	ac7b9f28-5b1e-4b4a-9d7e-85d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	18	pending	2025-11-08 12:45:04.873461+08
\.


--
-- Data for Name: TradeListing; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."TradeListing" (id, offering_farmer_id, offered_item_id, offered_item_quantity, desired_items, status, created_at, expires_at, image_url) FROM stdin;
1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	b7f2c6d4-1aeb-4f5b-9c2b-222222222222	20	Looking for long-grain white rice.	open	2025-11-05 12:04:04.200954+08	2025-11-12 12:04:04.200954+08	resources/images/no-image.jpg
2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e	9c6a8e18-4a0d-4a3a-8c6d-75d50d7a758f	f6c8d9e0-0b77-4d55-a789-666666666666	5	Fresh bananas, around 30kg.	completed	2025-10-28 12:04:04.200954+08	2025-11-04 12:04:04.200954+08	resources/images/no-image.jpg
3c4d5e6f-7a8b-9c0d-1e2f-3a4b5c6d7e8f	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	a3e1b9f2-7d94-4d3a-9b4a-111111111111	50	Willing to trade for other fruits.	cancelled	2025-11-06 12:04:04.200954+08	\N	resources/images/no-image.jpg
83b09723-7268-4d49-afcb-e350ecc7b9ab	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	30	Testing status updates	open	2025-11-07 16:27:02.060228+08	\N	resources/images/no-image.jpg
e5b360f4-d777-4865-9d49-ab961a36fee7	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	cancelled	2025-11-07 13:15:21.474432+08	2025-11-14 13:15:21.469012+08	resources/images/no-image.jpg
d13b21f4-a5ee-47b8-a47c-8ef7c34e64f3	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	open	2025-11-07 13:46:15.771153+08	2025-11-14 13:46:15.763768+08	resources/images/no-image.jpg
f2f22e14-a250-4b5c-99db-497e31606c2a	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 13:46:15.783297+08	2025-11-14 13:46:15.777798+08	resources/images/no-image.jpg
6792084e-0365-4f97-9c29-f8eabcb2240d	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	b7f2c6d4-1aeb-4f5b-9c2b-222222222222	40	Looking for corn or wheat	completed	2025-11-07 16:27:04.805389+08	\N	resources/images/no-image.jpg
99fb53b7-dce2-4430-827e-17e7d9d76f0c	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	30	Testing status updates	open	2025-11-07 13:46:15.795257+08	2025-11-14 13:46:15.790372+08	resources/images/no-image.jpg
f8c5ce4b-ac4c-4285-a69f-2725037d81cf	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	b7f2c6d4-1aeb-4f5b-9c2b-222222222222	40	Looking for corn or wheat	open	2025-11-07 13:46:15.808464+08	2025-11-14 13:46:15.804039+08	resources/images/no-image.jpg
1db3d95b-4c5a-4ba5-ad65-897ad0780867	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	open	2025-11-07 13:46:34.408868+08	2025-11-14 13:46:34.403516+08	resources/images/no-image.jpg
df34d7b1-0a8e-4a36-a659-b79f096c1138	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 13:46:34.418584+08	2025-11-14 13:46:34.414556+08	resources/images/no-image.jpg
5951e4f9-518c-4742-930d-604fb876c9e2	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	30	Testing status updates	open	2025-11-07 13:46:34.428521+08	2025-11-14 13:46:34.423968+08	resources/images/no-image.jpg
04007d58-c72d-414e-87a3-be735b049e17	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	b7f2c6d4-1aeb-4f5b-9c2b-222222222222	40	Looking for corn or wheat	open	2025-11-07 13:46:34.441871+08	2025-11-14 13:46:34.437003+08	resources/images/no-image.jpg
529c5139-8d3f-4b4f-9b02-84b9e098a0be	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	open	2025-11-07 13:55:23.760849+08	2025-11-14 13:55:23.752467+08	resources/images/no-image.jpg
bfa8d185-4a0d-47bd-a461-34bc35b7b01f	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 13:55:23.775793+08	2025-11-14 13:55:23.766581+08	resources/images/no-image.jpg
8e963c15-b5bb-4093-8e69-682c1236a6d4	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	cancelled	2025-11-08 14:47:10.114734+08	2025-11-15 14:47:10.110217+08	resource/images/no-image.jpg
9f47f937-b9ff-4825-8533-d3b0266095a3	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	30	Testing status updates	open	2025-11-07 13:55:23.796094+08	2025-11-14 13:55:23.785413+08	resources/images/no-image.jpg
24d8d1e0-797f-4386-b576-99ab88d8596e	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	b7f2c6d4-1aeb-4f5b-9c2b-222222222222	40	Looking for corn or wheat	open	2025-11-07 13:55:23.822253+08	2025-11-14 13:55:23.812817+08	resources/images/no-image.jpg
59cbe248-e7fa-4563-9717-55be81ba3154	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 14:06:19.689212+08	2025-11-14 14:06:19.684129+08	resources/images/no-image.jpg
7e91fc20-c0ab-4be1-9e05-ac0a4489e8d3	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 14:07:50.067746+08	2025-11-14 14:07:50.062816+08	resources/images/no-image.jpg
ad114ab0-5fa0-45f0-a805-9bd5d6200f08	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 15:25:54.414071+08	2025-11-14 15:25:54.407804+08	resources/images/no-image.jpg
93997606-0f25-47dd-b722-25171985b81e	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 15:28:22.462548+08	2025-11-14 15:28:22.456118+08	resources/images/no-image.jpg
80832a67-4220-4300-a9f0-dcfa3b4c24ad	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	b7f2c6d4-1aeb-4f5b-9c2b-222222222222	40	Looking for corn or wheat	completed	2025-11-07 15:28:26.506122+08	2025-11-14 15:28:26.500773+08	resources/images/no-image.jpg
0aebaa77-41cd-4856-8087-03debdaa83e6	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	open	2025-11-07 16:26:54.680406+08	\N	resources/images/no-image.jpg
1893223c-7e6b-4099-af0b-66cdbfdc193a	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 16:26:58.045173+08	\N	resources/images/no-image.jpg
\.


--
-- Data for Name: User; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."User" (id, username, password, email, profile_pic_url) FROM stdin;
6a24dd2b-d441-4b39-ab85-8fa2bd61065e	StewardLittle	ratatouille	mouse@example.com	resources/images/blank_profile.svg
c6554794-849f-4338-87c5-6db2e2f76514	DanielGaliego	SirJ0elB@cay	danielgaliego@example.com	resources/images/blank_profile.svg
d30869ec-fb97-46d8-85a3-82608c01f803	JohnDoe	P@ssw0rd	JohnDoe@example.com	resources/images/blank_profile.svg
65ddb076-744c-40be-b90b-0b281149c9fa	JohnDoe	password123	JohnDoeee5d801c-ea99-402b-a739-cf834907db10@example.com	resources/images/JohnDoe_pfp.png
9ae195a0-05ff-446b-99c0-e6f09a0150d1	JohnDoe	password123	JohnDoe42b2c16b-8779-4d26-8585-0e203f6b0dcb@example.com	resources/images/JohnDoe_pfp.png
91254998-107c-4c10-b7f3-5ee4c21cebc2	JohnDoe	password123	JohnDoe785cf6da-2406-40f4-a4f5-7b72b4b0cbe0@example.com	resources/images/JohnDoe_pfp.png
543255dd-5325-4d3f-bcd2-ee6f8ac87e2e	JohnDoe	password123	JohnDoe2ed440bb-ec2b-47b7-b4f5-59b3940a6cf1@example.com	resources/images/JohnDoe_pfp.png
5b331f37-cdb2-45e5-8687-f7b326fc1c2b	JohnDoe	password123	JohnDoec145cef4-19bf-4491-9188-6be31f998a59@example.com	resources/images/JohnDoe_pfp.png
\.


--
-- Data for Name: UserBalance; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."UserBalance" (id, user_id, balance, currency, updated_at, created_at) FROM stdin;
10000002-0000-0000-0000-000000000002	c6554794-849f-4338-87c5-6db2e2f76514	3500.50	PHP	2025-11-08 12:10:00+08	2025-10-31 09:00:00+08
10000003-0000-0000-0000-000000000003	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	8250.75	PHP	2025-11-08 12:10:00+08	2025-11-01 10:00:00+08
10000001-0000-0000-0000-000000000001	d30869ec-fb97-46d8-85a3-82608c01f803	28800.00	PHP	2025-11-08 15:43:06.375586+08	2025-10-30 08:00:00+08
\.


--
-- Data for Name: UserDetail; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."UserDetail" (first_name, last_name, contact_number, created_date, address, id) FROM stdin;
John	Doe	+1-555-0123	2025-10-30	123 Main Street, Springfield	e1f2a3b4-5c6d-7e8f-9a0b-1c2d3e4f5a6b
John	Doe	1234567890	2025-10-31	123 Main St	b442af69-04c2-44cd-9f4e-46961ea73f4a
John	Doe	1234567890	2025-11-05	123 Main St	36faa21e-03c6-4ecc-9cdf-8e659267ea1d
John	Doe	1234567890	2025-11-05	123 Main St	accbba21-06fd-4113-b23d-5b944ae51f3d
John	Doe	1234567890	2025-11-05	123 Main St	99c2ec91-2152-4ff9-81cf-860526631d6e
John	Doe	1234567890	2025-11-08	123 Main St	40590c89-90b3-4daa-8f0f-aa77894f3bcf
\.


--
-- Data for Name: UserUserDetail; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."UserUserDetail" (user_id, detail_id) FROM stdin;
d30869ec-fb97-46d8-85a3-82608c01f803	e1f2a3b4-5c6d-7e8f-9a0b-1c2d3e4f5a6b
65ddb076-744c-40be-b90b-0b281149c9fa	b442af69-04c2-44cd-9f4e-46961ea73f4a
9ae195a0-05ff-446b-99c0-e6f09a0150d1	36faa21e-03c6-4ecc-9cdf-8e659267ea1d
91254998-107c-4c10-b7f3-5ee4c21cebc2	accbba21-06fd-4113-b23d-5b944ae51f3d
543255dd-5325-4d3f-bcd2-ee6f8ac87e2e	99c2ec91-2152-4ff9-81cf-860526631d6e
5b331f37-cdb2-45e5-8687-f7b326fc1c2b	40590c89-90b3-4daa-8f0f-aa77894f3bcf
\.


--
-- Data for Name: UsersItem; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."UsersItem" (user_id, item_id) FROM stdin;
d30869ec-fb97-46d8-85a3-82608c01f803	a3e1b9f2-7d94-4d3a-9b4a-111111111111
d30869ec-fb97-46d8-85a3-82608c01f803	b7f2c6d4-1aeb-4f5b-9c2b-222222222222
c6554794-849f-4338-87c5-6db2e2f76514	c9d3e8a1-55b2-4f66-a123-333333333333
\.


--
-- Name: CartItem CartItem_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."CartItem"
    ADD CONSTRAINT "CartItem_pkey" PRIMARY KEY (id);


--
-- Name: Cart Cart_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Cart"
    ADD CONSTRAINT "Cart_pkey" PRIMARY KEY (id);


--
-- Name: Farmer Farmer_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Farmer"
    ADD CONSTRAINT "Farmer_pkey" PRIMARY KEY (id);


--
-- Name: Farmer Farmer_user_id_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Farmer"
    ADD CONSTRAINT "Farmer_user_id_unique" UNIQUE (user_id);


--
-- Name: Item Item_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Item"
    ADD CONSTRAINT "Item_pkey" PRIMARY KEY (id);


--
-- Name: OrderItem OrderItem_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."OrderItem"
    ADD CONSTRAINT "OrderItem_pkey" PRIMARY KEY (id);


--
-- Name: Order Order_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Order"
    ADD CONSTRAINT "Order_pkey" PRIMARY KEY (id);


--
-- Name: PaymentTransaction PaymentTransaction_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."PaymentTransaction"
    ADD CONSTRAINT "PaymentTransaction_pkey" PRIMARY KEY (id);


--
-- Name: PaymentTransaction PaymentTransaction_reference_number_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."PaymentTransaction"
    ADD CONSTRAINT "PaymentTransaction_reference_number_key" UNIQUE (reference_number);


--
-- Name: Review Review_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Review"
    ADD CONSTRAINT "Review_pkey" PRIMARY KEY (id);


--
-- Name: TradeBid TradeBid_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."TradeBid"
    ADD CONSTRAINT "TradeBid_pkey" PRIMARY KEY (id);


--
-- Name: TradeListing TradeListing_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."TradeListing"
    ADD CONSTRAINT "TradeListing_pkey" PRIMARY KEY (id);


--
-- Name: UserBalance UserBalance_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."UserBalance"
    ADD CONSTRAINT "UserBalance_pkey" PRIMARY KEY (id);


--
-- Name: UserBalance UserBalance_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."UserBalance"
    ADD CONSTRAINT "UserBalance_user_id_key" UNIQUE (user_id);


--
-- Name: UserDetail UserDetail_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."UserDetail"
    ADD CONSTRAINT "UserDetail_pkey" PRIMARY KEY (id);


--
-- Name: User User_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."User"
    ADD CONSTRAINT "User_email_key" UNIQUE (email);


--
-- Name: User User_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."User"
    ADD CONSTRAINT "User_pkey" PRIMARY KEY (id);


--
-- Name: UsersItem UsersItem_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."UsersItem"
    ADD CONSTRAINT "UsersItem_pkey" PRIMARY KEY (user_id, item_id);


--
-- Name: Comment comment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Comment"
    ADD CONSTRAINT comment_pkey PRIMARY KEY (id);


--
-- Name: Farm farm_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Farm"
    ADD CONSTRAINT farm_pkey PRIMARY KEY (id);


--
-- Name: Likes likes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Likes"
    ADD CONSTRAINT likes_pkey PRIMARY KEY (id);


--
-- Name: Posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Posts"
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- Name: idx_payment_transaction_created_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_payment_transaction_created_at ON public."PaymentTransaction" USING btree (created_at DESC);


--
-- Name: idx_payment_transaction_order_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_payment_transaction_order_id ON public."PaymentTransaction" USING btree (order_id);


--
-- Name: idx_payment_transaction_reference; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_payment_transaction_reference ON public."PaymentTransaction" USING btree (reference_number);


--
-- Name: idx_payment_transaction_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_payment_transaction_status ON public."PaymentTransaction" USING btree (status);


--
-- Name: idx_payment_transaction_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_payment_transaction_user_id ON public."PaymentTransaction" USING btree (user_id);


--
-- Name: idx_user_balance_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_balance_user_id ON public."UserBalance" USING btree (user_id);


--
-- Name: CartItem CartItem_cart_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."CartItem"
    ADD CONSTRAINT "CartItem_cart_id_fkey" FOREIGN KEY (cart_id) REFERENCES public."Cart"(id) ON DELETE CASCADE;


--
-- Name: CartItem CartItem_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."CartItem"
    ADD CONSTRAINT "CartItem_item_id_fkey" FOREIGN KEY (item_id) REFERENCES public."Item"(id);


--
-- Name: Cart Cart_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Cart"
    ADD CONSTRAINT "Cart_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."User"(id);


--
-- Name: OrderItem OrderItem_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."OrderItem"
    ADD CONSTRAINT "OrderItem_item_id_fkey" FOREIGN KEY (item_id) REFERENCES public."Item"(id);


--
-- Name: OrderItem OrderItem_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."OrderItem"
    ADD CONSTRAINT "OrderItem_order_id_fkey" FOREIGN KEY (order_id) REFERENCES public."Order"(id) ON DELETE CASCADE;


--
-- Name: Order Order_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Order"
    ADD CONSTRAINT "Order_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."User"(id);


--
-- Name: PaymentTransaction PaymentTransaction_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."PaymentTransaction"
    ADD CONSTRAINT "PaymentTransaction_order_id_fkey" FOREIGN KEY (order_id) REFERENCES public."Order"(id) ON DELETE SET NULL;


--
-- Name: PaymentTransaction PaymentTransaction_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."PaymentTransaction"
    ADD CONSTRAINT "PaymentTransaction_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Review Review_itemid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Review"
    ADD CONSTRAINT "Review_itemid_fkey" FOREIGN KEY (itemid) REFERENCES public."Item"(id);


--
-- Name: Review Review_userid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Review"
    ADD CONSTRAINT "Review_userid_fkey" FOREIGN KEY (userid) REFERENCES public."User"(id);


--
-- Name: TradeBid TradeBid_bid_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."TradeBid"
    ADD CONSTRAINT "TradeBid_bid_item_id_fkey" FOREIGN KEY (bid_item_id) REFERENCES public."Item"(id) ON DELETE CASCADE;


--
-- Name: TradeBid TradeBid_bidding_farmer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."TradeBid"
    ADD CONSTRAINT "TradeBid_bidding_farmer_id_fkey" FOREIGN KEY (bidding_farmer_id) REFERENCES public."Farmer"(id) ON DELETE CASCADE;


--
-- Name: TradeBid TradeBid_trade_listing_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."TradeBid"
    ADD CONSTRAINT "TradeBid_trade_listing_id_fkey" FOREIGN KEY (trade_listing_id) REFERENCES public."TradeListing"(id) ON DELETE CASCADE;


--
-- Name: TradeListing TradeListing_offered_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."TradeListing"
    ADD CONSTRAINT "TradeListing_offered_item_id_fkey" FOREIGN KEY (offered_item_id) REFERENCES public."Item"(id) ON DELETE CASCADE;


--
-- Name: TradeListing TradeListing_offering_farmer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."TradeListing"
    ADD CONSTRAINT "TradeListing_offering_farmer_id_fkey" FOREIGN KEY (offering_farmer_id) REFERENCES public."Farmer"(id) ON DELETE CASCADE;


--
-- Name: UserBalance UserBalance_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."UserBalance"
    ADD CONSTRAINT "UserBalance_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: UserUserDetail UserUserDetail_detail_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."UserUserDetail"
    ADD CONSTRAINT "UserUserDetail_detail_id_fkey" FOREIGN KEY (detail_id) REFERENCES public."UserDetail"(id);


--
-- Name: UserUserDetail UserUserDetail_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."UserUserDetail"
    ADD CONSTRAINT "UserUserDetail_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."User"(id);


--
-- Name: UsersItem UsersItem_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."UsersItem"
    ADD CONSTRAINT "UsersItem_item_id_fkey" FOREIGN KEY (item_id) REFERENCES public."Item"(id);


--
-- Name: UsersItem UsersItem_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."UsersItem"
    ADD CONSTRAINT "UsersItem_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."User"(id);


--
-- Name: Comment comment_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Comment"
    ADD CONSTRAINT comment_post_id_fkey FOREIGN KEY (post_id) REFERENCES public."Posts"(id);


--
-- Name: Comment comment_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Comment"
    ADD CONSTRAINT comment_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."User"(id);


--
-- Name: Farm farm_farmer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Farm"
    ADD CONSTRAINT farm_farmer_id_fkey FOREIGN KEY (farmer_id) REFERENCES public."User"(id);


--
-- Name: Likes likes_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Likes"
    ADD CONSTRAINT likes_post_id_fkey FOREIGN KEY (post_id) REFERENCES public."Posts"(id);


--
-- Name: Likes likes_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Likes"
    ADD CONSTRAINT likes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."User"(id);


--
-- Name: Posts posts_farm_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Posts"
    ADD CONSTRAINT posts_farm_id_fkey FOREIGN KEY (farm_id) REFERENCES public."Farm"(id);


--
-- Name: Posts posts_farmer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Posts"
    ADD CONSTRAINT posts_farmer_id_fkey FOREIGN KEY (farmer_id) REFERENCES public."User"(id);


--
-- Name: Farmer user_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Farmer"
    ADD CONSTRAINT user_fk FOREIGN KEY (user_id) REFERENCES public."User"(id);


--
-- PostgreSQL database dump complete
--

\unrestrict HZJt83l80Aozg7JzSKKGhIVIixGpi0DHKDFyqox4tcR2hm8cKaxiPbGPAYROK4n

