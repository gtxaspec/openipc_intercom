send audio to an OpenIPC /play_audio endpoint from a web browser or Home Assistant card.

Runs ON-DEVICE (mipsle) or on another server

Tested with Ingenic T31 devices.

compile with:

`./compile.sh`

install UPX for smaller binary sizes

run:

rename `config.json.example` to `config.json`, make your changes, then run

`./intercom`

for debug output:
`./intercom --debug`

based off `https://github.com/addpipe/simple-recorderjs-demo/tree/master`
