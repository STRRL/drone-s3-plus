FROM plugins/base:multiarch

LABEL org.label-schema.name="Drone S3 Plus"

ADD release/linux/amd64/drone-s3-plus /bin/
ENTRYPOINT ["/bin/drone-s3-plus"]
