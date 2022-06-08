FROM golang:1.18-bullseye as build
WORKDIR /build
COPY . .
RUN go build -o letter

FROM texlive/texlive:latest
WORKDIR /app
RUN adduser texuser
USER texuser

COPY --from=build /build/letter /app/letter

CMD ["/app/letter"]
