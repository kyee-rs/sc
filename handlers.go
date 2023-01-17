package main

import (
	"bytes"
	"io"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Data struct {
	gorm.Model
	Buffer []byte
	ID     string
	Name   string
	MIME   string
}

// Upload a file, save and attribute a hash
func upload(c *gin.Context, db *gorm.DB, scheme string) {

	file, err := c.FormFile("file")
	if err != nil {
		c.String(400, "400: Bad request!\n")
		return
	}

	if file.Size > (int64(config.Size_limit) * 1024 * 1024) {
		c.String(413, "413: Request entity too large!\n")
		return
	}

	uuid := uuid.New().String()
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

	mimes := mimetype.Detect(buffer).String()

	if mimes == "" {
		mimes = "application/octet-stream"
	}
	data := Data{
		ID:     uuid,
		Name:   file.Filename,
		MIME:   mimes,
		Buffer: buffer,
	}

	db.Create(&data)

	c.String(200, scheme+"://"+c.Request.Host+"/"+uuid+"\n")
}

// Gets the file using the provided UUID on the URL
func getFile(uuid string) ([]byte, string, string) {
	db, err := gorm.Open(sqlite.Open(config.DB_path), &gorm.Config{})
	if err != nil {
		panic("Connection to database failed. Please check your configuration.")
	}

	var data Data
	db.First(&data, "ID = ?", uuid)

	if len(data.ID) <= 0 {
		return nil, "", ""
	}

	return data.Buffer, data.MIME, data.Name
}
