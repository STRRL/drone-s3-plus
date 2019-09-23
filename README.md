# drone-s3-plus

fork from https://github.com/drone-plugins/drone-s3

Drone plugin to publish files and artifacts to Amazon S3 or Minio

## Changes
- add overwrite option(default true)
- upload concurrently, it will save you a lot time if you have a lot artifacts
- better output, artifacts size and time for upload will be print
- support dry-run, debug only

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

or in your .drone.yml
```
  - name: release
    image: f1shl3gs/drone-s3-plus
    settings:
      bucket: <bucket>
      path_style: true
      access_key: <access_key>
      secret_key: <secret_key>
      source: release/*
      strip_prefix: release/
      endpoint:   http://<your_damoin>
```
