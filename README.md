send audio to an OpenIPC /play_audio endpoint from a web browser or Home Assistant card.

Runs ON-DEVICE (mipsle) or on another server

Tested with Ingenic T31 devices.

compile with:

`GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o intercom_mipsle main.go`

based off `https://github.com/addpipe/simple-recorderjs-demo/tree/master`
