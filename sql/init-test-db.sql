--
-- PostgreSQL database dump
--

\restrict KKbWjV9VDSMfGYeh2oRT9409BAiDIOEVo5LBiAkP0udkgKJIQV2IQPWKkrfMyQA

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

--
-- Name: check_duplicate_direct_conversation(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.check_duplicate_direct_conversation() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    participant_count INT;
    existing_conversation UUID;
BEGIN
    -- Only check for direct conversations
    IF (SELECT type FROM "Conversation" WHERE id = NEW.conversation_id) = 'direct' THEN
        -- Count participants in this conversation
        SELECT COUNT(*) INTO participant_count
        FROM "ConversationParticipant"
        WHERE conversation_id = NEW.conversation_id;
        
        -- Direct conversations should have exactly 2 participants
        IF participant_count >= 2 THEN
            RAISE EXCEPTION 'Direct conversations can only have 2 participants';
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.check_duplicate_direct_conversation() OWNER TO postgres;

--
-- Name: update_conversation_timestamp(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_conversation_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE "Conversation"
    SET updated_at = NEW.created_at
    WHERE id = NEW.conversation_id;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_conversation_timestamp() OWNER TO postgres;

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
-- Name: Conversation; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Conversation" (
    id uuid NOT NULL,
    title character varying(255),
    type character varying(20) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "Conversation_type_check" CHECK (((type)::text = ANY ((ARRAY['direct'::character varying, 'group'::character varying])::text[])))
);


ALTER TABLE public."Conversation" OWNER TO postgres;

--
-- Name: TABLE "Conversation"; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public."Conversation" IS 'Stores conversation metadata for both direct and group chats';


--
-- Name: COLUMN "Conversation".type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public."Conversation".type IS 'Type of conversation: direct (1-on-1) or group (multiple users)';


--
-- Name: ConversationParticipant; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."ConversationParticipant" (
    conversation_id uuid NOT NULL,
    user_id uuid NOT NULL,
    joined_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    last_read_at timestamp without time zone,
    is_muted boolean DEFAULT false
);


ALTER TABLE public."ConversationParticipant" OWNER TO postgres;

--
-- Name: TABLE "ConversationParticipant"; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public."ConversationParticipant" IS 'Junction table linking users to conversations with read tracking';


--
-- Name: COLUMN "ConversationParticipant".last_read_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public."ConversationParticipant".last_read_at IS 'Timestamp of last read message for unread count calculation';


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
-- Name: Message; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Message" (
    id uuid NOT NULL,
    conversation_id uuid NOT NULL,
    sender_id uuid NOT NULL,
    content text NOT NULL,
    message_type character varying(20) DEFAULT 'text'::character varying,
    attachment_url character varying(500),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    edited_at timestamp without time zone,
    is_deleted boolean DEFAULT false,
    CONSTRAINT "Message_message_type_check" CHECK (((message_type)::text = ANY ((ARRAY['text'::character varying, 'image'::character varying, 'file'::character varying, 'system'::character varying])::text[])))
);


ALTER TABLE public."Message" OWNER TO postgres;

--
-- Name: TABLE "Message"; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public."Message" IS 'Individual messages within conversations with support for text, images, and files';


--
-- Name: COLUMN "Message".message_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public."Message".message_type IS 'Type of message content: text, image, file, or system notification';


--
-- Name: COLUMN "Message".is_deleted; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public."Message".is_deleted IS 'Soft delete flag - messages are never hard deleted for history';


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
-- Name: UserPresence; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."UserPresence" (
    user_id uuid NOT NULL,
    status character varying(20) DEFAULT 'offline'::character varying NOT NULL,
    last_seen timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "UserPresence_status_check" CHECK (((status)::text = ANY ((ARRAY['online'::character varying, 'offline'::character varying, 'away'::character varying])::text[])))
);


ALTER TABLE public."UserPresence" OWNER TO postgres;

--
-- Name: TABLE "UserPresence"; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public."UserPresence" IS 'Tracks user online/offline status for real-time presence indicators';


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
-- Data for Name: Conversation; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Conversation" (id, title, type, created_at, updated_at) FROM stdin;
c3d4e5f6-a7b8-9c0d-1e2f-3c4d5e6f7a8b	Farm Cooperative Group	group	2024-01-10 08:00:00	2024-01-17 16:20:00
d4e5f6a7-b8c9-0d1e-2f3a-4d5e6f7a8b9c	Organic Farmers Alliance	group	2024-01-12 13:00:00	2024-01-18 10:15:00
8539d77f-99d2-4a9f-b73d-98002b1241b1	Test Group	group	2025-11-10 16:03:13.017842	2025-11-10 16:03:13.017842
6cdafb97-cb8a-489a-beef-9740c9ff3e9e	Project Team	group	2025-11-10 17:18:45.104973	2025-11-10 17:18:45.128228
\.


--
-- Data for Name: ConversationParticipant; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."ConversationParticipant" (conversation_id, user_id, joined_at, last_read_at, is_muted) FROM stdin;
c3d4e5f6-a7b8-9c0d-1e2f-3c4d5e6f7a8b	d30869ec-fb97-46d8-85a3-82608c01f803	2024-01-10 08:00:00	2024-01-17 16:20:00	f
c3d4e5f6-a7b8-9c0d-1e2f-3c4d5e6f7a8b	c6554794-849f-4338-87c5-6db2e2f76514	2024-01-10 08:05:00	2024-01-17 16:15:00	f
c3d4e5f6-a7b8-9c0d-1e2f-3c4d5e6f7a8b	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	2024-01-10 08:10:00	2024-01-17 16:00:00	t
c3d4e5f6-a7b8-9c0d-1e2f-3c4d5e6f7a8b	543255dd-5325-4d3f-bcd2-ee6f8ac87e2e	2024-01-10 08:15:00	2024-01-17 14:00:00	f
d4e5f6a7-b8c9-0d1e-2f3a-4d5e6f7a8b9c	d30869ec-fb97-46d8-85a3-82608c01f803	2024-01-12 13:00:00	2024-01-18 10:15:00	f
d4e5f6a7-b8c9-0d1e-2f3a-4d5e6f7a8b9c	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	2024-01-12 13:00:00	2024-01-18 10:10:00	f
d4e5f6a7-b8c9-0d1e-2f3a-4d5e6f7a8b9c	9ae195a0-05ff-446b-99c0-e6f09a0150d1	2024-01-12 13:00:00	2024-01-18 09:00:00	f
8539d77f-99d2-4a9f-b73d-98002b1241b1	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-10 16:03:13.017842	\N	f
8539d77f-99d2-4a9f-b73d-98002b1241b1	c6554794-849f-4338-87c5-6db2e2f76514	2025-11-10 16:03:13.017842	\N	f
6cdafb97-cb8a-489a-beef-9740c9ff3e9e	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-10 17:18:45.104973	\N	f
6cdafb97-cb8a-489a-beef-9740c9ff3e9e	c6554794-849f-4338-87c5-6db2e2f76514	2025-11-10 17:18:45.104973	\N	f
6cdafb97-cb8a-489a-beef-9740c9ff3e9e	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	2025-11-10 17:18:45.104973	\N	f
6cdafb97-cb8a-489a-beef-9740c9ff3e9e	543255dd-5325-4d3f-bcd2-ee6f8ac87e2e	2025-11-10 17:18:45.112098	\N	f
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
469a635f-3ba1-4f49-863c-0ff9a0806360	Test Item	This is a test item	50	10	dairy	resources/images/Test Item_pic.png
\.


--
-- Data for Name: Likes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Likes" (id, user_id, post_id, created_at) FROM stdin;
\.


--
-- Data for Name: Message; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Message" (id, conversation_id, sender_id, content, message_type, attachment_url, created_at, edited_at, is_deleted) FROM stdin;
102f9f0d-0f1e-4347-9336-d1e3ac020f55	6cdafb97-cb8a-489a-beef-9740c9ff3e9e	d30869ec-fb97-46d8-85a3-82608c01f803	Welcome to the team!	text	\N	2025-11-10 17:18:45.114835	\N	f
a225eaec-b5f1-47e2-a195-a2a043563f10	6cdafb97-cb8a-489a-beef-9740c9ff3e9e	c6554794-849f-4338-87c5-6db2e2f76514	Thanks for adding me!	text	\N	2025-11-10 17:18:45.128228	\N	f
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
d9c70e2f-6591-45cf-9319-a5d196e8a76b	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-10 13:37:53.150379+08	5.00	Shipping
3f4bc7a6-9bdb-4448-9ad1-2edd83b72a0f	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-10 13:58:34.178301+08	5.00	Shipping
72ef9d28-9b31-472b-82ef-30fab819e1c5	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-10 16:03:13.878554+08	5.00	Shipping
8cb57388-d953-4ad1-afa0-132443ffe5e4	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-10 16:03:33.960114+08	5.00	Shipping
0b494f91-6220-4844-a199-5c7eac95160f	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-10 17:26:40.562416+08	5.00	Shipping
daec4211-a930-4432-87c5-ff11776fce16	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-10 17:29:29.76203+08	5.00	Shipping
fec01151-5bfd-4793-a22d-c5a989e7bf0b	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-10 17:29:36.383773+08	5.00	Shipping
250c7cf6-06a0-4ec5-9af3-d8edcbf15932	d30869ec-fb97-46d8-85a3-82608c01f803	2025-11-10 17:33:17.537957+08	5.00	Shipping
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
3f7430aa-3057-4c61-aab3-be7e6dbd0821	d9c70e2f-6591-45cf-9319-a5d196e8a76b	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
c6a5850c-5982-406b-8a50-13190e5042c2	d9c70e2f-6591-45cf-9319-a5d196e8a76b	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
a4198a63-de7a-43a2-b418-ecc029617f3c	3f4bc7a6-9bdb-4448-9ad1-2edd83b72a0f	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
41dc5551-61a1-46b2-b26f-f626fe9cfda1	3f4bc7a6-9bdb-4448-9ad1-2edd83b72a0f	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
2f597083-a59d-45f7-8f03-97e8a4f075d7	72ef9d28-9b31-472b-82ef-30fab819e1c5	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
eac74ca4-6044-4d2f-a599-e793f4c49127	72ef9d28-9b31-472b-82ef-30fab819e1c5	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
e8ff56f7-89bc-4384-ba54-87c3f2b52661	8cb57388-d953-4ad1-afa0-132443ffe5e4	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
f0287232-dcef-4ff5-89d6-81b3d8257e53	8cb57388-d953-4ad1-afa0-132443ffe5e4	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
d42f8c38-0823-4947-bec0-83a9b77156e8	0b494f91-6220-4844-a199-5c7eac95160f	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
c1ac45ee-f721-40cf-966e-0bb73a2fd163	0b494f91-6220-4844-a199-5c7eac95160f	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
91724a94-2580-4c42-8f82-180d8900121c	daec4211-a930-4432-87c5-ff11776fce16	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
7d65b0fb-afbb-47c9-8522-95b10cced4bd	daec4211-a930-4432-87c5-ff11776fce16	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
69473e35-4240-4ab8-a66d-ea8f183dd4e3	fec01151-5bfd-4793-a22d-c5a989e7bf0b	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
7c89de23-5fec-42be-8b5c-f71386f4c0c4	fec01151-5bfd-4793-a22d-c5a989e7bf0b	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
59d933c3-ab08-4ce9-9fb7-f3fcb56e2deb	250c7cf6-06a0-4ec5-9af3-d8edcbf15932	a3e1b9f2-7d94-4d3a-9b4a-111111111111	2	1.00
a18a366e-d723-42d6-a8bd-fec9f85bbe80	250c7cf6-06a0-4ec5-9af3-d8edcbf15932	c9d3e8a1-55b2-4f66-a123-333333333333	1	3.00
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
0ffd41cb-dd26-4ee3-bbb5-81d83d485c76	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	28800.00	29300.00	completed	gcash	GC27cd04e6	Test deposit	\N	2025-11-10 13:37:52.827846+08	2025-11-10 13:37:52.827846+08	2025-11-10 13:37:52.827846+08
93c18335-b1bb-4918-9918-b5ea2fdcad77	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	29300.00	30300.00	completed	test	TEST3d37d3de	Setup funds	\N	2025-11-10 13:37:52.839303+08	2025-11-10 13:37:52.839303+08	2025-11-10 13:37:52.839303+08
87d64852-efe7-4141-b753-d3d4a78f4df7	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	order_payment	250.00	30300.00	30050.00	completed	balance	\N	Payment for test order	\N	2025-11-10 13:37:52.851873+08	2025-11-10 13:37:52.851873+08	2025-11-10 13:37:52.851873+08
6e66e310-c70c-4c3e-a9f8-8cb3d049a12c	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	refund	250.00	30050.00	30300.00	completed	balance	\N	Test cleanup	\N	2025-11-10 13:37:52.862569+08	2025-11-10 13:37:52.862569+08	2025-11-10 13:37:52.862569+08
48034ddc-ea3f-45c0-8b7c-7092c733e95c	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	150.00	30300.00	30450.00	completed	balance	\N	Test refund for cancelled order	\N	2025-11-10 13:37:52.897693+08	2025-11-10 13:37:52.897693+08	2025-11-10 13:37:52.897693+08
291a2187-17a2-4ca3-9835-c56f3a87aa15	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	30450.00	30950.00	completed	test	TEST2c720718	Setup funds	\N	2025-11-10 13:37:52.911593+08	2025-11-10 13:37:52.911593+08	2025-11-10 13:37:52.911593+08
47bbd2c8-6939-4089-9e88-a7aafeae73b4	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	200.00	30950.00	30750.00	completed	bank_transfer	BT851f14a5	Test withdrawal	\N	2025-11-10 13:37:52.919592+08	2025-11-10 13:37:52.919592+08	2025-11-10 13:37:52.919592+08
d5835283-49f7-47c0-a926-5a839acd15ac	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	300.00	30750.00	31050.00	completed	test	TEST76c1618a	Test deposit	\N	2025-11-10 13:37:52.945843+08	2025-11-10 13:37:52.945843+08	2025-11-10 13:37:52.945843+08
a6d72fe5-959b-4417-8e88-03e42aab7a68	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	31050.00	31150.00	completed	test	7036928e-62b3-47ce-9cb8-c7f7e228e322	Test deposit 1	\N	2025-11-10 13:37:52.967153+08	2025-11-10 13:37:52.967153+08	2025-11-10 13:37:52.967153+08
259cb295-14b3-4bf7-8d7d-7124dc599821	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	31150.00	31350.00	completed	test	9119ad4e-3f85-42d1-a2d2-d77f066e643f	Test deposit 2	\N	2025-11-10 13:37:52.971798+08	2025-11-10 13:37:52.971798+08	2025-11-10 13:37:52.971798+08
25e8686d-3ae3-4ae4-ba8c-248d74cd0e6e	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	50.00	31350.00	31400.00	completed	test	9959b6ac	Test deposit	\N	2025-11-10 13:37:52.983438+08	2025-11-10 13:37:52.983438+08	2025-11-10 13:37:52.983438+08
b91ab133-40df-40bd-af9b-7be9bbaa3255	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	31400.00	31500.00	completed	test	d6e24b11	Test deposit	\N	2025-11-10 13:37:52.987901+08	2025-11-10 13:37:52.987901+08	2025-11-10 13:37:52.987901+08
28458085-d46d-4e69-9be6-2a91d3be39fc	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	150.00	31500.00	31650.00	completed	test	37315b41	Test deposit	\N	2025-11-10 13:37:52.989913+08	2025-11-10 13:37:52.989913+08	2025-11-10 13:37:52.989913+08
fc57ccbd-5e6d-44a8-aade-cfefafe2363f	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	31650.00	31850.00	completed	test	577e0c27	Test deposit	\N	2025-11-10 13:37:52.991652+08	2025-11-10 13:37:52.991652+08	2025-11-10 13:37:52.991652+08
f1a06158-c83a-4274-a0e0-3364ebdf3665	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	250.00	31850.00	32100.00	completed	test	41c563d2	Test deposit	\N	2025-11-10 13:37:52.993481+08	2025-11-10 13:37:52.993481+08	2025-11-10 13:37:52.993481+08
f239ce0a-a915-46b1-a483-4d6442ff4460	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	32100.00	33100.00	completed	gcash	6c4aa6f9-b25a-49d4-9c9b-f9174442a02f	Initial deposit	\N	2025-11-10 13:37:53.005065+08	2025-11-10 13:37:53.005065+08	2025-11-10 13:37:53.005065+08
8ec106dd-cfa2-4a33-8821-826ee4403831	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	order_payment	400.00	33100.00	32700.00	completed	balance	\N	Order payment	\N	2025-11-10 13:37:53.009707+08	2025-11-10 13:37:53.009707+08	2025-11-10 13:37:53.009707+08
c8774c7e-e023-4d15-a363-362080af6f35	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	400.00	32700.00	33100.00	completed	balance	\N	Order cancelled	\N	2025-11-10 13:37:53.012873+08	2025-11-10 13:37:53.012873+08	2025-11-10 13:37:53.012873+08
f2435fbf-05c0-464f-942c-5d5df2fc9bb9	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	500.00	33100.00	32600.00	completed	bank	f01d35d4-98e7-44dc-9fae-c4877ac0d344	Withdrawal to bank	\N	2025-11-10 13:37:53.015272+08	2025-11-10 13:37:53.015272+08	2025-11-10 13:37:53.015272+08
5eab9e46-523c-4275-9ee0-03bb404bfff8	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	32600.00	33100.00	completed	gcash	GC5e94d7a8	Test deposit	\N	2025-11-10 16:03:13.536968+08	2025-11-10 16:03:13.536968+08	2025-11-10 16:03:13.536968+08
fa08785d-1dca-4f30-9047-c9b9ceb57285	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	33100.00	34100.00	completed	test	TESTf4e179fe	Setup funds	\N	2025-11-10 16:03:13.552664+08	2025-11-10 16:03:13.552664+08	2025-11-10 16:03:13.552664+08
32a81193-eda9-4f3e-8e44-b74cb44d117d	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	order_payment	250.00	34100.00	33850.00	completed	balance	\N	Payment for test order	\N	2025-11-10 16:03:13.559836+08	2025-11-10 16:03:13.559836+08	2025-11-10 16:03:13.559836+08
7dba1b49-5d5b-4607-9441-e555951e57dc	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	refund	250.00	33850.00	34100.00	completed	balance	\N	Test cleanup	\N	2025-11-10 16:03:13.562975+08	2025-11-10 16:03:13.562975+08	2025-11-10 16:03:13.562975+08
eaaa9ee2-894f-4f05-94e0-cf38cf9c9be0	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	150.00	34100.00	34250.00	completed	balance	\N	Test refund for cancelled order	\N	2025-11-10 16:03:13.583546+08	2025-11-10 16:03:13.583546+08	2025-11-10 16:03:13.583546+08
1b71ffff-16e9-452c-ae68-3fbe31b5a97d	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	34250.00	34750.00	completed	test	TEST37443b0f	Setup funds	\N	2025-11-10 16:03:13.598256+08	2025-11-10 16:03:13.598256+08	2025-11-10 16:03:13.598256+08
d0861098-af13-40a3-a599-dfe66ffc208a	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	200.00	34750.00	34550.00	completed	bank_transfer	BTb63cff18	Test withdrawal	\N	2025-11-10 16:03:13.602341+08	2025-11-10 16:03:13.602341+08	2025-11-10 16:03:13.602341+08
727100a6-f912-4e52-80e3-a29dd7e7d0d3	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	300.00	34550.00	34850.00	completed	test	TEST4dfd9564	Test deposit	\N	2025-11-10 16:03:13.624956+08	2025-11-10 16:03:13.624956+08	2025-11-10 16:03:13.624956+08
c97e882d-41ff-400c-b076-3301b882eb70	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	34850.00	34950.00	completed	test	9fbf9257-a754-4f7b-9e8e-a150ce7f1bbb	Test deposit 1	\N	2025-11-10 16:03:13.648337+08	2025-11-10 16:03:13.648337+08	2025-11-10 16:03:13.648337+08
0f32afc4-5d67-4f0d-83fa-d7605c9c4d31	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	34950.00	35150.00	completed	test	520cd9f1-b0de-468e-b7d3-3ae50615baeb	Test deposit 2	\N	2025-11-10 16:03:13.651739+08	2025-11-10 16:03:13.651739+08	2025-11-10 16:03:13.651739+08
2493dc22-8ede-4a5d-8357-65b95bd46490	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	50.00	35150.00	35200.00	completed	test	7b235b0c	Test deposit	\N	2025-11-10 16:03:13.663784+08	2025-11-10 16:03:13.663784+08	2025-11-10 16:03:13.663784+08
c61de8b0-0e64-4409-9fb8-d13a5b6e0aef	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	35200.00	35300.00	completed	test	59fb1373	Test deposit	\N	2025-11-10 16:03:13.668226+08	2025-11-10 16:03:13.668226+08	2025-11-10 16:03:13.668226+08
92f54aa2-87e9-4717-bacf-91064c9b43ef	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	150.00	35300.00	35450.00	completed	test	eb5d2a20	Test deposit	\N	2025-11-10 16:03:13.670579+08	2025-11-10 16:03:13.670579+08	2025-11-10 16:03:13.670579+08
11ae3c57-e4b5-4981-94a6-17a6203f071f	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	35450.00	35650.00	completed	test	8ad8f05c	Test deposit	\N	2025-11-10 16:03:13.672569+08	2025-11-10 16:03:13.672569+08	2025-11-10 16:03:13.672569+08
7d054520-80df-4116-bc92-d4cf2a3ecf73	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	250.00	35650.00	35900.00	completed	test	99664d97	Test deposit	\N	2025-11-10 16:03:13.674445+08	2025-11-10 16:03:13.674445+08	2025-11-10 16:03:13.674445+08
70461aed-664c-433f-bb44-07cd07a46f6a	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	35900.00	36900.00	completed	gcash	5bc78d84-c022-4fff-a73a-43bfb1eacdcd	Initial deposit	\N	2025-11-10 16:03:13.686451+08	2025-11-10 16:03:13.686451+08	2025-11-10 16:03:13.686451+08
fd20e765-b278-4e66-83e1-19507f7525f1	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	order_payment	400.00	36900.00	36500.00	completed	balance	\N	Order payment	\N	2025-11-10 16:03:13.69015+08	2025-11-10 16:03:13.69015+08	2025-11-10 16:03:13.69015+08
ad20db89-9187-491b-9050-7b27f268dfdf	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	400.00	36500.00	36900.00	completed	balance	\N	Order cancelled	\N	2025-11-10 16:03:13.692328+08	2025-11-10 16:03:13.692328+08	2025-11-10 16:03:13.692328+08
a0f862d0-d834-426a-89a6-33b6d63f86e3	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	500.00	36900.00	36400.00	completed	bank	db693d8c-2f87-4e4d-95df-41b7c069051e	Withdrawal to bank	\N	2025-11-10 16:03:13.694473+08	2025-11-10 16:03:13.694473+08	2025-11-10 16:03:13.694473+08
4350da85-8890-45b8-9357-364631139828	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	36400.00	36900.00	completed	gcash	GC4a283f3f	Test deposit	\N	2025-11-10 17:26:40.21532+08	2025-11-10 17:26:40.21532+08	2025-11-10 17:26:40.21532+08
44bca31c-44dd-4a16-85be-53c31e5d0a2b	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	36900.00	37900.00	completed	test	TEST2b1dade2	Setup funds	\N	2025-11-10 17:26:40.231223+08	2025-11-10 17:26:40.231223+08	2025-11-10 17:26:40.231223+08
ce05df38-16d2-4e2e-b516-e577d44c1e5c	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	order_payment	250.00	37900.00	37650.00	completed	balance	\N	Payment for test order	\N	2025-11-10 17:26:40.23616+08	2025-11-10 17:26:40.23616+08	2025-11-10 17:26:40.23616+08
add9d3ab-bfaa-46c2-a423-c9c7eb84f256	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	refund	250.00	37650.00	37900.00	completed	balance	\N	Test cleanup	\N	2025-11-10 17:26:40.239434+08	2025-11-10 17:26:40.239434+08	2025-11-10 17:26:40.239434+08
6ae3918b-cc99-4c14-99a1-263f274768fc	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	150.00	37900.00	38050.00	completed	balance	\N	Test refund for cancelled order	\N	2025-11-10 17:26:40.26204+08	2025-11-10 17:26:40.26204+08	2025-11-10 17:26:40.26204+08
12ffaebe-b767-4d83-b26d-305e6d1e4697	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	38050.00	38550.00	completed	test	TEST1a2e82f5	Setup funds	\N	2025-11-10 17:26:40.271926+08	2025-11-10 17:26:40.271926+08	2025-11-10 17:26:40.271926+08
6171b338-a78d-452a-9904-f0793d473a4d	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	200.00	38550.00	38350.00	completed	bank_transfer	BT3852c09d	Test withdrawal	\N	2025-11-10 17:26:40.275349+08	2025-11-10 17:26:40.275349+08	2025-11-10 17:26:40.275349+08
1e7251f8-9997-479b-994e-c4b94b97c289	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	300.00	38350.00	38650.00	completed	test	TEST2b14a604	Test deposit	\N	2025-11-10 17:26:40.297849+08	2025-11-10 17:26:40.297849+08	2025-11-10 17:26:40.297849+08
c857cb27-dad1-4131-b141-8aeb23eaab44	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	38650.00	38750.00	completed	test	888b9d50-9f86-4df3-ae41-28e05cce03f4	Test deposit 1	\N	2025-11-10 17:26:40.317219+08	2025-11-10 17:26:40.317219+08	2025-11-10 17:26:40.317219+08
42dff30d-f759-408e-b4da-108734ebd093	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	38750.00	38950.00	completed	test	1dbefc4d-4e30-4e1b-9d65-4aee2d25d219	Test deposit 2	\N	2025-11-10 17:26:40.320284+08	2025-11-10 17:26:40.320284+08	2025-11-10 17:26:40.320284+08
c3911cdc-15de-44b9-909c-178a82aceda7	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	50.00	38950.00	39000.00	completed	test	16f1306f	Test deposit	\N	2025-11-10 17:26:40.33143+08	2025-11-10 17:26:40.33143+08	2025-11-10 17:26:40.33143+08
5b9c2f65-c0bf-41c1-9feb-22eb4b2a9d95	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	39000.00	39100.00	completed	test	2119275f	Test deposit	\N	2025-11-10 17:26:40.335337+08	2025-11-10 17:26:40.335337+08	2025-11-10 17:26:40.335337+08
59312b62-372c-4848-8ccb-823adbf2ff94	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	150.00	39100.00	39250.00	completed	test	891dc6b5	Test deposit	\N	2025-11-10 17:26:40.337406+08	2025-11-10 17:26:40.337406+08	2025-11-10 17:26:40.337406+08
fba07b28-875c-4ff7-b5f5-c039b4b40a0d	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	39250.00	39450.00	completed	test	60520542	Test deposit	\N	2025-11-10 17:26:40.339637+08	2025-11-10 17:26:40.339637+08	2025-11-10 17:26:40.339637+08
6a3a2ddd-47c2-4754-bb17-0e36f0f02829	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	250.00	39450.00	39700.00	completed	test	54fc46ea	Test deposit	\N	2025-11-10 17:26:40.342422+08	2025-11-10 17:26:40.342422+08	2025-11-10 17:26:40.342422+08
b472adbf-e21d-4fb7-8506-f6c37d2b4aca	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	39700.00	40700.00	completed	gcash	06057111-3b00-4596-a102-a451765db8bd	Initial deposit	\N	2025-11-10 17:26:40.353356+08	2025-11-10 17:26:40.353356+08	2025-11-10 17:26:40.353356+08
b73a6b39-f74f-4929-98cf-a79b65c5da12	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	order_payment	400.00	40700.00	40300.00	completed	balance	\N	Order payment	\N	2025-11-10 17:26:40.356856+08	2025-11-10 17:26:40.356856+08	2025-11-10 17:26:40.356856+08
8356c68e-2bf1-4be7-8e21-be7a05b0dbe9	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	400.00	40300.00	40700.00	completed	balance	\N	Order cancelled	\N	2025-11-10 17:26:40.360248+08	2025-11-10 17:26:40.360248+08	2025-11-10 17:26:40.360248+08
9cdd086a-9690-4e0d-9f0d-cc80359ae017	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	500.00	40700.00	40200.00	completed	bank	36edded9-4f55-4181-99c2-de0e61589068	Withdrawal to bank	\N	2025-11-10 17:26:40.36337+08	2025-11-10 17:26:40.36337+08	2025-11-10 17:26:40.36337+08
06433e96-798e-49b8-8ed1-047173ef7238	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	40200.00	40700.00	completed	gcash	GC1eb0710d	Test deposit	\N	2025-11-10 17:29:36.098828+08	2025-11-10 17:29:36.098828+08	2025-11-10 17:29:36.098828+08
d125b1e6-ea25-402a-94ba-7c79fbe12312	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	40700.00	41700.00	completed	test	TEST3f2c9209	Setup funds	\N	2025-11-10 17:29:36.113205+08	2025-11-10 17:29:36.113205+08	2025-11-10 17:29:36.113205+08
9f02a7dc-8e66-4eed-bfd3-4e6326745eed	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	order_payment	250.00	41700.00	41450.00	completed	balance	\N	Payment for test order	\N	2025-11-10 17:29:36.118702+08	2025-11-10 17:29:36.118702+08	2025-11-10 17:29:36.118702+08
e05aac9a-6c23-41ee-a8c4-927329e142f7	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	refund	250.00	41450.00	41700.00	completed	balance	\N	Test cleanup	\N	2025-11-10 17:29:36.122237+08	2025-11-10 17:29:36.122237+08	2025-11-10 17:29:36.122237+08
f4bb286f-e668-419b-9864-5bdc7508ea68	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	150.00	41700.00	41850.00	completed	balance	\N	Test refund for cancelled order	\N	2025-11-10 17:29:36.143011+08	2025-11-10 17:29:36.143011+08	2025-11-10 17:29:36.143011+08
00149a18-c4ff-4103-8f58-e6e9c86491bd	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	41850.00	42350.00	completed	test	TEST3d8828d8	Setup funds	\N	2025-11-10 17:29:36.153406+08	2025-11-10 17:29:36.153406+08	2025-11-10 17:29:36.153406+08
350d346a-ee8f-4fbb-93ac-fa4204a4d42a	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	200.00	42350.00	42150.00	completed	bank_transfer	BT746055e9	Test withdrawal	\N	2025-11-10 17:29:36.157926+08	2025-11-10 17:29:36.157926+08	2025-11-10 17:29:36.157926+08
38b58b38-a420-43a3-9057-99eb7d62ecd3	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	300.00	42150.00	42450.00	completed	test	TEST26b387e3	Test deposit	\N	2025-11-10 17:29:36.179476+08	2025-11-10 17:29:36.179476+08	2025-11-10 17:29:36.179476+08
9f25c01f-890b-44e2-bf5b-549142bb8f38	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	42450.00	42550.00	completed	test	05f4290d-4c4a-4aab-a1a7-c0e2006262a9	Test deposit 1	\N	2025-11-10 17:29:36.198057+08	2025-11-10 17:29:36.198057+08	2025-11-10 17:29:36.198057+08
3d3dedad-f679-4eaa-a5d5-ff81ec83e073	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	42550.00	42750.00	completed	test	9748c74e-91c0-4632-b0fb-7d20e274a3a5	Test deposit 2	\N	2025-11-10 17:29:36.201568+08	2025-11-10 17:29:36.201568+08	2025-11-10 17:29:36.201568+08
f69b9909-ecd4-4519-b70f-28df504bb546	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	50.00	42750.00	42800.00	completed	test	345e720b	Test deposit	\N	2025-11-10 17:29:36.211393+08	2025-11-10 17:29:36.211393+08	2025-11-10 17:29:36.211393+08
8491fdb6-35bb-4759-9bc5-058784adc103	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	42800.00	42900.00	completed	test	2318a5eb	Test deposit	\N	2025-11-10 17:29:36.214686+08	2025-11-10 17:29:36.214686+08	2025-11-10 17:29:36.214686+08
1f6bc4d4-decd-4608-a153-a1f102c9edcb	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	150.00	42900.00	43050.00	completed	test	fdf69808	Test deposit	\N	2025-11-10 17:29:36.216802+08	2025-11-10 17:29:36.216802+08	2025-11-10 17:29:36.216802+08
167dc79e-b29b-4444-9381-58e9896a48aa	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	43050.00	43250.00	completed	test	15a3a7b2	Test deposit	\N	2025-11-10 17:29:36.21877+08	2025-11-10 17:29:36.21877+08	2025-11-10 17:29:36.21877+08
7268e549-8bec-4967-985a-54fc16d1f69b	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	250.00	43250.00	43500.00	completed	test	64254a60	Test deposit	\N	2025-11-10 17:29:36.220952+08	2025-11-10 17:29:36.220952+08	2025-11-10 17:29:36.220952+08
109b8437-be87-41a6-8cc4-788682b370f5	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	43500.00	44500.00	completed	gcash	2a9474f4-0fab-40ff-bcac-31034bdb5c97	Initial deposit	\N	2025-11-10 17:29:36.232192+08	2025-11-10 17:29:36.232192+08	2025-11-10 17:29:36.232192+08
fc8987f8-5631-4393-93b5-04eeee7dd84e	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	order_payment	400.00	44500.00	44100.00	completed	balance	\N	Order payment	\N	2025-11-10 17:29:36.23699+08	2025-11-10 17:29:36.23699+08	2025-11-10 17:29:36.23699+08
d131899e-8f72-4252-aa46-25170e735ac2	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	400.00	44100.00	44500.00	completed	balance	\N	Order cancelled	\N	2025-11-10 17:29:36.239661+08	2025-11-10 17:29:36.239661+08	2025-11-10 17:29:36.239661+08
3a00e778-77a1-4232-9fbd-1d6c651fe354	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	500.00	44500.00	44000.00	completed	bank	c285b6b6-b79d-459a-959d-d094f2dc5987	Withdrawal to bank	\N	2025-11-10 17:29:36.24247+08	2025-11-10 17:29:36.24247+08	2025-11-10 17:29:36.24247+08
05de16dc-734c-45f9-bb20-1a70565ecdc5	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	44000.00	44500.00	completed	gcash	GC2cadd980	Test deposit	\N	2025-11-10 17:33:17.257827+08	2025-11-10 17:33:17.257827+08	2025-11-10 17:33:17.257827+08
3310c63f-94a9-4666-b22a-e29bd2c6dc4e	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	44500.00	45500.00	completed	test	TEST7fcce4d8	Setup funds	\N	2025-11-10 17:33:17.270558+08	2025-11-10 17:33:17.270558+08	2025-11-10 17:33:17.270558+08
12536f55-e40d-48d9-9466-e17d86c0ac5d	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	order_payment	250.00	45500.00	45250.00	completed	balance	\N	Payment for test order	\N	2025-11-10 17:33:17.274716+08	2025-11-10 17:33:17.274716+08	2025-11-10 17:33:17.274716+08
6f4ff96f-4660-4346-a9b4-da3e127684c9	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-1111-1111-1111-111111111111	refund	250.00	45250.00	45500.00	completed	balance	\N	Test cleanup	\N	2025-11-10 17:33:17.278573+08	2025-11-10 17:33:17.278573+08	2025-11-10 17:33:17.278573+08
5f6df22e-63f3-4a79-a65f-109d2bda59a0	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	150.00	45500.00	45650.00	completed	balance	\N	Test refund for cancelled order	\N	2025-11-10 17:33:17.299757+08	2025-11-10 17:33:17.299757+08	2025-11-10 17:33:17.299757+08
4c5174fd-1337-43b1-ad6a-93446415a89d	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	500.00	45650.00	46150.00	completed	test	TESTa8db6622	Setup funds	\N	2025-11-10 17:33:17.313205+08	2025-11-10 17:33:17.313205+08	2025-11-10 17:33:17.313205+08
ce724069-7d04-4978-8980-ad0e4753367a	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	200.00	46150.00	45950.00	completed	bank_transfer	BT3f46d4a3	Test withdrawal	\N	2025-11-10 17:33:17.316801+08	2025-11-10 17:33:17.316801+08	2025-11-10 17:33:17.316801+08
ecfa3d29-443e-4bf9-bc1f-e25354cdf0f9	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	300.00	45950.00	46250.00	completed	test	TESTb235eec0	Test deposit	\N	2025-11-10 17:33:17.338086+08	2025-11-10 17:33:17.338086+08	2025-11-10 17:33:17.338086+08
2b03cad9-38af-4fc6-aba4-0261c67cb8af	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	46250.00	46350.00	completed	test	29cf76e7-96ae-48ce-8263-20e1b5494681	Test deposit 1	\N	2025-11-10 17:33:17.356929+08	2025-11-10 17:33:17.356929+08	2025-11-10 17:33:17.356929+08
738acdbc-0d74-4049-9c69-004f88108889	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	46350.00	46550.00	completed	test	1402180c-753e-454f-a6ae-7bd489103752	Test deposit 2	\N	2025-11-10 17:33:17.361217+08	2025-11-10 17:33:17.361217+08	2025-11-10 17:33:17.361217+08
f48b57b8-c588-4ada-84f5-419a78e67aa5	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	50.00	46550.00	46600.00	completed	test	279cfe52	Test deposit	\N	2025-11-10 17:33:17.372008+08	2025-11-10 17:33:17.372008+08	2025-11-10 17:33:17.372008+08
e29496da-f650-4796-ae8e-750b659a09a2	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	100.00	46600.00	46700.00	completed	test	38afb84d	Test deposit	\N	2025-11-10 17:33:17.375332+08	2025-11-10 17:33:17.375332+08	2025-11-10 17:33:17.375332+08
7a6e4873-fd55-4ead-9b37-916cb940e61a	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	150.00	46700.00	46850.00	completed	test	ff56301a	Test deposit	\N	2025-11-10 17:33:17.377813+08	2025-11-10 17:33:17.377813+08	2025-11-10 17:33:17.377813+08
5c405efe-894e-412b-94fc-946ac7688b5f	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	200.00	46850.00	47050.00	completed	test	6b5dbac0	Test deposit	\N	2025-11-10 17:33:17.380875+08	2025-11-10 17:33:17.380875+08	2025-11-10 17:33:17.380875+08
799f4553-964e-4810-8474-e959c8222df0	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	250.00	47050.00	47300.00	completed	test	74e13ee9	Test deposit	\N	2025-11-10 17:33:17.383131+08	2025-11-10 17:33:17.383131+08	2025-11-10 17:33:17.383131+08
79e5021f-dd3a-456d-9142-9b5359b1f689	d30869ec-fb97-46d8-85a3-82608c01f803	\N	deposit	1000.00	47300.00	48300.00	completed	gcash	8bec7939-2506-4574-9825-7e03f2a25a50	Initial deposit	\N	2025-11-10 17:33:17.39578+08	2025-11-10 17:33:17.39578+08	2025-11-10 17:33:17.39578+08
6dd28278-2b75-44c0-9c35-ef8882165cee	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	order_payment	400.00	48300.00	47900.00	completed	balance	\N	Order payment	\N	2025-11-10 17:33:17.401096+08	2025-11-10 17:33:17.401096+08	2025-11-10 17:33:17.401096+08
f225cbce-94c0-4a26-9ee1-718e0e6b2d29	d30869ec-fb97-46d8-85a3-82608c01f803	983a8b8f-2fcf-45ce-b26f-5965ae5fef8b	refund	400.00	47900.00	48300.00	completed	balance	\N	Order cancelled	\N	2025-11-10 17:33:17.403562+08	2025-11-10 17:33:17.403562+08	2025-11-10 17:33:17.403562+08
ba2ee199-1b5c-45e6-ab68-07a92a657ac0	d30869ec-fb97-46d8-85a3-82608c01f803	\N	withdrawal	500.00	48300.00	47800.00	completed	bank	5320802d-210b-4e7a-af53-df5ad75717b8	Withdrawal to bank	\N	2025-11-10 17:33:17.405609+08	2025-11-10 17:33:17.405609+08	2025-11-10 17:33:17.405609+08
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
4ef9a372-7078-4171-94ff-249a553fad6d	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	This is a test post from JohnDoe.	resources/images/4ef9a372-7078-4171-94ff-249a553fad6d_post.png	2025-11-07 16:19:33.484823
4683f7a1-72d4-40b7-a725-19387799bcb9	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	This is a test post from JohnDoe.	resources/images/4683f7a1-72d4-40b7-a725-19387799bcb9_post.png	2025-11-04 17:18:56.113634
9034d32a-589e-46d7-b4fd-b2e98778acd2	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	This is a test post from JohnDoe.	resources/images/9034d32a-589e-46d7-b4fd-b2e98778acd2_post.png	2025-11-10 13:37:48.698166
721b120c-a7ba-4f8e-bd90-db879713f475	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	This is a test post from JohnDoe.	resources/images/721b120c-a7ba-4f8e-bd90-db879713f475_post.png	2025-11-04 17:31:06.697544
d37e8c39-6329-426c-b1e5-8b26a7826a89	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	This is a test post from JohnDoe.	resources/images/d37e8c39-6329-426c-b1e5-8b26a7826a89_post.png	2025-11-04 17:58:12.951623
11111111-1111-1111-1111-111111111111	d30869ec-fb97-46d8-85a3-82608c01f803	11111111-aaaa-aaaa-aaaa-111111111111	Updated content	resources/images/updated-image-url.jpg	2025-10-30 14:51:36
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
60e73a9c-2b63-4189-9605-fc71567698b2	1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	ac7b9f28-5b1e-4b4a-9d7e-85d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	18	pending	2025-11-10 13:37:53.236883+08
233ae3ad-00c2-4831-a755-fc79b94a3085	1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	ac7b9f28-5b1e-4b4a-9d7e-85d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	18	rejected	2025-11-10 16:03:14.031253+08
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
9c5dbf59-fb2c-4bda-a5cf-9dfb08ff6d15	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	cancelled	2025-11-10 16:03:14.014379+08	2025-11-17 16:03:14.000755+08	resource/images/no-image.jpg
5951e4f9-518c-4742-930d-604fb876c9e2	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	30	Testing status updates	open	2025-11-07 13:46:34.428521+08	2025-11-14 13:46:34.423968+08	resources/images/no-image.jpg
04007d58-c72d-414e-87a3-be735b049e17	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	b7f2c6d4-1aeb-4f5b-9c2b-222222222222	40	Looking for corn or wheat	open	2025-11-07 13:46:34.441871+08	2025-11-14 13:46:34.437003+08	resources/images/no-image.jpg
529c5139-8d3f-4b4f-9b02-84b9e098a0be	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	open	2025-11-07 13:55:23.760849+08	2025-11-14 13:55:23.752467+08	resources/images/no-image.jpg
bfa8d185-4a0d-47bd-a461-34bc35b7b01f	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 13:55:23.775793+08	2025-11-14 13:55:23.766581+08	resources/images/no-image.jpg
8e963c15-b5bb-4093-8e69-682c1236a6d4	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	cancelled	2025-11-08 14:47:10.114734+08	2025-11-15 14:47:10.110217+08	resource/images/no-image.jpg
73a0f675-a1e5-4818-a10c-36304e759e06	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	open	2025-11-10 13:37:48.796383+08	\N	resources/images/no-image.jpg
9f47f937-b9ff-4825-8533-d3b0266095a3	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	30	Testing status updates	open	2025-11-07 13:55:23.796094+08	2025-11-14 13:55:23.785413+08	resources/images/no-image.jpg
24d8d1e0-797f-4386-b576-99ab88d8596e	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	b7f2c6d4-1aeb-4f5b-9c2b-222222222222	40	Looking for corn or wheat	open	2025-11-07 13:55:23.822253+08	2025-11-14 13:55:23.812817+08	resources/images/no-image.jpg
59cbe248-e7fa-4563-9717-55be81ba3154	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 14:06:19.689212+08	2025-11-14 14:06:19.684129+08	resources/images/no-image.jpg
7e91fc20-c0ab-4be1-9e05-ac0a4489e8d3	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 14:07:50.067746+08	2025-11-14 14:07:50.062816+08	resources/images/no-image.jpg
ad114ab0-5fa0-45f0-a805-9bd5d6200f08	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 15:25:54.414071+08	2025-11-14 15:25:54.407804+08	resources/images/no-image.jpg
93997606-0f25-47dd-b722-25171985b81e	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 15:28:22.462548+08	2025-11-14 15:28:22.456118+08	resources/images/no-image.jpg
80832a67-4220-4300-a9f0-dcfa3b4c24ad	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	b7f2c6d4-1aeb-4f5b-9c2b-222222222222	40	Looking for corn or wheat	completed	2025-11-07 15:28:26.506122+08	2025-11-14 15:28:26.500773+08	resources/images/no-image.jpg
0aebaa77-41cd-4856-8087-03debdaa83e6	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	open	2025-11-07 16:26:54.680406+08	\N	resources/images/no-image.jpg
1893223c-7e6b-4099-af0b-66cdbfdc193a	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-07 16:26:58.045173+08	\N	resources/images/no-image.jpg
379c2bbb-4460-42b2-af55-d41acae75c21	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	25	Looking for fruits	completed	2025-11-10 13:37:48.814641+08	\N	resources/images/no-image.jpg
ddbb39c9-7943-487e-89b5-22f6b24ea30f	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	cancelled	2025-11-10 16:03:37.246116+08	2025-11-17 16:03:37.241787+08	resource/images/no-image.jpg
5e2f5891-caad-4b57-a33e-19365e276d2f	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	30	Testing status updates	open	2025-11-10 13:37:48.826085+08	\N	resources/images/no-image.jpg
664bf742-7e91-494c-8743-244ab9088e08	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	b7f2c6d4-1aeb-4f5b-9c2b-222222222222	40	Looking for corn or wheat	completed	2025-11-10 13:37:48.839825+08	\N	resources/images/no-image.jpg
58a826bf-a10c-4e0f-9ad7-986ff5d0c135	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	cancelled	2025-11-10 13:37:53.223442+08	2025-11-17 13:37:53.217737+08	resource/images/no-image.jpg
2de7e192-b185-4ce0-abea-8b22525af3c8	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	cancelled	2025-11-10 15:09:50.081897+08	2025-11-17 15:09:50.076211+08	resource/images/no-image.jpg
7d7410c7-30a8-4c8f-95b4-f8be85b5eca9	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	cancelled	2025-11-10 16:26:48.257262+08	2025-11-17 16:26:48.25265+08	resource/images/no-image.jpg
2531c3a9-22b0-4105-b4e9-0fd96976d939	8c8c73e8-0a16-4d3a-826d-75d50d7a758f	c9d3e8a1-55b2-4f66-a123-333333333333	50	Looking for fresh vegetables, preferably tomatoes or lettuce.	cancelled	2025-11-10 17:26:40.65121+08	2025-11-17 17:26:40.644797+08	resource/images/no-image.jpg
\.


--
-- Data for Name: User; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."User" (id, username, password, email, profile_pic_url) FROM stdin;
6a24dd2b-d441-4b39-ab85-8fa2bd61065e	StewardLittle	ratatouille	mouse@example.com	resources/images/blank_profile.svg
c6554794-849f-4338-87c5-6db2e2f76514	DanielGaliego	SirJ0elB@cay	danielgaliego@example.com	resources/images/blank_profile.svg
d30869ec-fb97-46d8-85a3-82608c01f803	JohnDoe	P@ssw0rd	JohnDoe@example.com	resources/images/blank_profile.svg
65ddb076-744c-40be-b90b-0b281149c9fa	JohnDoe	password123	JohnDoeee5d801c-ea99-402b-a739-cf834907db10@example.com	resources/images/JohnDoe_pfp.png
91254998-107c-4c10-b7f3-5ee4c21cebc2	JohnDoe	password123	JohnDoe785cf6da-2406-40f4-a4f5-7b72b4b0cbe0@example.com	resources/images/JohnDoe_pfp.png
5b331f37-cdb2-45e5-8687-f7b326fc1c2b	JohnDoe	password123	JohnDoec145cef4-19bf-4491-9188-6be31f998a59@example.com	resources/images/JohnDoe_pfp.png
0831e138-2c50-4d59-8703-9e5493f17325	JohnDoe	password123	JohnDoe52f6c886-658f-4005-9775-37fb2bd7f748@example.com	resources/images/JohnDoe_pfp.png
543255dd-5325-4d3f-bcd2-ee6f8ac87e2e	MatthewConcepcion	password123	JohnDoe2ed440bb-ec2b-47b7-b4f5-59b3940a6cf1@example.com	resources/images/JohnDoe_pfp.png
9ae195a0-05ff-446b-99c0-e6f09a0150d1	MaryGrace	password123	JohnDoe42b2c16b-8779-4d26-8585-0e203f6b0dcb@example.com	resources/images/JohnDoe_pfp.png
\.


--
-- Data for Name: UserBalance; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."UserBalance" (id, user_id, balance, currency, updated_at, created_at) FROM stdin;
10000002-0000-0000-0000-000000000002	c6554794-849f-4338-87c5-6db2e2f76514	3500.50	PHP	2025-11-08 12:10:00+08	2025-10-31 09:00:00+08
10000003-0000-0000-0000-000000000003	6a24dd2b-d441-4b39-ab85-8fa2bd61065e	8250.75	PHP	2025-11-08 12:10:00+08	2025-11-01 10:00:00+08
10000001-0000-0000-0000-000000000001	d30869ec-fb97-46d8-85a3-82608c01f803	47800.00	PHP	2025-11-10 17:33:17.405609+08	2025-10-30 08:00:00+08
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
John	Doe	1234567890	2025-11-10	123 Main St	107fee0a-1c61-4e71-a733-00b24f949d20
\.


--
-- Data for Name: UserPresence; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."UserPresence" (user_id, status, last_seen, updated_at) FROM stdin;
6a24dd2b-d441-4b39-ab85-8fa2bd61065e	offline	2025-11-10 18:22:46.633398	2025-11-10 18:22:46.633398
543255dd-5325-4d3f-bcd2-ee6f8ac87e2e	offline	2025-11-10 14:02:32.076416	2025-11-10 14:02:32.076416
9ae195a0-05ff-446b-99c0-e6f09a0150d1	offline	2025-11-09 16:02:32.076416	2025-11-09 16:02:32.076416
c6554794-849f-4338-87c5-6db2e2f76514	offline	2025-11-10 19:12:48.982914	2025-11-10 19:12:48.982914
d30869ec-fb97-46d8-85a3-82608c01f803	online	2025-11-10 19:12:49.009765	2025-11-10 19:12:49.009765
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
0831e138-2c50-4d59-8703-9e5493f17325	107fee0a-1c61-4e71-a733-00b24f949d20
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
-- Name: ConversationParticipant ConversationParticipant_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ConversationParticipant"
    ADD CONSTRAINT "ConversationParticipant_pkey" PRIMARY KEY (conversation_id, user_id);


--
-- Name: Conversation Conversation_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Conversation"
    ADD CONSTRAINT "Conversation_pkey" PRIMARY KEY (id);


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
-- Name: Message Message_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message"
    ADD CONSTRAINT "Message_pkey" PRIMARY KEY (id);


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
-- Name: UserPresence UserPresence_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."UserPresence"
    ADD CONSTRAINT "UserPresence_pkey" PRIMARY KEY (user_id);


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
-- Name: idx_conversation_participant_user; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_conversation_participant_user ON public."ConversationParticipant" USING btree (user_id);


--
-- Name: idx_conversation_updated_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_conversation_updated_at ON public."Conversation" USING btree (updated_at DESC);


--
-- Name: idx_message_conversation; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_message_conversation ON public."Message" USING btree (conversation_id, created_at DESC);


--
-- Name: idx_message_sender; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_message_sender ON public."Message" USING btree (sender_id);


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
-- Name: idx_user_presence_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_presence_status ON public."UserPresence" USING btree (status);


--
-- Name: ConversationParticipant trigger_check_duplicate_direct_conversation; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trigger_check_duplicate_direct_conversation BEFORE INSERT ON public."ConversationParticipant" FOR EACH ROW EXECUTE FUNCTION public.check_duplicate_direct_conversation();


--
-- Name: Message trigger_update_conversation_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trigger_update_conversation_timestamp AFTER INSERT ON public."Message" FOR EACH ROW EXECUTE FUNCTION public.update_conversation_timestamp();


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
-- Name: ConversationParticipant ConversationParticipant_conversation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ConversationParticipant"
    ADD CONSTRAINT "ConversationParticipant_conversation_id_fkey" FOREIGN KEY (conversation_id) REFERENCES public."Conversation"(id) ON DELETE CASCADE;


--
-- Name: ConversationParticipant ConversationParticipant_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."ConversationParticipant"
    ADD CONSTRAINT "ConversationParticipant_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."User"(id) ON DELETE CASCADE;


--
-- Name: Message Message_conversation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message"
    ADD CONSTRAINT "Message_conversation_id_fkey" FOREIGN KEY (conversation_id) REFERENCES public."Conversation"(id) ON DELETE CASCADE;


--
-- Name: Message Message_sender_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Message"
    ADD CONSTRAINT "Message_sender_id_fkey" FOREIGN KEY (sender_id) REFERENCES public."User"(id) ON DELETE CASCADE;


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
-- Name: UserPresence UserPresence_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."UserPresence"
    ADD CONSTRAINT "UserPresence_user_id_fkey" FOREIGN KEY (user_id) REFERENCES public."User"(id) ON DELETE CASCADE;


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

\unrestrict KKbWjV9VDSMfGYeh2oRT9409BAiDIOEVo5LBiAkP0udkgKJIQV2IQPWKkrfMyQA

