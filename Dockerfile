FROM alpine:latest
RUN apk --no-cache add ca-certificates
EXPOSE 8080
WORKDIR /

ADD _output/storage .
ADD _output/itineraries-server .
ADD _output/ui .
ADD _output/schedules-generator .
