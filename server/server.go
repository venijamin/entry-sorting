package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	os.MkdirAll("./files", 0755) // 0755 is the permission mode
	fs := http.FileServer(http.Dir("./sites"))

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
}

func (h *FileUploadHandler) FileUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader, so we can get the Filename,
	// the Header and the size of the file
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		fmt.Fprintf(w, "Failed to Upload File\n")
		return
	}
	defer file.Close()
	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := os.CreateTemp("./files", "*.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!

	fileName := sortFile(tempFile.Name())
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, fileName)
	http.Redirect(w, r, "localhost:3333/", 301)
}

func sortFile(filepath string) string {
	start := time.Now()
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	sort.Sort(CSV(records))

	filename := strings.TrimSuffix(path.Base(filepath), path.Ext(filepath))
	sortedFileName := "./files/" + filename + "-sorted.csv"
	sortedFile, err := os.OpenFile(sortedFileName, os.O_CREATE|os.O_WRONLY, 0644)
	writer := csv.NewWriter(sortedFile)
	writer.WriteAll(records)

	fmt.Printf("Sorting time: %s\n", time.Since(start))
	return sortedFileName

}

type CSV [][]string

// Determine if one CSV line at index i comes before the line at index j.
func (data CSV) Less(i, j int) bool {
	dateColumnIndex := 0
	date1 := data[i][dateColumnIndex]
	date2 := data[j][dateColumnIndex]
	timeT1, _ := time.Parse("2006-01-02 15:04:05", date1)
	timeT2, _ := time.Parse("2006-01-02 15:04:05", date2)

	return timeT1.Before(timeT2)
}

// Other functions required for sort.Sort.
func (data CSV) Len() int {
	return len(data)
}
func (data CSV) Swap(i, j int) {
	data[i], data[j] = data[j], data[i]
}
