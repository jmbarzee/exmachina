build: $(shell find . -regex ".*\.go" -type f)
	# go build cmd/run/main.go 

start:
	go run cmd/run/main.go  $(DOMINION_IP) $(DOMINION_PORT) $(DOMAIN_UUID) $(SERVICE_PORT)