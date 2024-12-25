package main

import (
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate a UUID for the file name
	newUUID := uuid.New().String()
	dirStructure := filepath.Join(newUUID[:8], newUUID[9:13], newUUID[14:18], newUUID[19:23])
	dirPath := filepath.Join("uploads", dirStructure)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		http.Error(w, "Failed to create directories", http.StatusInternalServerError)
		return
	}

	filename := filepath.Join(dirPath, newUUID[24:]+filepath.Ext(header.Filename))

	// Determine file type (image or PDF)
	switch {
	case header.Filename[len(header.Filename)-4:] == ".pdf":
		if err := savePDF(file, filename); err != nil {
			http.Error(w, "Failed to save PDF", http.StatusInternalServerError)
			return
		}
	case isImage(header.Filename):
		if err := processAndSaveImage(file, filename); err != nil {
			http.Error(w, "Failed to process image", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
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
	ext := filename[len(filename)-4:]
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
