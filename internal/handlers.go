package main

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type Data struct {
	gorm.Model
	Buffer    []byte
	ID        string
	Name      string
	Size      int64
	Mime      string
	CreatedAt time.Time
}

func jsonOrString(c echo.Context, status int, message string, error bool) error {
	if (c.Request().Header.Get("Accept")) == "application/json" {
		return c.JSON(status, map[string]interface{}{
			"error":   error,
			"status":  status,
			"message": message,
		})
	} else {
		return c.String(status, message+"\n")
	}
}

// Upload a file, save and attribute an ID to it.
func upload(c echo.Context, db *gorm.DB) error {

	file, err := c.FormFile("file")
	if err != nil {
		return jsonOrString(c, http.StatusBadRequest, "400: Bad request.", true)
	}

	if file.Size > (int64(config.MaxSize) * 1024 * 1024) {
		return jsonOrString(c, http.StatusRequestEntityTooLarge, "413: Request entity too large.", true)
	}

	id := xid.New().String()
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

	data := Data{
		ID:        id,
		Name:      file.Filename,
		Buffer:    buffer,
		Size:      file.Size,
		Mime:      mime.TypeByExtension(file.Filename),
		CreatedAt: time.Now().UTC(),
	}

	db.Create(&data)

	if (c.Request().Header.Get("Accept")) == "application/json" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  http.StatusOK,
			"message": "200: File uploaded successfully.",
			"url":     fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, id),
		})
	} else {
		return c.String(http.StatusOK, fmt.Sprintf("%s://%s/%s\n", c.Scheme(), c.Request().Host, id))
	}
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
