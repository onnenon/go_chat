# build image
FROM golang as builder

WORKDIR /opt

COPY . .

# compile
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o client cmd/client/* 

# deployment image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

LABEL author="Stephen Onnen"

WORKDIR /root/

COPY --from=builder /opt/ .

CMD ["./client", "-server=server:9000"]

EXPOSE 9000 