package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	"google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func GoogleSpeechToText(audioFilePath string) string {
	// Replace with the path to your service account key JSON file.
	serviceAccountKeyPath := "key.json"

	// Read the service account key JSON file.
	keyJSON, err := ioutil.ReadFile(serviceAccountKeyPath)
	if err != nil {
		log.Fatalf("Error reading service account key: %v", err)
	}

	// Set up the Speech-to-Text client with the service account key.
	ctx := context.Background()
	client, err := speech.NewClient(ctx, option.WithCredentialsJSON(keyJSON))
	if err != nil {
		log.Fatalf("Error creating Speech-to-Text client: %v", err)
	}

	// Read the audio file.
	audioData, err := ioutil.ReadFile(audioFilePath)
	if err != nil {
		log.Fatalf("Error reading audio file: %v", err)
	}

	// Perform speech recognition on the audio data.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:          speechpb.RecognitionConfig_WEBM_OPUS,
			SampleRateHertz:   48000,   // Adjust based on the audio file's sample rate.
			LanguageCode:      "en-IN", // Language code for English (United States).
			AudioChannelCount: 2,
			MaxAlternatives:   0,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: audioData},
		},
	})
	if err != nil {
		log.Fatalf("Error performing speech recognition: %v", err)
	}

	transcript := ""
	i := 0
	// Print the recognized text.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Println("Transcript Alternatives:", i, alt.Transcript)
			transcript += alt.Transcript
		}
		i++
	}
	fmt.Println("Transcript Done")
	return transcript
}
