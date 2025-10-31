package main

import (
	"bukidlink/db"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func get100ItemsHandler(c *gin.Context) {

	blockP, err := strconv.Atoi(c.Param("block"))
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	items, err := db.QueryAllItem100(blockP)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	// Embed images for each item
	itemsWithImages := embedItemImages(items)

	c.JSON(http.StatusOK, itemsWithImages)
}

func getItemByCategory(c *gin.Context) {
	cat := c.Param("category")

	items, err := db.QueryItembyCategory(cat)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	// Embed images for each item
	itemsWithImages := embedItemImages(items)

	c.JSON(http.StatusOK, itemsWithImages)
}

func getItembyId(c *gin.Context) {
	itemid := c.Query("id")

	item, err := db.QueryItemByID(itemid)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	// Embed image for the item
	if item.ImgPath != "" {
		if imgData := readItemImage(item.ImgPath); imgData != nil {
			c.JSON(http.StatusOK, gin.H{
				"item": item,
				"item_pic": gin.H{
					"content_type": imgData.ContentType,
					"base64":       imgData.Base64,
				},
			})
			return
		}
	}

	// Default: return item only (no image or image read failed)
	c.JSON(http.StatusOK, item)
}

func postItemHandler(c *gin.Context) {
	// Expect payload shape: { "item": {..}, "item_pic": {"base64":"...","content_type":"..."} }
	var req struct {
		Item    db.Item `json:"item"`
		ItemPic *struct {
			Base64      string `json:"base64"`
			ContentType string `json:"content_type"`
		} `json:"item_pic"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If item picture provided, decode and save to resources/images/<itemname>_pic.<ext>
	if req.ItemPic != nil && req.ItemPic.Base64 != "" {
		itemName := req.Item.Name
		if itemName == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "item name is required when uploading item_pic",
			})
			return
		}

		data, err := base64.StdEncoding.DecodeString(req.ItemPic.Base64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid base64 item_pic",
			})
			return
		}

		// Determine extension from content type
		ext := "jpg"
		switch req.ItemPic.ContentType {
		case "image/jpeg", "image/jpg":
			ext = "jpg"
		case "image/png":
			ext = "png"
		case "image/gif":
			ext = "gif"
		case "image/webp":
			ext = "webp"
		default:
			// Try to detect from decoded data
			contentType := http.DetectContentType(data)
			if contentType == "image/png" {
				ext = "png"
			} else if contentType == "image/gif" {
				ext = "gif"
			}
		}

		// Ensure directory exists
		imgsDir := filepath.Join("resources", "images")
		if err := os.MkdirAll(imgsDir, 0o755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create images directory",
			})
			return
		}

		filename := itemName + "_pic." + ext
		fp := filepath.Join(imgsDir, filename)

		// Check if file already exists, if not save it
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			// File doesn't exist, save it
			if err := os.WriteFile(fp, data, 0o644); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to save item image",
				})
				return
			}
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to check existing item image",
			})
			return
		}

		// Set relative path in item struct so it will be stored in DB
		req.Item.ImgPath = filepath.ToSlash(filepath.Join("resources/images", filename))
	}

	// Insert item into DB
	itemId, err := db.InsertItem(req.Item)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "item created", "item_id": itemId})
}

// readItemImage reads an image from the filesystem and returns base64-encoded data
func readItemImage(imgPath string) *struct {
	Base64      string
	ContentType string
} {
	const maxImageSize = 1 << 20 // 1 MiB limit

	fp := imgPath
	if !filepath.IsAbs(fp) {
		fp = filepath.Join(".", fp)
	}

	f, err := os.Open(fp)
	if err != nil {
		return nil
	}
	defer f.Close()

	data, err := io.ReadAll(io.LimitReader(f, maxImageSize+1))
	if err != nil || int64(len(data)) > maxImageSize {
		return nil
	}

	contentType := http.DetectContentType(data)
	b64 := base64.StdEncoding.EncodeToString(data)

	return &struct {
		Base64      string
		ContentType string
	}{
		Base64:      b64,
		ContentType: contentType,
	}
}

// embedItemImages processes a slice of items and returns them with embedded images
func embedItemImages(items []db.Item) []gin.H {
	result := make([]gin.H, 0, len(items))

	for _, item := range items {
		itemData := gin.H{"item": item}

		if item.ImgPath != "" {
			if imgData := readItemImage(item.ImgPath); imgData != nil {
				itemData["item_pic"] = gin.H{
					"content_type": imgData.ContentType,
					"base64":       imgData.Base64,
				}
			}
		}

		result = append(result, itemData)
	}

	return result
}
