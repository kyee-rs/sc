package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Data struct {
	gorm.Model
	Buffer []byte
	ID     string
	Name   string
}

// Handles and processes the home page
func home(w http.ResponseWriter, r *http.Request) {
	config := loadConfig(*config_file)
	if config.Fake_SSL || config.SSL_ {
		tmpl.Execute(w, template.HTML(fmt.Sprintf(`https://%s/`, r.Host)))
	} else {
		tmpl.Execute(w, template.HTML(fmt.Sprintf(`http://%s/`, r.Host)))
	}
}

// Upload a file, save and attribute a hash
func upload(w http.ResponseWriter, r *http.Request) {
	config := loadConfig(*config_file)
	logger := log.WithFields(log.Fields{
		"time":    time.Now(),
		"service": "upload",
		"file":    "handlers.go",
	})

	r.Body = http.MaxBytesReader(w, r.Body, int64(config.Size_limit)*1024*1024)
	if err := r.ParseMultipartForm(int64(config.Size_limit) * 1024 * 1024); err != nil {
		logger.Errorf("Error parsing form.")
		logger.Errorf("Error: %s", err.Error())
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		fmt.Fprintf(w, "413: File too large. Max size is %sMB.", fmt.Sprint(config.Size_limit))
		return
	}
	db, err := gorm.Open(sqlite.Open(config.DB_path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	logger.Info("Upload request recieved")

	var uuid string = GenerateUUID()
	buf := bytes.NewBuffer(nil)

	// Prepare to get the file
	if file, header, err := r.FormFile("file"); err != nil {
		logger.Errorf("Error uploading file.")
		logger.Errorf("Error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400: Bad request.")
		return
	} else {
		defer func() {
			file.Close()
			logger.Infof(`File "%s" closed.`, header.Filename)
		}()

		if _, err := io.Copy(buf, file); err != nil {
			w.WriteHeader(http.StatusInsufficientStorage)
			fmt.Fprintf(w, "Insufficient Storage. Error storing file.")
			return
		}

		var data Data
		db.Where(Data{Buffer: buf.Bytes()}).Attrs(Data{ID: uuid, Name: header.Filename}).FirstOrCreate(&data)

		if config.Fake_SSL || config.SSL_ {
			fmt.Fprintf(w, `https://%s/%s`, r.Host, uuid)
		} else {
			fmt.Fprintf(w, `http://%s/%s`, r.Host, uuid)
		}
	}
}

// Gets the file using the provided UUID on the URL
func getFile(w http.ResponseWriter, r *http.Request) {
	config := loadConfig(*config_file)
	logger := log.WithFields(log.Fields{
		"time":    time.Now(),
		"service": "getFile",
		"file":    "handlers.go",
	})

	logger.Info("Retrieve request received")
	db, err := gorm.Open(sqlite.Open(config.DB_path), &gorm.Config{})
	if err != nil {
		panic("Connection to database failed. Please check your configuration.")
	}
	var uuid string = strings.Replace(r.URL.Path[1:], "/", "", -1)

	logger.Infof(`Route "%s"`, r.URL.Path)
	logger.Infof(`Retrieving UUID "%s"`, uuid)

	var data Data
	db.First(&data, "ID = ?", uuid)

	if len(data.ID) <= 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404: File not found.\n")
		return
	}

	var filename = data.Name
	logger.Infof(`Retrieving Filename "%s"`, fmt.Sprintf("./%s", filename))

	w.Header().Set("Content-Disposition", fmt.Sprintf("filename=%s", filename))
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(data.Buffer))
}
