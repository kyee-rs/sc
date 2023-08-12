package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	"github.com/voxelin/sc/sqlc_gen"
)

// upload a file, save, and attribute an ID to it.
func upload(c *fiber.Ctx) error {
	formFile, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Segmentation Fault")
	}

	buffer := func() []byte {
		f, err := formFile.Open()
		if err != nil {
			return nil
		}
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, f); err != nil {
			return nil
		}
		return buf.Bytes()
	}()

	hashBytes := sha256.Sum256(buffer)
	hash := hex.EncodeToString(hashBytes[:])

	hashChecker, err := db.GetFileHash(ctx, hash)
	if err == nil {
		return c.SendString(c.BaseURL() + "/" + hashChecker.ID)
	}

	mimes := mimetype.Detect(buffer)
	ext := ".unknown"
	contentType := mimes.String()
	if mimes.Extension() != "" {
		ext = mimes.Extension()
	}

	shortID, err := sid.Generate()
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Segmentation Fault")
	}

	id := shortID + ext

	file, err := db.CreateFile(ctx, postgresql.CreateFileParams{
		ID:     id,
		Name:   formFile.Filename,
		Mime:   contentType,
		Size:   formFile.Size,
		Buffer: buffer,
		Hash:   hash,
	})
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Segmentation Fault")
	}

	return c.SendString(c.BaseURL() + "/" + file.ID)
}

// loadResponse - Loads a File Response
func loadResponse(c *fiber.Ctx) error {
	data, err := db.GetFile(ctx, c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Segmentation Fault")
	}

	c.Response().Header.Set("Content-Disposition", "filename=\""+data.Name+"\"")
	c.Response().Header.Set("Content-Type", data.Mime)
	
	return c.Send(data.Buffer)
}
