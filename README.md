# OpenIPC Audio Endpoint Interface

This repository provides tools to send audio to an OpenIPC `/play_audio` endpoint directly from a web browser or a Home Assistant card.

## Compatibility

- Designed primarily for Ingenic T31 devices.
- Can be deployed directly on-device or on an external server.
- Written in Golang, cross-compile for most processor architectures.
- Binaries for amd64 & mipsle included.

## Setup & Installation

### Compilation

Compile the source code with the provided script:

`./compile.sh`

**Note**: For optimal results and reduced binary sizes, ensure you have UPX (UPX 4.0.2 or later for mipsle) installed, then run:

`./compile.sh upx`

### Configuration

1. Rename the sample configuration file:

`mv config.json.example config.json`

2. Update `config.json` with your specific settings.

### Running the Server

To start the server, use:

`./intercom`

For debug output, use the `--debug` flag:

`./intercom --debug`

### Accessing the Interface

Tested on Chrome.  Add the server url to `chrome://flags/#unsafely-treat-insecure-origin-as-secure` in order to support microphone capability (since the server is not https).

Once the server is running, access the interface from your browser:

`http://<ip-address>:3333/index.html`

## Credits

This project is inspired by and based on the [simple-recorderjs-demo](https://github.com/addpipe/simple-recorderjs-demo/tree/master) from addpipe.
