FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/expense-tracker main.go

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=builder /app/bin/expense-tracker ./expense-tracker

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/expense-tracker"]
