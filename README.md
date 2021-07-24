# Anthole
Anthole is a simple intranet penetration tool implemented in Golang. It uses the Cobra scaffolding library to build the command line, which is very simple to use.

## Building the source
First, you need a Golang compilation environment. Please find the relevant documents on how to install Golang, then clone this repository and switch to the `cmd\anthole` folder locally to compile

## Usage

### Server
Run anthole server -c Configuration file Path

### Client
Run anthole client -c Configuration file Path

## Configuration file

Create a yaml file locally, add the following configuration items, and modify the host and port in the configuration

```
common:
  token: xiaolaji
server:
  host: 139.10.33.85
  port: 15555
client:
  services:
  - local_host: 127.0.0.1
    local_port: 3306
    remote_port: 16666
    type: tcp
  - local_host: 127.0.0.1
    local_port: 6279
    remote_port: 17777
    type: tcp
```

## Running

<p><img src="https://github.com/danzhuxia/Anthole/blob/main/images/server.png" alt="cmd-markdown-logo" title="Server Running" /></p>
<p><img src="https://github.com/danzhuxia/Anthole/blob/main/images/client.png" alt="cmd-markdown-logo" title="Client Running" /></p>