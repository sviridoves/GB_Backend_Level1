package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type helloHandler struct {
	subject string
}

type Handler struct {
}

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

func (h *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", h.subject)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		name := r.FormValue("name")
		fmt.Fprintf(w, "Parsed query-param with key \"name\": %s", name)
	case http.MethodPost:
		//body, err := ioutil.ReadAll(r.Body)
		//if err != nil {
		//	http.Error(w, "Unable to read request body", http.StatusBadRequest)
		//	return
		//}
		//defer r.Body.Close()

		var employee Employee

		contentType := r.Header.Get("Content-Type")
		//err = json.Unmarshal(body, &employee)
		//if err != nil {
		//	http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
		//	return
		//}
		switch contentType {
		case "application/json":
			err := json.NewDecoder(r.Body).Decode(&employee)
			if err != nil {
				http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
				return
			}
		case "application/xml":
			err := xml.NewDecoder(r.Body).Decode(&employee)
			if err != nil {
				http.Error(w, "Unable to unmarshal XML", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "Unknown content type", http.StatusBadRequest)
			return
		}
		//fmt.Fprintf(w, "Parsed request body: %s\n", string(body))
		fmt.Fprintf(w, "Got a new employee!\nName: %s\nAge: %dy.o.\nSalary %0.2f\n",
			employee.Name,
			employee.Age,
			employee.Salary,
		)
	}
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
	//fmt.Fprintf(w, "File %s has been successfully uploaded", header.Filename)
}

func main() {
	//http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "Hello World!")
	//})
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "Ya got the wrong place, pal")
	//})
	uploadHandler := &UploadHandler{
		UploadDir: "upload",
		HostAddr:  "http://localhost:8080",
	}
	http.Handle("/upload", uploadHandler)

	handler := &Handler{}
	http.Handle("/", handler)

	worldHandler := &helloHandler{"World"}
	roomHandler := &helloHandler{"Mark"}

	http.Handle("/world", worldHandler)
	http.Handle("/room", roomHandler)

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
	//http.ListenAndServe(":80", nil)
}
