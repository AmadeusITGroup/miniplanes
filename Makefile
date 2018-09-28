
all:

OUTPUTDIR=output

output:
	mkdir -p $(OUTPUTDIR)

validate-itineraries-server:
	cd itineraries-server/swagger && swagger validate ./swagger.yaml

build-itineraries-server: output
	go build  -o $(OUTPUTDIR)/itineraries-server itineraries-server/cmd/itineraries-server/main.go

gen-itineraries-server:
	cd itineraries-server/swagger && swagger generate server --target ../pkg --name itineraries --spec ./swagger.yaml

clean: $(OUTPUTDIR)
	rm -rf $(OUTPUTDIR)
