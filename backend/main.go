package main

import (
	"CloudChildren/services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

const (
	port = 3000
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/final", finalHandler).Methods("POST")
	log.Println("Go service is running on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

func finalHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	//Access the image key
	imageFile, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	//Access the audio key
	audioFile, _, err := r.FormFile("audio")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	texts, err := services.SpeechToText(audioFile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	objects, err := services.DetectObjects(imageFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result := compare(objects, texts)

	resultBytes, err := json.Marshal(result)
	if err != nil {
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
	IsSuccess bool             `json:"isSuccess"`
	Objects []*services.Object `json:"objects"`
}
func compare(objects []*services.Object, texts []string) Result {
	fmt.Println(texts)
	fmt.Println(objects)

	var result Result
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
	result.Objects = objects

	return result
}
