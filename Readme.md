

# About Pipcas gateway

Pipcas Gateway sits in front of a Pipcas node containing multiple shards. Gateway handles routing of read and write requests among the shards. Optimized connection management with high throughput.  


# Quickstart

## Deploy the binaries on your local machine directly onto your host os.

make sure Go 1.18 is installed and it's path variable is configured

clone the repository
```
    git clone https://github.com/raja-dettex/pipcas-lb
```
```
    cd pipcas-lb
```

```
    go mod tidy
```

```
    make build
```

```
    make run
```

## Deploy onto a docker cotainer

after cloning the repository as mentioined previously,

Build a Docker image from the Docker file

```
    docker build -t <your image name> .
```

Start the docker container from the built image

```
    docker run -d -p <port to access contiainer process>:<port according to listen address> -e  LISTEN_ADDR=<port> <your image name>
```

## latest release: 
    pipcas-lb:1.0
