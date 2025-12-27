FROM golang:1.25 AS deps
WORKDIR /hello-api
COPY *.mod *.sum ./
RUN go mod download

FROM deps AS dev
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -X main.docker=true" -o api cmd/main.go
CMD ["/hello-api/api"]

FROM scratch AS prod
WORKDIR /
EXPOSE 8080
COPY --from=dev /hello-api/api /
CMD ["/api"]

