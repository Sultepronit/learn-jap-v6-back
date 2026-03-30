FROM alpine:3.19

WORKDIR /app

# COPY . .

EXPOSE 8080

CMD ["./app"]