FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache gcc musl-dev make bash

CMD ["tail","-f","/dev/null"]
