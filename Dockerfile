FROM golang:1.19 as builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /pokegen

FROM scratch
COPY --from=builder /pokegen /pokegen
EXPOSE 8080
CMD ["/pokegen"]
