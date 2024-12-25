package main

import (
	"encoding/json"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type Response struct {
	Data    []string `json:"data,omitempty"`
	Status  string   `json:"status"`
	Message string   `json:"message"`
}

func main() {
	http.HandleFunc("/upload", uploadHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle OPTIONS request for preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		responseJSON(w, Response{
			Status:  "error",
			Message: "Invalid request method",
		}, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // Limit file size to 10MB per file
	if err != nil {
		responseJSON(w, Response{
			Status:  "error",
			Message: "Failed to parse form",
		}, http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 || len(files) > 12 {
		responseJSON(w, Response{
			Status:  "error",
			Message: "Please upload between 1 and 12 files",
		}, http.StatusBadRequest)
		return
	}

	var uploadedPaths []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			responseJSON(w, Response{
				Status:  "error",
				Message: "Failed to open file",
			}, http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Generate UUID-based path
		newUUID := uuid.New().String()
		dirPath := filepath.Join("uploads", newUUID[:8], newUUID[9:13], newUUID[14:18], newUUID[19:23])
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			responseJSON(w, Response{
				Status:  "error",
				Message: "Failed to create directories",
			}, http.StatusInternalServerError)
			return
		}

		// Determine file extension
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		filename := filepath.Join(dirPath, newUUID[24:]+ext)

		switch {
		case ext == ".pdf":
			// Save PDF
			if err := savePDF(file, filename); err != nil {
				responseJSON(w, Response{
					Status:  "error",
					Message: "Failed to save PDF",
				}, http.StatusInternalServerError)
				return
			}
		case ext == ".epub", ext == ".docx", ext == ".xlsx":
			// Save other document formats
			if err := saveDocument(file, filename); err != nil {
				responseJSON(w, Response{
					Status:  "error",
					Message: "Failed to save document",
				}, http.StatusInternalServerError)
				return
			}
		case isImage(fileHeader.Filename):
			// Process and save image
			if err := processAndSaveImage(file, filename); err != nil {
				responseJSON(w, Response{
					Status:  "error",
					Message: "Failed to process image",
				}, http.StatusInternalServerError)
				return
			}
		default:
			responseJSON(w, Response{
				Status:  "error",
				Message: "Unsupported file type",
			}, http.StatusBadRequest)
			return
		}

		// Remove "uploads/" prefix for response
		relativePath := strings.TrimPrefix(filename, "uploads/")
		uploadedPaths = append(uploadedPaths, relativePath)
	}

	responseJSON(w, Response{
		Data:    uploadedPaths,
		Status:  "success",
		Message: "ok",
	}, http.StatusOK)
}

func saveDocument(file io.Reader, filename string) error {
	dest, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, file)
	return err
}

func savePDF(file io.Reader, filename string) error {
	dest, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, file)
	return err
}

func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".png" || ext == ".jpeg"
}

func processAndSaveImage(file io.Reader, filename string) error {
	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Resize the image if necessary
	img = resizeImage(img)

	// Save the resized image
	dest, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dest.Close()

	if format == "jpeg" || format == "jpg" {
		return jpeg.Encode(dest, img, nil)
	}

	return nil
}

func resizeImage(img image.Image) image.Image {
	maxSize := uint(1080)
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	if width > int(maxSize) || height > int(maxSize) {
		return resize.Resize(maxSize, 0, img, resize.Lanczos3)
	}

	return img
}

func responseJSON(w http.ResponseWriter, response Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
