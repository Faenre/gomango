# Go-Man-Go

A small server that opens a simple web listener to a mongodb 

## Installation

- Build the server
- Config the config.yml file (right now, only the fully-qualified MongoDB URI is supported)
- Make sure all expected tracer sources are included in the yml file

## Usage

- Send traffic (GET, POST, ...) to the server, using the port specified in the config file
- The message headers are copied, and the content is dumped as a raw string into the db
