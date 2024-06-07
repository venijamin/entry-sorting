package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func main() {
	fs := http.FileServer(http.Dir("./server/sites"))

	mux := http.NewServeMux()
	mux.Handle("/", fs)
	mux.Handle("/upload", &FileUploadHandler{})

	log.Print("Listening on :3333...")
	err := http.ListenAndServe(":3333", mux)
	if err != nil {
		log.Fatal(err)
	}
}

type FileUploadHandler struct{}

var (
	FileUpload = regexp.MustCompile(`^/upload/*$`)
)

func (h *FileUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && FileUpload.MatchString(r.URL.Path):
		h.FileUpload(w, r)
	default:
		return
	}

	// TODO: results page
}

func (h *FileUploadHandler) FileUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		fmt.Fprintf(w, "Failed to Upload File\n")
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("./", "*.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
