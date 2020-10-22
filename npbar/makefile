
UNAME := $(shell uname)

build: $(shell find . -regex ".*\.go" -type f)
ifeq ($(UNAME), Linux)
	sudo apt-get install scons

	cd /tmp; \
	rm -rf rpi_ws281x; \
	git clone https://github.com/jgarff/rpi_ws281x --quiet; \
	cd rpi_ws281x; \
	scons; \
	find . -name '*.a' | sudo cpio -pdm /usr/local/bin; \
	find . -name '*.h' | sudo cpio -pdm /usr/local/include

	
	go build -v cmd/run/main.go 
else 
ifeq ($(UNAME), Darwin)
	go build -v cmd/runtest/main.go 
else 
	$(error Must build on a mac or raspberry pi)
endif
endif

start:
	./main $(DOMINION_IP) $(DOMINION_PORT) $(DOMAIN_UUID) $(SERVICE_PORT)