package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Database model
type Database struct {
	ID        string `gorm:"primaryKey,uniqueIndex"`
	Hash      string `gorm:"uniqueIndex"`
	Buffer    []byte
	Name      string
	Size      int64
	Mime      string
	CreatedAt time.Time
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// Error returns an error in JSON or plaintext depending on the Accept header.
func Error(c echo.Context, status int, message string) error {
	if c.Request().Header.Get("Accept") == "application/json" {
		return c.JSON(status, map[string]interface{}{
			"error":   true,
			"status":  status,
			"message": message,
		})
	} else {
		return c.String(status, fmt.Sprintf("%d: %s\n", status, message))
	}
}

// UploadedFile returns a link to the uploaded file in JSON or plaintext depending on the Accept header.
func UploadedFile(c echo.Context, id string) error {
	if c.Request().Header.Get("Accept") == "application/json" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"uploaded": true,
			"status":   http.StatusOK,
			"message":  fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, id),
		})
	} else {
		return c.String(http.StatusOK, fmt.Sprintf("%s://%s/%s\n", c.Scheme(), c.Request().Host, id))
	}
}

// Upload a file, save and attribute an ID to it.
func upload(c echo.Context, db *gorm.DB) error {
	file, err := c.FormFile("file")
	if err != nil {
		return Error(c, http.StatusBadRequest, ts.HTTPErrors.BadRequest)
	}

	if config.MaxSize > 0 {
		if file.Size > (int64(config.MaxSize) * 1024 * 1024) {
			return Error(c, http.StatusRequestEntityTooLarge, ts.HTTPErrors.FileTooLarge)
		}
	}

	buffer := func() []byte {
		f, err := file.Open()
		if err != nil {
			return nil
		}
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, f); err != nil {
			return nil
		}
		return buf.Bytes()
	}()

	HashBytes := sha256.Sum256(buffer)
	hash := hex.EncodeToString(HashBytes[:])

	// Check if the file is already in the database by comparing the SHA256 hash
	// of the file with the ones in the database.
	var hashChecker Database
	db.Where("hash = ?", fmt.Sprintf("%x", sha256.Sum256(buffer))).First(&hashChecker)
	if hashChecker.ID != "" {
		return UploadedFile(c, hashChecker.ID)
	}

	var mimes = mimetype.Detect(buffer)
	var ext = ".bin"
	var contentType = mimes.String()
	if mimes.Extension() != "" {
		ext = mimes.Extension()
	}

	if c.Request().Header.Get("Parse_HTML") == "yes" && strings.HasPrefix(contentType, "text/plain") {
		contentType = "text/html; charset=utf-8"
		ext = ".html"
	}

	id := randSeq(5) + ext
	for {
		var count int64
		db.Model(&Database{}).Where("id = ?", id).Count(&count)
		if count == 0 {
			break
		}
		id = randSeq(5) + ext
	}

	data := Database{
		ID:        id,
		Name:      file.Filename,
		Buffer:    buffer,
		Hash:      hash,
		Size:      file.Size,
		Mime:      contentType,
		CreatedAt: time.Now().UTC(),
	}

	db.Create(&data)

	return UploadedFile(c, id)
}

// Gets the file using the provided UUID on the URL
func getFile(uuid string, db *gorm.DB) ([]byte, string, string) {
	if len(uuid) <= 0 {
		return nil, "", ""
	}

	uuid = strings.TrimSpace(uuid)

	var data Database
	db.Where("ID = ?", uuid).First(&data)

	if len(data.ID) <= 0 {
		return nil, "", ""
	}

	return data.Buffer, data.Name, data.Mime
}
