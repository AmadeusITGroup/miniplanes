
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

gen-itineraries-client:
	cd itineraries-server/swagger && swagger generate client --target ../pkg --name itineraries --spec ./swagger.yaml

build-ui: $(OUTPUTDIR)
	cd ui && go-bindata -o=assets/bindata.go --nocompress --nometadata --pkg=assets templates/... static/...
	go build -o $(OUTPUTDIR)/ui ui/cmd/main.go

clean: $(OUTPUTDIR)
	rm -rf $(OUTPUTDIR)


#TMP target for local tests
start-ui: ${OUTPUTDIR}
	cd ui && ./ui

start-itineraries-server: $(OUTPUTDIR)
	output/itineraries-server --port=41807
