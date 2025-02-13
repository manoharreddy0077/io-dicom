FROM alpine:latest

RUN mkdir /app

WORKDIR /app

ADD build/linux/amd64/io-dicom /app

ENTRYPOINT [ "/app/io-dicom", "-scp", "-calledae", "DICOM_SCP", "-port", "1040", "-datastore", "/datastore" ]
