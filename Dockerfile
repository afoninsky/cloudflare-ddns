FROM golang:1.18.3 AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /service .
RUN useradd -u 10001 scratchuser


FROM scratch
COPY --from=builder /service /service
COPY --from=builder /etc/passwd /etc/passwd
USER 10001

ENTRYPOINT [ "/service" ]
