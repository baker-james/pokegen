FROM golang:1.19
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /pokegen
EXPOSE 8080
CMD ["/pokegen"]
