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

type Result struct {
	IsSuccess  bool               `json:"isSuccess"`
	Objects    []*services.Object `json:"objects"`
	SpeechTexts []string             `json:"speechTexts"`
}

type SpeechTextResult struct {
	texts  []string
	speechText    string
	err error
}

type ObjectResult struct {
	objects  []*services.Object
	err error
}

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

	speechChan := make(chan SpeechTextResult)

	go func() {
		texts, speechText, err := services.SpeechToText(audioFile)
		speechChan <- SpeechTextResult{texts, speechText, err}
	}()

	objectChan := make(chan ObjectResult)

	go func() {
		objects, err := services.DetectObjects(imageFile)
		objectChan <- ObjectResult{objects, err}
	}()

	speechTextResult := <- speechChan

	objectResult := <- objectChan

	if speechTextResult.err != nil {
		log.Printf("SpeechToText: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if objectResult.err != nil {
		log.Printf("DetectObjects: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result := compare(objectResult.objects, speechTextResult.texts)
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

func compare(objects []*services.Object, texts []string) Result {

	fmt.Println(texts)

	var result Result

	count := 0
	result.IsSuccess = false
	if objects != nil {
		result.Objects = objects
		count ++
	}
	if len(texts) != 0 {
		result.SpeechTexts = texts
		count ++
	}
	if count != 2 {
		return result
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
		}
	}

	return result
}
