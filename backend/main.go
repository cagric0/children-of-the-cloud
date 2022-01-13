package main

import (
	"CloudChildren/services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// PORT environment variable is provided by Cloud Run.
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	r := mux.NewRouter()
	r.HandleFunc("/final", finalHandler).Methods("POST")

	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}

func finalHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		log.Printf("ParseMultipartForm: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	//Access the image key
	imageFile, _, err := r.FormFile("image")
	if err != nil {
		log.Printf("Image FormFile: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	//Access the audio key
	audioFile, _, err := r.FormFile("audio")
	if err != nil {
		log.Printf("Audio FormFile: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	texts, speechText, err := services.SpeechToText(audioFile)
	if err != nil {
		log.Printf("SpeechToText: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	objects, err := services.DetectObjects(imageFile)
	if err != nil {
		log.Printf("DetectObjects: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result := compare(objects, texts)
	result.SpeechText = speechText
	resultBytes, err := json.Marshal(result)
	if err != nil {
		log.Printf("Marshal: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultBytes)
	return
}

type Result struct {
	IsSuccess  bool               `json:"isSuccess"`
	Objects    []*services.Object `json:"objects"`
	SpeechText string             `json:"speechText"`
	ErrorMsg   string             `json:"errorMsg,omitempty"`
}

func compare(objects []*services.Object, texts []string) Result {

	fmt.Println(texts)

	var result Result

	result.ErrorMsg = "No match"
	result.IsSuccess = false
	if objects == nil {
		result.ErrorMsg = "Objects not detected"
		return result
	}

	result.Objects = objects
	if len(texts) == 0 {
		result.ErrorMsg = "Speech not recognized"
	}

	textMap := make(map[string]int, len(texts))
	for _, text := range texts {
		textMap[strings.ToLower(text)] = 1
	}

	for _, object := range objects {
		_, ok := textMap[strings.ToLower(object.Name)]
		object.DetectedByUser = ok
		if ok {
			result.IsSuccess = true
			result.ErrorMsg = ""
		}
	}

	return result
}
