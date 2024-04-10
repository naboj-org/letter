FROM golang:1.21-bookworm as build
WORKDIR /build
COPY . .
RUN go build -o letter

FROM dxjoke/tectonic-docker:latest
WORKDIR /app
RUN adduser texuser
USER texuser

COPY --from=build /build/letter /app/letter

CMD ["/app/letter"]
