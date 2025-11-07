# Trading API Endpoints

## Overview
The trading feature allows farmers to create listings offering items for trade and receive bids from other farmers. This implements an auction-style barter system.

## Endpoints

### 1. Get Trade Listings (Batch)
Retrieve a paginated list of trade listings with embedded base64-encoded images.

**Endpoint:** `GET /trades/batch?block=0`

**Query Parameters:**
- `block` (optional, default: 0) - Page number for pagination (100 listings per page)

**Response:** `200 OK`

Returns an array of enriched listing objects. Each object contains:
- `listing` - The trade listing data
- `image_base64` (optional) - Base64-encoded image data if the listing has an image
- `image_content_type` (optional) - MIME type of the image (e.g., "image/png", "image/jpeg")

```json
[
  {
    "listing": {
      "id": "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
      "offering_farmer_id": "8c8c73e8-0a16-4d3a-826d-75d50d7a758f",
      "offered_item_id": "b7f2c6d4-1aeb-4f5b-9c2b-222222222222",
      "offered_item_quantity": 20,
      "desired_items": "Looking for long-grain white rice.",
      "status": "open",
      "image_url": "resources/images/1a2b3c4d_trade.png",
      "created_at": "2025-11-05T10:30:00Z",
      "expires_at": "2025-11-12T10:30:00Z"
    },
    "image_base64": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
    "image_content_type": "image/png"
  }
]
```

**Example:**
```bash
curl -i http://localhost:8080/trades/batch?block=0
```

---

### 2. Get Single Trade Listing
Retrieve a specific trade listing with all its bids and embedded image data.

**Endpoint:** `GET /trade?id=<uuid>`

**Query Parameters:**
- `id` (required) - UUID of the trade listing

**Response:** `200 OK`

Returns an object containing:
- `listing` - The trade listing data
- `bids` - Array of all bids for this listing
- `image_base64` (optional) - Base64-encoded image data if the listing has an image
- `image_content_type` (optional) - MIME type of the image

```json
{
  "listing": {
    "id": "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
    "offering_farmer_id": "8c8c73e8-0a16-4d3a-826d-75d50d7a758f",
    "offered_item_id": "b7f2c6d4-1aeb-4f5b-9c2b-222222222222",
    "offered_item_quantity": 20,
    "desired_items": "Looking for long-grain white rice.",
    "status": "open",
    "image_url": "resources/images/1a2b3c4d_trade.jpg",
    "created_at": "2025-11-05T10:30:00Z",
    "expires_at": "2025-11-12T10:30:00Z"
  },
  "bids": [
    {
      "id": "4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a",
      "trade_listing_id": "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
      "bidding_farmer_id": "9c6a8e18-4a0d-4a3a-8c6d-75d50d7a758f",
      "bid_item_id": "c9d3e8a1-55b2-4f66-a123-333333333333",
      "bid_item_quantity": 15,
      "status": "pending",
      "created_at": "2025-11-06T14:20:00Z"
    }
  ],
  "image_base64": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
  "image_content_type": "image/jpeg"
}
```

**Response:** `404 Not Found` (if listing doesn't exist)
```json
{
  "error": "trade listing not found"
}
```

**Example:**
```bash
curl -i http://localhost:8080/trade?id=1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d
```

---

### 3. Create Trade Listing
Create a new trade listing offering an item for trade, with optional image upload.

**Endpoint:** `POST /trade`

**Headers:**
- `Content-Type: application/json`

**Request Body:**

The request uses an envelope format with two main objects:
- `listing` (required) - The trade listing data
- `listing_image` (optional) - Base64-encoded image data

```json
{
  "listing": {
    "offering_farmer_id": "8c8c73e8-0a16-4d3a-826d-75d50d7a758f",
    "offered_item_id": "b7f2c6d4-1aeb-4f5b-9c2b-222222222222",
    "offered_item_quantity": 20,
    "desired_items": "Looking for long-grain white rice.",
    "expires_at": "2025-11-12T23:59:59Z"
  },
  "listing_image": {
    "base64": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
    "content_type": "image/png"
  }
}
```

**Request Fields:**

`listing` object (all fields required except `expires_at`):
- `offering_farmer_id` (string, UUID) - ID of the farmer offering the trade
- `offered_item_id` (string, UUID) - ID of the item being offered
- `offered_item_quantity` (number) - Quantity of the offered item
- `desired_items` (string) - Description of desired items in exchange
- `expires_at` (string, RFC3339, optional) - When the listing expires

`listing_image` object (optional):
- `base64` (string) - Base64-encoded image data
- `content_type` (string) - MIME type of the image (e.g., "image/png", "image/jpeg", "image/gif", "image/webp")

**Notes:**
- `id` is auto-generated and returned in the response
- `status` defaults to "open"
- `image_url` is auto-generated when an image is provided (saved as `<listing_id>_trade.<ext>`)
- If no image is provided, defaults to `"resources/images/no-image.jpg"`
- `created_at` is auto-generated by database
- Image is decoded from base64 and saved to `resources/images/` directory

**Response:** `201 Created`
```json
{
  "status": "trade listing created",
  "listing_id": "7b8c9d0e-1f2a-3b4c-5d6e-7f8a9b0c1d2e"
}
```

**Example without image:**
```bash
curl -i -X POST http://localhost:8080/trade \
  -H "Content-Type: application/json" \
  -d '{
    "listing": {
      "offering_farmer_id": "8c8c73e8-0a16-4d3a-826d-75d50d7a758f",
      "offered_item_id": "b7f2c6d4-1aeb-4f5b-9c2b-222222222222",
      "offered_item_quantity": 20,
      "desired_items": "Looking for long-grain white rice.",
      "expires_at": "2025-11-12T23:59:59Z"
    }
  }'
```

**Example with image:**
```bash
curl -i -X POST http://localhost:8080/trade \
  -H "Content-Type: application/json" \
  -d '{
    "listing": {
      "offering_farmer_id": "8c8c73e8-0a16-4d3a-826d-75d50d7a758f",
      "offered_item_id": "b7f2c6d4-1aeb-4f5b-9c2b-222222222222",
      "offered_item_quantity": 20,
      "desired_items": "Looking for long-grain white rice.",
      "expires_at": "2025-11-12T23:59:59Z"
    },
    "listing_image": {
      "base64": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJ...",
      "content_type": "image/png"
    }
  }'
```

---

### 4. Create Trade Bid
Place a bid on an existing trade listing.

**Endpoint:** `POST /bid`

**Headers:**
- `Content-Type: application/json`

**Request Body:**
```json
{
  "trade_listing_id": "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
  "bidding_farmer_id": "9c6a8e18-4a0d-4a3a-8c6d-75d50d7a758f",
  "bid_item_id": "c9d3e8a1-55b2-4f66-a123-333333333333",
  "bid_item_quantity": 15
}
```

**Notes:**
- `id` is auto-generated
- `status` defaults to "pending" if not provided
- `created_at` is auto-generated by database

**Response:** `201 Created`
```json
{
  "status": "trade bid created",
  "bid_id": "8c9d0e1f-2a3b-4c5d-6e7f-8a9b0c1d2e3f"
}
```

**Example:**
```bash
curl -i -X POST http://localhost:8080/bid \
  -H "Content-Type: application/json" \
  -d '{
    "trade_listing_id": "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
    "bidding_farmer_id": "9c6a8e18-4a0d-4a3a-8c6d-75d50d7a758f",
    "bid_item_id": "c9d3e8a1-55b2-4f66-a123-333333333333",
    "bid_item_quantity": 15
  }'
```

---

### 5. Update Trade Listing Status
Change the status of a trade listing.

**Endpoint:** `PATCH /trade/:id?updated_status=<status>`

**Path Parameters:**
- `id` (required) - UUID of the trade listing

**Query Parameters:**
- `updated_status` (required) - New status value
  - Valid values: `open`, `completed`, `cancelled`

**Response:** `200 OK`
```json
{
  "status": "trade listing status updated",
  "trade_id": "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
  "new_status": "completed"
}
```

**Response:** `400 Bad Request` (if invalid status)
```json
{
  "error": "invalid status. Must be: open, completed, or cancelled"
}
```

**Examples:**
```bash
# Mark listing as completed
curl -i -X PATCH http://localhost:8080/trade/1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d?updated_status=completed

# Cancel a listing
curl -i -X PATCH http://localhost:8080/trade/1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d?updated_status=cancelled

# Reopen a listing
curl -i -X PATCH http://localhost:8080/trade/1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d?updated_status=open
```

---

## Trade Bid CRUD Operations

### 6. Get Trade Bid by ID
Retrieve a specific trade bid.

**Endpoint:** `GET /bid?id=<uuid>`

**Query Parameters:**
- `id` (required) - UUID of the trade bid

**Response:** `200 OK`
```json
{
  "id": "4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a",
  "trade_listing_id": "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
  "bidding_farmer_id": "9c6a8e18-4a0d-4a3a-8c6d-75d50d7a758f",
  "bid_item_id": "c9d3e8a1-55b2-4f66-a123-333333333333",
  "bid_item_quantity": 15,
  "status": "pending",
  "created_at": "2025-11-06T14:20:00Z"
}
```

**Response:** `404 Not Found`
```json
{
  "error": "trade bid not found"
}
```

**Example:**
```bash
curl -i http://localhost:8080/bid?id=4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a
```

---

### 7. Get Bids by Farmer
Retrieve all bids made by a specific farmer.

**Endpoint:** `GET /bids/farmer/:farmer_id`

**Path Parameters:**
- `farmer_id` (required) - UUID of the farmer

**Response:** `200 OK`
```json
[
  {
    "id": "4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a",
    "trade_listing_id": "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
    "bidding_farmer_id": "9c6a8e18-4a0d-4a3a-8c6d-75d50d7a758f",
    "bid_item_id": "c9d3e8a1-55b2-4f66-a123-333333333333",
    "bid_item_quantity": 15,
    "status": "pending",
    "created_at": "2025-11-06T14:20:00Z"
  }
]
```

**Example:**
```bash
curl -i http://localhost:8080/bids/farmer/9c6a8e18-4a0d-4a3a-8c6d-75d50d7a758f
```

---

### 8. Update Trade Bid
Update a trade bid's item and quantity (only for pending bids).

**Endpoint:** `PUT /bid/:id`

**Path Parameters:**
- `id` (required) - UUID of the trade bid

**Headers:**
- `Content-Type: application/json`

**Request Body:**
```json
{
  "bid_item_id": "c9d3e8a1-55b2-4f66-a123-333333333333",
  "bid_item_quantity": 20
}
```

**Response:** `200 OK`
```json
{
  "status": "trade bid updated",
  "bid_id": "4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a"
}
```

**Response:** `404 Not Found` (if bid doesn't exist or not pending)
```json
{
  "error": "bid not found or not in pending status"
}
```

**Example:**
```bash
curl -i -X PUT http://localhost:8080/bid/4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a \
  -H "Content-Type: application/json" \
  -d '{
    "bid_item_id": "c9d3e8a1-55b2-4f66-a123-333333333333",
    "bid_item_quantity": 20
  }'
```

---

### 9. Update Trade Bid Status
Change the status of a trade bid.

**Endpoint:** `PATCH /bid/:id/status?updated_status=<status>`

**Path Parameters:**
- `id` (required) - UUID of the trade bid

**Query Parameters:**
- `updated_status` (required) - New status value
  - Valid values: `pending`, `accepted`, `rejected`

**Response:** `200 OK`
```json
{
  "status": "trade bid status updated",
  "bid_id": "4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a",
  "new_status": "accepted"
}
```

**Response:** `400 Bad Request` (if invalid status)
```json
{
  "error": "invalid status. Must be: pending, accepted, or rejected"
}
```

**Examples:**
```bash
# Accept a bid
curl -i -X PATCH http://localhost:8080/bid/4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a/status?updated_status=accepted

# Reject a bid
curl -i -X PATCH http://localhost:8080/bid/4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a/status?updated_status=rejected
```

---

### 10. Delete Trade Bid
Delete a trade bid (only for pending bids).

**Endpoint:** `DELETE /bid/:id`

**Path Parameters:**
- `id` (required) - UUID of the trade bid

**Response:** `200 OK`
```json
{
  "status": "trade bid deleted",
  "bid_id": "4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a"
}
```

**Response:** `404 Not Found` (if bid doesn't exist or not pending)
```json
{
  "error": "bid not found or not in pending status"
}
```

**Example:**
```bash
curl -i -X DELETE http://localhost:8080/bid/4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a
```

---

## Image Handling

The Trading API uses base64 encoding for images to simplify client integration and avoid multipart/form-data complexity.

### On GET Requests (Image Encoding)

When retrieving trade listings:
- If a listing has an `image_url`, the API automatically reads the image file
- The image is encoded to base64 and included in the response
- `image_base64` field contains the base64-encoded image data
- `image_content_type` field indicates the MIME type (detected from file extension)
- Supported formats: PNG, JPEG, GIF, WebP

**Example GET Response with Image:**
```json
{
  "listing": {
    "id": "...",
    "image_url": "resources/images/abc123_trade.png",
    ...
  },
  "image_base64": "iVBORw0KGgoAAAANSUhEUg...",
  "image_content_type": "image/png"
}
```

### On POST Requests (Image Decoding)

When creating a trade listing:
- Include `listing_image` object with `base64` and `content_type` fields
- The API decodes the base64 string and saves the image to disk
- Image is saved as `resources/images/<listing_id>_trade.<ext>`
- File extension is determined by `content_type`
- If no image is provided, defaults to `resources/images/no-image.jpg`

**Example POST Request with Image:**
```json
{
  "listing": { ... },
  "listing_image": {
    "base64": "iVBORw0KGgoAAAANSUhEUg...",
    "content_type": "image/png"
  }
}
```

### Image File Naming Convention
- Trade listing images: `<listing_id>_trade.<ext>`
- Example: `7b8c9d0e-1f2a-3b4c-5d6e-7f8a9b0c1d2e_trade.png`

### Supported Image Formats
- PNG (`.png`) - `image/png`
- JPEG (`.jpg`, `.jpeg`) - `image/jpeg`
- GIF (`.gif`) - `image/gif`
- WebP (`.webp`) - `image/webp`

---

## Status Values

### Trade Listing Status
- `open` - Active listing accepting bids
- `completed` - Trade has been completed
- `cancelled` - Listing has been cancelled by the farmer

### Trade Bid Status
- `pending` - Bid is waiting for farmer's decision
- `accepted` - Bid has been accepted (trade agreed)
- `rejected` - Bid has been rejected

---

## Database Schema Reference

### TradeListing Table
```sql
CREATE TABLE public."TradeListing" (
    id uuid NOT NULL,
    offering_farmer_id uuid NOT NULL,
    offered_item_id uuid NOT NULL,
    offered_item_quantity numeric NOT NULL,
    desired_items text,
    status public.trade_listing_status NOT NULL,
    image_url character varying(255),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamp with time zone,
    CONSTRAINT "TradeListing_pkey" PRIMARY KEY (id)
);
```

### TradeBid Table
```sql
CREATE TABLE public."TradeBid" (
    id uuid NOT NULL,
    trade_listing_id uuid NOT NULL,
    bidding_farmer_id uuid NOT NULL,
    bid_item_id uuid NOT NULL,
    bid_item_quantity numeric NOT NULL,
    status public.trade_bid_status NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "TradeBid_pkey" PRIMARY KEY (id)
);
```

---

### 5. Update Trade Listing Status

1. **Farmer creates a listing:**
   ```bash
   POST /trade
   ```

2. **Other farmers browse available trades:**
   ```bash
   GET /trades/batch?block=0
   ```

3. **Farmer views specific listing details:**
   ```bash
   GET /trade?id=<listing_id>
   ```

4. **Farmer places a bid:**
   ```bash
   POST /bid
   ```

5. **Original farmer reviews bids and accepts one:**
   ```bash
   # Update bid status (requires bid status endpoint - not yet implemented)
   # Then mark listing as completed:
   PATCH /trade/<listing_id>?updated_status=completed
   ```

---

## Notes

- All timestamps are in RFC3339 format with timezone
- UUIDs must be valid UUID v4 format
- Foreign key constraints ensure data integrity
- Cascade deletes: If a farmer, item, or listing is deleted, related records are automatically removed
