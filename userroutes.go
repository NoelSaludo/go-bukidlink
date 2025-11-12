package main

import (
	"bukidlink/db"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func postUserHandler(c *gin.Context) {
	// Expect payload shape: { "user": {..}, "profile_pic": {"base64":"...","content_type":"..."} }
	var req struct {
		User       db.User `json:"user"`
		ProfilePic *struct {
			Base64      string `json:"base64"`
			ContentType string `json:"content_type"`
		} `json:"profile_pic"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If profile picture provided, decode and save to resources/images/<username>_pfp.<ext>
	if req.ProfilePic != nil && req.ProfilePic.Base64 != "" {
		username := req.User.Username
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "username is required when uploading profile_pic",
			})
			return
		}

		data, err := base64.StdEncoding.DecodeString(req.ProfilePic.Base64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid base64 profile_pic",
			})
			return
		}

		// Determine extension from content type, fallback to first mime.ExtensionsByType result
		ext := "bin"
		if req.ProfilePic.ContentType != "" {
			// try mime.ExtensionsByType
			if exts, _ := mime.ExtensionsByType(req.ProfilePic.ContentType); len(exts) > 0 {
				// exts contain leading dot, strip it
				ext = strings.TrimPrefix(exts[0], ".")
			} else {
				// fallback to subtype (image/png -> png, image/svg+xml -> svg)
				parts := strings.Split(req.ProfilePic.ContentType, "/")
				if len(parts) == 2 {
					sub := strings.Split(parts[1], "+")[0]
					if sub != "" {
						ext = sub
					}
				}
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

		filename := fmt.Sprintf("%s_pfp.%s", username, ext)
		fp := filepath.Join(imgsDir, filename)

		// Check if file already exists, if not save it
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			// File doesn't exist, save it
			if err := os.WriteFile(fp, data, 0o644); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to save profile image",
				})
				return
			}
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to check existing profile image",
			})
			return
		}
		// If file exists, just use the existing one

		// Set relative path in user struct so it will be stored in DB
		req.User.ProfilePicPath = filepath.ToSlash(filepath.Join("resources/images", filename))
	}

	// Insert user into DB
	if err := db.InsertUser(req.User); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Success"})
}

func getUserHandler(c *gin.Context) {
	var temp db.User

	usernameP := c.Param("username")

	temp, err := db.QueryUser(usernameP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if temp.Username != "" {
		// Attempt to fetch and embed profile picture (non-fatal).
		if temp.ProfilePicPath != "" {
			const maxImageSize = 1 << 20 // 1 MiB limit
			var data []byte
			var contentType string

			// Treat ProfilePicPath as a local filesystem path. Resolve relative
			// paths from the process working directory and read the file.
			fp := temp.ProfilePicPath
			if !filepath.IsAbs(fp) {
				fp = filepath.Join(".", fp)
			}
			f, err := os.Open(fp)
			if err == nil {
				defer f.Close()
				data, err = io.ReadAll(io.LimitReader(f, maxImageSize+1))
				if err == nil && int64(len(data)) <= maxImageSize {
					contentType = http.DetectContentType(data)
				} else {
					data = nil
				}
			}

			if len(data) > 0 {
				b64 := base64.StdEncoding.EncodeToString(data)
				c.JSON(http.StatusOK, gin.H{
					"user": temp,
					"profile_pic": gin.H{
						"content_type": contentType,
						"base64":       b64,
					},
				})
				return
			}
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": "failed to read profile image"})
			return
		}

		// Default: return user only
		c.JSON(http.StatusOK, temp)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}
