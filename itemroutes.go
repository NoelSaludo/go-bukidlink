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
