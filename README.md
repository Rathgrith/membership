# ECE428_MP2

## Project Layout

`_test/`: Supports commands via UDP requests, including `list_mem`, `list_id`, `leave`,`enable/disable_suspicion`

`cmd/`:  Applications for this MP, contains main start function `gossip.go`

`config/`: Configuration file(s) and their load functions. It also implements validate functions and some Getters for the user.

`pkg/`: Library code that can be used in application code (and also be used inside the pkg package), including request handler for the client, client-server interface definition, network tools, and logger.

`scripts/`: Scripts to perform build, deploy, performance measure, etc operations.

## Overview

This repository currently contains MP2 code for this class. This repository consists of one major application:

\- `cmd/gossip.go`: It can dispatch `grep` commands (with full support of argument options) to all servers and return aggregated results to the user.

## Getting started

To run the program, please make sure that current directory in located at /ece428_mp1 such that it can correctly loads the config yaml.

Then, call the following commands to activate corresponding processes.

```go
go run ./cmd/gossip.go
```

You can also build your .exe binaries to run the service, where we have make one automated server builder `./scripts/server_build.sh` that builds the most updated version of codes on remote git repository automatically. (recommended). You should already have an SSH key on your VM that can access your Gitlab repository.



## Deployment

- Simply use ssh to clone your repository using ssh.

  ```shell
  ssh dl58@sftp dl58fa23-cs425-48XX.cs.illinois.edu
  ssh> cd ./go
  ssh> git clone https://gitlab.engr.illinois.edu/dl58/ece428_mp2.git
  ```
  
- We have made an automatic deployment script `./scripts/deploy.sh` to automate this troublesome process. And you should run this script under `scripts` folder. 

