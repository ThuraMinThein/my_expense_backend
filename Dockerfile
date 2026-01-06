From golang:1.25

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.* ./
RUN go mod download

COPY . .

# dev
CMD ["air", "-c", ".air.toml"]

# prod
# RUN go build -o main ./cmd/

# CMD ["./main"]