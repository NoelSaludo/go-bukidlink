--
-- PostgreSQL database dump
--

\restrict FhezmbIePXo9E5xQqv1CdPbadJAazYuHNcUIyXDT20SAXUqg3l7MiqfJoHmEYM0

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
-- Data for Name: User; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."User" (id, username, password, email, profile_pic_url) FROM stdin;
6a24dd2b-d441-4b39-ab85-8fa2bd61065e	StewardLittle	ratatouille	mouse@example.com	resources/images/blank_profile.svg
c6554794-849f-4338-87c5-6db2e2f76514	DanielGaliego	SirJ0elB@cay	danielgaliego@example.com	resources/images/blank_profile.svg
d30869ec-fb97-46d8-85a3-82608c01f803	JohnDoe	P@ssw0rd	JohnDoe@example.com	resources/images/blank_profile.svg
65ddb076-744c-40be-b90b-0b281149c9fa	JohnDoe	password123	JohnDoeee5d801c-ea99-402b-a739-cf834907db10@example.com	resources/images/JohnDoe_pfp.png
\.


--
-- Data for Name: UserDetail; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."UserDetail" (first_name, last_name, contact_number, created_date, address, id) FROM stdin;
John	Doe	+1-555-0123	2025-10-30	123 Main Street, Springfield	e1f2a3b4-5c6d-7e8f-9a0b-1c2d3e4f5a6b
John	Doe	1234567890	2025-10-31	123 Main St	b442af69-04c2-44cd-9f4e-46961ea73f4a
\.


--
-- Data for Name: UserUserDetail; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."UserUserDetail" (user_id, detail_id) FROM stdin;
d30869ec-fb97-46d8-85a3-82608c01f803	e1f2a3b4-5c6d-7e8f-9a0b-1c2d3e4f5a6b
65ddb076-744c-40be-b90b-0b281149c9fa	b442af69-04c2-44cd-9f4e-46961ea73f4a
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
-- Name: Review Review_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Review"
    ADD CONSTRAINT "Review_pkey" PRIMARY KEY (id);


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
-- PostgreSQL database dump complete
--

\unrestrict FhezmbIePXo9E5xQqv1CdPbadJAazYuHNcUIyXDT20SAXUqg3l7MiqfJoHmEYM0

