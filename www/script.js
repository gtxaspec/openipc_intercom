let audioContext;
let recorder;
let recordButton = document.getElementById("recordButton");
let stream;  // Add this line to keep a reference to the media stream

recordButton.addEventListener("mousedown", startRecording);
recordButton.addEventListener("mouseup", stopRecording);

function startRecording() {
    console.log("Attempting to start recording...");
    
    // Initialize the AudioContext inside the function
    if (!audioContext) {
        audioContext = new (window.AudioContext || window.webkitAudioContext)();
    }

    navigator.mediaDevices.getUserMedia({ audio: true })
        .then(mediaStream => {
            stream = mediaStream;  // Store the media stream
            let source = audioContext.createMediaStreamSource(mediaStream);
            recorder = new Recorder(source, { numChannels: 1 });
            recorder.record();  // Start recording
            recordButton.innerText = ". . . Recording . . .";
            console.log("Recording started.");
        })
        .catch(err => {
            console.error("Error accessing the microphone:", err);
        });
}

function stopRecording() {
    console.log("Attempting to stop recording...");
    recorder.stop();
    recordButton.innerText = "Push and Hold to Record";

    // Stop the MediaStream tracks (this releases the microphone)
    stream.getTracks().forEach(track => track.stop());

    // Export the WAV audio
    recorder.exportWAV(sendDataToServer);

    // Clear the recorder and release resources
    recorder.clear();
    console.log("Recording stopped. Sending data...");
}

function sendDataToServer(blob) {
    console.log("Sending WAV data to server...");
    const formData = new FormData();
    formData.append('audio', blob, 'audio.wav');  // Send as WAV

    fetch("/upload", {
        method: "POST",
        body: formData
    })
    .then(response => response.text())
    .then(data => {
        console.log("Server Response:", data);
    })
    .catch(error => {
        console.error("There was an error uploading the audio:", error);
    });
}

