# drone-s3-plus

fork from https://github.com/drone-plugins/drone-s3

## Changes
- add overwrite option(default true)

Drone plugin to publish files and artifacts to Amazon S3 or Minio. For the
usage information and a listing of the available options please take a look at
[the docs](http://plugins.drone.io/drone-plugins/drone-s3/).

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the Docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm=true -t plugins/s3 .
```

Please note incorrectly building the image for the correct x64 linux and with
CGO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-s3' not found or does not exist..
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e PLUGIN_SOURCE=<source> \
  -e PLUGIN_TARGET=<target> \
  -e PLUGIN_BUCKET=<bucket> \
  -e AWS_ACCESS_KEY_ID=<token> \
  -e AWS_SECRET_ACCESS_KEY=<secret> \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/s3 --dry-run
```

```
  AWS_ACCESS_KEY_ID=K3V27FPQO9FCN58IK4FY AWS_SECRET_ACCESS_KEY=bBNVauNe9EAIXH89NkBmJW2Z82huT9DvuDw85kir
```