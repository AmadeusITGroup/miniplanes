FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /
ADD itineraries-server .

#ENTRYPOINT ["/itineraries-server"]
#CMD ["/itineraries-server", "--port=4838"]