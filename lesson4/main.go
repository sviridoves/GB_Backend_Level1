package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

type GetFilesListHandler struct {
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}

	filePath := h.UploadDir + "\\" + header.Filename

	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	fileLink := h.HostAddr + "/" + header.Filename

	req, err := http.NewRequest(http.MethodGet, fileLink, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to check file", http.StatusInternalServerError)
		return
	}

	cli := &http.Client{}

	resp, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to check file", http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, fileLink)
}

func (h *GetFilesListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, `<html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8"><title>Files</title></head><body>`)
	extension := r.FormValue("ext")
	if len(extension) > 0 {
		fmt.Fprintf(w, `<h2>Отображаются файлы с расширением: %s</h2>`, extension)
		fmt.Fprint(w, `<table width="40%" border="1" cellspacing="0"><tr><th>File name</th><th>File type</th><th>File size, bytes</th></tr>`)
		for _, file := range files {
			if !file.IsDir() {
				fileType := filepath.Ext(file.Name())
				//fileType = fileType[1:]
				if len(fileType) > 0 {
					if fileType[1:] == extension {
						fmt.Fprintf(w, `<tr><td>%s</td><td>%s</td><td>%d</td></tr>`, file.Name(), fileType, file.Size())
					}
				}
			}
		}
	} else {
		fmt.Fprint(w, `<table width="40%" border="1" cellspacing="0"><tr><th>File name</th><th>File type</th><th>File size, bytes</th></tr>`)
		for _, file := range files {
			if !file.IsDir() {
				fileType := filepath.Ext(file.Name())
				fmt.Fprintf(w, `<tr><td>%s</td><td>%s</td><td>%d</td></tr>`, file.Name(), fileType, file.Size())
			}
		}
	}
	fmt.Fprint(w, `</table></body></html>`)
}

func main() {
	uploadHandler := &UploadHandler{
		UploadDir: "upload",
		HostAddr:  "http://localhost:8080",
	}
	http.Handle("/upload", uploadHandler)

	http.Handle("/files", &GetFilesListHandler{})
	srv := &http.Server{
		Addr:         ":80",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go srv.ListenAndServe()

	dirToServe := http.Dir(uploadHandler.UploadDir)

	fs := &http.Server{
		Addr:         ":8080",
		Handler:      http.FileServer(dirToServe),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fs.ListenAndServe()
}
