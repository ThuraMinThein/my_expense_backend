From golang:1.24

WORKDIR /app

COPY go.* ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]