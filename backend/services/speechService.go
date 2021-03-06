package services

import (
	"bytes"
	speech "cloud.google.com/go/speech/apiv1"
	"context"
	"fmt"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"mime/multipart"
	"strings"
)

func SpeechToText(file multipart.File) ([]string, string, error) {
	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, "", fmt.Errorf("ReadFrom File: %v", err)
	}

	// retrieve a byte slice from bytes.Buffer
	data := buf.Bytes()

	// 	Encoding:        speechpb.RecognitionConfig_LINEAR16,
	// 	SampleRateHertz: 48000,
	// 	UseEnhanced:     true,
	//	// A model must be specified to use enhanced model.
	// 	Model: "phone_call",
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_ENCODING_UNSPECIFIED,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})

	if err != nil {
		return nil, "", fmt.Errorf("Recognize: %v", err)
	}

	if len(resp.Results) == 0 {
		return nil, "", nil
	}

	receivedText := ""
	wordList := make([]string, 0)
	for _, result := range resp.Results {
		for _, alternative := range result.Alternatives {
			transcript := strings.TrimSpace(alternative.Transcript)
			words := strings.Split(transcript, " ")
			wordList = appendUnique(wordList, words)
			receivedText += " " + transcript
		}
	}

	return wordList, strings.TrimSpace(receivedText), nil
}

func appendUnique(a []string, b []string) []string {

	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter, _ := range check {
		res = append(res, letter)
	}

	return res
}
