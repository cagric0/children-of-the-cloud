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
		port = "8080"
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
	fmt.Println("aaaa")
	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		fmt.Println("1", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	//Access the image key
	imageFile, _, err := r.FormFile("image")
	if err != nil {
		fmt.Println("2", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	//Access the audio key
	audioFile, _, err := r.FormFile("audio")
	if err != nil {
		fmt.Println("3", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	texts, err := services.SpeechToText(audioFile)
	if err != nil {
		fmt.Println("4", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	objects, err := services.DetectObjects(imageFile)
	if err != nil {
		fmt.Println("5", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result := compare(objects, texts)

	resultBytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println("6", err)
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
	if objects == nil {
		fmt.Println("No object detected")
	}
	fmt.Println(texts)

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
