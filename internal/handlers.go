package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jaevor/go-nanoid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Database model.
type Data struct {
	ID        string `gorm:"primaryKey,uniqueIndex"`
	Hash      string `gorm:"uniqueIndex"`
	Buffer    []byte
	Name      string
	Size      int64
	Mime      string
	CreatedAt time.Time
}

// MakeError returns an error in JSON or plaintext depending on the Accept header.
func MakeError(c echo.Context, status int, message string) error {
	if (c.Request().Header.Get("Accept")) == "application/json" {
		return c.JSON(status, map[string]interface{}{
			"error":   true,
			"status":  status,
			"message": message,
		})
	} else {
		return c.String(status, fmt.Sprintf("%d: %s\n", status, message))
	}
}

// MakeUploadedFile returns a link to the uploaded file in JSON or plaintext depending on the Accept header.
func MakeUploadedFile(c echo.Context, id string) error {
	if (c.Request().Header.Get("Accept")) == "application/json" {
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
		return MakeError(c, http.StatusBadRequest, ts.HTTPErrors.BadRequest)
	}

	generator, err := nanoid.ASCII(5)

	if err != nil {
		panic(err)
	}

	if config.MaxSize > 0 {
		if file.Size > (int64(config.MaxSize) * 1024 * 1024) {
			return MakeError(c, http.StatusRequestEntityTooLarge, ts.HTTPErrors.FileTooLarge)
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

	// Check if the file is already in the database by comparing the SHA256 hash
	// of the file with the ones in the database.
	var hashChecker Data
	db.Where("hash = ?", fmt.Sprintf("%x", sha256.Sum256(buffer))).First(&hashChecker)
	if hashChecker.ID != "" {
		MakeUploadedFile(c, hashChecker.ID)
	}

	mimes := mimetype.Detect(buffer)
	ext := ".bin"

	if mimes.Extension() != "" {
		ext = mimes.Extension()
	}

	id := generator() + ext
	for {
		var data Data
		db.Where("ID = ?", id).First(&data) // check for duplicates

		if len(data.ID) <= 0 {
			break
		}

		id = generator() + ext
	}
	data := Data{
		ID:        id,
		Name:      file.Filename,
		Buffer:    buffer,
		Hash:      fmt.Sprintf("%x", sha256.Sum256(buffer)),
		Size:      file.Size,
		Mime:      mimes.String(),
		CreatedAt: time.Now().UTC(),
	}

	db.Create(&data)

	return MakeUploadedFile(c, id)
}

// Gets the file using the provided UUID on the URL
func getFile(uuid string, db *gorm.DB) ([]byte, string, string) {
	if len(uuid) <= 0 {
		return nil, "", ""
	}

	uuid = strings.TrimSpace(uuid)

	var data Data
	db.Where("ID = ?", uuid).First(&data)

	if len(data.ID) <= 0 {
		return nil, "", ""
	}

	return data.Buffer, data.Name, data.Mime
}
