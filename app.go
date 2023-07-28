package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ffmpeg -i file.webm -ac 2 -f wav file.wav

var audioBuffer bytes.Buffer
var bufferMutex sync.Mutex

var websocketConn = make(map[*websocket.Conn]string)
var audioSegemetUser = make(map[string]int)

const (
	audioBufferSize = 1024 * 1024 // Adjust the buffer size as needed
)

type audioChunkStruct struct {
	UserId     string `json:"userId"`
	AudioChunk int    `json:"audiochunk"`
}

type speechResponse struct {
	ErrorRegions [][]int
	Transcript   string
}

func deleteFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		// If the file doesn't exist, os.Remove will return an error with message "remove path/to/your/file.txt: no such file or directory".
		// You can check the error to determine if the file didn't exist and ignore the error in that case.
		fmt.Println("Error deleting file:", err)
		return
	}
}

func main() {
	deleteFile("file.wav")
	deleteFile("output")

	// log.Println("Speaking Server has started")
	// http.HandleFunc("/ws", handleWebSocket)
	// log.Fatal(http.ListenAndServe(":8000", nil))

	// GoogleSpeechToText("output")

	// RabbitMQConnect()

	// start := time.Now()
	// testSTT()
	// elapsed := time.Since(start)

	// fmt.Printf("page took %s", elapsed)

	// start = time.Now()
	// // RabbitMQ_SendAudio()
	// elapsed = time.Since(start)

	// fmt.Printf("page took2 %s", elapsed)

	//7 characters 11
	//12 characters 7
	// var databaseText string = "People Love to comment on the self made businessman. The person who has come up the hard way. "
	// var userTranscript string = "people love to comment on the self made Businessman the person who is come up the Hard Way"
	// compareText(databaseText, userTranscript)

	// userTranscript := pythonTranscribe()
	// log.Println(userTranscript)
	var databaseText string = "People Love to comment on the self-made businessman. The person who has come up the hard way. "
	userTranscript := "Transcript Alternatives: 0 people love to comment on the Businessman the person who has come up soft now in stock for the next five to ten seconds in C transcribes world test of speaking and see how well matching algorithm works in can identify extra words given the initial words were correct"
	// // userTranscript := "people love to comment on the self made Businessman the person who has come up the Hard Way"
	// userTranscript := "people love to comment on the self-made businessman the person who has come up the soft way mountain people love to comment on the self-made businessman the person who has come up to soft way now I'm going to be continued to talk or think to reach about 30 seconds minute it does means it has been successfully working and yeah let's test if this works"
	errorRegion := compareText(databaseText, userTranscript)
	log.Println("errorRegion", errorRegion)

}

func handleEndAudioRecording(savePath string, conn *websocket.Conn) {

	//Convert file
	// outputFilePath := "file.wav"
	// err := convertToWav(savePath, outputFilePath, "00:00:00", "00:00:05")
	// if err != nil {
	// 	log.Fatal("Failed to convert audio file to wav:", err)
	// }

	userTranscript := GoogleSpeechToText(savePath)

	// var databaseText string = "People Love to comment on the self-made businessman. The person who has come up the hard way. This is a person who typically has reached the top without any support of qualifications. And the reason it is so comment-worthy is that it is very hard. A qualification, in whatever field you chose to study, is the foundation that underpins and support your entire career. If your career is a rocket, then the degree or diploma you read for is the launch pad. And in these modern times, you can get qualifications in almost any field. In the old days, it was just the professions like medicine, engineering or law, etc that saw you read for a degree. But now you can study for a diploma of logistics, a Bachelor of media studies or a qualification in retail management. Whatever your chosen field you can almost certainly find a course of study to help launch you on your way."
	// var userTranscript string = "People to comment on the businessman The person who has come the hard way This is a person who typically has reached the top without any support of qualifications And the reason it is so commentworthy is that it is very hard A qualification in whatever field you chose to study is the foundation that underpins and support your entire career If your career is a rocket then the degree or diploma you read for is the launch pad And in these modern times you can get qualifications in almost any field In the old days it was just the professions like medicine engineering or law etc that saw you read for a degree But now you can study for a diploma of logistics a Bachelor of media studies or a qualification in retail management Whatever your chosen field you can almost certainly find a course of study to help launch you on your way"

	var databaseText string = "People Love to comment on the self made businessman. The person who has come up the hard way. "
	// var databaseText string = "Online education, also known as e-learning or distance learning, refers to the delivery of educational content and instruction through digital platforms and the internet. It allows students to access courses, resources, and interact with instructors remotely, without the need for physical presence in a traditional classroom setting. Online education has become increasingly popular and widespread, offering numerous benefits and opportunities for learners of all ages and backgrounds. It allows flexibility, easy accessibility, variety of courses and programs, Self-paced learning, cost and time saving.	"
	// var userTranscript string = "people love to comment on the self made Businessman the person who is come up the Hard Way"
	ErrorRegions := compareText(databaseText, userTranscript)

	var tempRes speechResponse = speechResponse{
		ErrorRegions: ErrorRegions,
		Transcript:   userTranscript,
	}

	// Convert string array to JSON
	jsonData, err := json.Marshal(tempRes)
	if err != nil {
		fmt.Println("Error converting string array to JSON:", err)
		return
	}

	// Send JSON-encoded string array to the client
	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		fmt.Println("Error sending message:", err)
	}

}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	log.Println("WS Connected", r.Method)
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket connection:", err)
		return
	}
	defer conn.Close()

	//Get the browser name and version
	userAgent := r.Header.Get("User-Agent")
	pattern := `^([a-zA-Z]*)[\\]([0-9][.][0-9])`
	regx_browser, _ := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling the regular expression:", err)
		return
	}

	// Find the first match of the pattern in the input string
	match := regx_browser.FindStringSubmatch(strings.Replace(userAgent, string(filepath.Separator), "\\", -1))
	if match != nil {
		// Group 0 is the entire match, group 1 is the first captured group, and group 2 is the second captured group, and so on.
		// regx_browser_name := match[1]
		// regx_browser_version := match[2]
		// fmt.Printf("Name: %s, Version: %s\n", regx_browser_name, regx_browser_version)
	}

	savePath := "output"

	//Start reading the websocket messages
	for {
		messageType, wsMessageBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("WebSocket connection closed with error:", err)
			}
			break
		}

		if messageType == websocket.TextMessage {

			var jsonData audioChunkStruct
			err := json.Unmarshal(wsMessageBytes, &jsonData)
			if err != nil {
				log.Println("Could not unmarshal json")
			}

			var isKeyPresent bool = false
			if _, isKeyPresent = websocketConn[conn]; !isKeyPresent {
				//Adding new WebSocket Connection

				websocketConn[conn] = jsonData.UserId
				audioSegemetUser[jsonData.UserId] = int(jsonData.AudioChunk)
			}

			if isKeyPresent {
				if int(jsonData.AudioChunk) != -1 {
					audioSegemetUser[jsonData.UserId] = int(jsonData.AudioChunk)
				} else {
					log.Println("End Recording")
					handleEndAudioRecording(savePath, conn)
				}
			}

		} else if messageType == websocket.BinaryMessage {
			log.Println("websocket.BinaryMessage", websocketConn[conn])
			appendAudioChunk(wsMessageBytes, websocketConn[conn], savePath)
			userTranscript := GoogleSpeechToText(savePath)

			var tempRes speechResponse = speechResponse{
				ErrorRegions: make([][]int, 0),
				Transcript:   userTranscript,
			}

			// Convert string array to JSON
			jsonData, err := json.Marshal(tempRes)
			if err != nil {
				fmt.Println("Error converting string array to JSON:", err)
				return
			}

			// Send JSON-encoded string array to the client
			err = conn.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
		}

	}

	// Save audio buffer as WebM file
	// saveWebMFile(audioBuffer)
}

func appendAudioChunk(audioChunk []byte, userId string, savePath string) {

	log.Println("appendAudioChunk")

	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	audioBuffer.Write(audioChunk)

	err := ioutil.WriteFile(savePath, audioBuffer.Bytes(), 0644)
	if err != nil {
		log.Fatal("Failed to save audio file:", err)
	}

	// outputFilePath := "file.wav"
	// err = convertToWav(savePath, outputFilePath, "00:00:00", "00:00:05")
	// if err != nil {
	// 	log.Fatal("Failed to convert audio file to wav:", err)
	// }

	// fmt.Printf("Audio file saved: %s\n", outputFilePath)

}

func convertToWav(inputFilePath string, outputFilePath string, startTime string, duration string) error {
	// cmd := exec.Command(ffmpegPath, "-i", inputFilePath, "-c:v", "libvpx-vp9", "-c:a", "libopus", "-b:v", "1M", "-b:a", "64k", "-deadline", "realtime", "-cpu-used", "4", outputFilePath)

	// cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-ss", startTime, "-t", duration, "-ac", "2", "-f", "wav", outputFilePath, "-y")
	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-ac", "2", "-f", "wav", outputFilePath, "-y")
	log.Println("cmd", cmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func pythonTranscribe() string {

	// Specify the virtual environment directory path
	virtualEnvDir := "../lq_go/myenv"

	// Specify the Python script path
	pythonScript := "app.py"

	// Get the path to the Python interpreter in the virtual environment
	pythonInterpreter := filepath.Join(virtualEnvDir, "bin", "python")

	// Create the command to execute the Python script using the virtual environment
	cmd := exec.Command(pythonInterpreter, pythonScript)

	// Set the environment variables for the command (if needed)
	cmd.Env = os.Environ()

	// Create byte buffers to capture the script output
	var stdout bytes.Buffer
	// var stderr bytes.Buffer
	cmd.Stdout = &stdout
	// cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		log.Fatal("Failed to execute Python script:", err)
	}

	// Retrieve the script output
	output := stdout.String()
	// str_errors := stderr.String()

	log.Println("output", output)

	return output
}

func generateRandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
