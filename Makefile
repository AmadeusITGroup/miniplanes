
all:

OUTPUTDIR=output

output:
	mkdir -p $(OUTPUTDIR)

build-itineraries-server: output
	go build  -o $(OUTPUTDIR)/itineraries-server itineraries-server/cmd/itineraries-server/main.go

gen-itineraries-server:
	cd itineraries-server/swagger && swagger generate server --target .. --name itineraries --spec ./swagger.yaml


clean: $(OUTPUTDIR)
	rm -rf $(OUTPUTDIR)
