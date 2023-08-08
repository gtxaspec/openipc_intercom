package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var config map[string]string

type WAVHeader struct {
	RIFFHeader    [4]byte
	RIFFSize      uint32
	WAVEHeader    [4]byte
	FMTHeader     [4]byte
	FMTSize       uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	DataHeader    [4]byte
	DataSize      uint32
}

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *debug {
		log.Println("Debug mode activated")
	}

	configBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/", http.FileServer(http.Dir(config["webPath"])))
	log.Printf("Server started on http://localhost:%s", config["port"])
	log.Fatal(http.ListenAndServe(":"+config["port"], nil))
}

func downsample(data []int16, oldRate, newRate int) []int16 {
	ratio := float64(oldRate) / float64(newRate)
	newSize := int(float64(len(data)) / ratio)
	downsampled := make([]int16, newSize)

	for i := 0; i < newSize; i++ {
		srcIdx := float64(i) * ratio
		leftIdx := int(srcIdx)
		rightIdx := leftIdx + 1

		if rightIdx >= len(data) {
			rightIdx = len(data) - 1
		}

		alpha := srcIdx - float64(leftIdx)
		downsampled[i] = int16((1-alpha)*float64(data[leftIdx]) + alpha*float64(data[rightIdx]))
	}

	return downsampled
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Received a request to upload handler")

	file, _, err := r.FormFile("audio")
	if err != nil {
		log.Printf("Error reading the file from the request: %v", err)
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Error reading the file bytes: %v", err)
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	if len(fileBytes) == 0 {
		http.Error(w, "Empty file", http.StatusBadRequest)
		return
	}

	// Parse WAV header
	var header WAVHeader
	 err = binary.Read(bytes.NewReader(fileBytes), binary.LittleEndian, &header)
	if err != nil {
		log.Printf("Error parsing WAV header: %v", err)
		http.Error(w, "Invalid WAV file", http.StatusBadRequest)
		return
	}

	// Validate WAV header
	if string(header.RIFFHeader[:]) != "RIFF" || string(header.WAVEHeader[:]) != "WAVE" || string(header.DataHeader[:]) != "data" {
		http.Error(w, "Invalid WAV file, header error.", http.StatusBadRequest)
		return
	}

	// Extract PCM data from WAV
	pcmDataStart := binary.Size(header)
	pcmData := make([]int16, header.DataSize/2)
	err = binary.Read(bytes.NewReader(fileBytes[pcmDataStart:]), binary.LittleEndian, &pcmData)
	if err != nil {
		log.Printf("Error extracting PCM data: %v", err)
		http.Error(w, "Invalid WAV file, PCM extraction error", http.StatusBadRequest)
		return
	}

	// Downsample the PCM data
	downsampledData := downsample(pcmData, int(header.SampleRate), 16000)

	// Pad with silence
	paddedData := make([]int16, 8000+len(downsampledData)+24000)
	copy(paddedData[8000:], downsampledData)

	// Convert padded data to bytes
	paddedBytes := make([]byte, len(paddedData)*2)
	for i, sample := range paddedData {
		binary.LittleEndian.PutUint16(paddedBytes[i*2:i*2+2], uint16(sample))
	}

	// Save the padded PCM data (without headers)
	outputPath := config["uploadPath"] + "padded.pcm"
	err = ioutil.WriteFile(outputPath, paddedBytes, 0644)
	if err != nil {
		log.Printf("Failed to save padded PCM data: %v", err)
		http.Error(w, "Failed to save padded PCM data", http.StatusInternalServerError)
		return
	}

	log.Printf("Size of padded PCM data written: %d bytes", len(paddedBytes))

	// Send the padded PCM data to the playAudioURL endpoint
	log.Printf("Attempting to send PCM data to: %s", config["playAudioURL"])

	req, err := http.NewRequest("POST", config["playAudioURL"], bytes.NewReader(paddedBytes))
	if err != nil {
		log.Printf("Failed to create POST request: %v", err)
		http.Error(w, "Failed to send audio data", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to send POST request: %v", err)
		http.Error(w, "Failed to send audio data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("Received response status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-200 response code: %d", resp.StatusCode)
		http.Error(w, "Failed to send audio data", http.StatusInternalServerError)
		return
	}

	log.Printf("Audio data sent successfully to %s", config["playAudioURL"])

	w.Write([]byte("Audio uploaded, downsampled, padded, and sent successfully."))
}
