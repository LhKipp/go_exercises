// curl -v -X POST localhost:8080 -F "alpha=30" -F "numberShapes=60" -F file=@png/mona.png > out.png

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func createPrimitivePic(inputPic string, alpha, numberShapes int) (outputPic string, err error) {
	outputPic = inputPic + ".primitive.png"
	err = exec.Command("primitive", "-i", inputPic, "-n", fmt.Sprint(numberShapes), "-a", fmt.Sprint(alpha), "-o", outputPic).Run()
	return
}

func createPic(w http.ResponseWriter, r *http.Request) {
	const MAX_UPLOAD_SIZE = 1024 * 1024 // 1 MB
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.ContentLength > MAX_UPLOAD_SIZE {
		http.Error(w, "The uploaded image is too big. Please use an image less than 1MB in size", http.StatusBadRequest)
		return
	}
	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return
	}
	alpha, _ := strconv.Atoi(r.FormValue("alpha"))
	numberShapes, _ := strconv.Atoi(r.FormValue("numberShapes"))

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dstFileName := fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	dst, err := os.Create(dstFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	primitivePic, err := createPrimitivePic(dstFileName, alpha, numberShapes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileBytes, err := ioutil.ReadFile(primitivePic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func main() {
	http.HandleFunc("/", createPic)
	http.ListenAndServe(":8080", nil)
}
