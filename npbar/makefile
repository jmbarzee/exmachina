
UNAME := $(shell uname)

build:
ifneq ($(UNAME), Linux)
	$(error Must build on a raspberry pi)
else 
	sudo apt-get install scons

	cd /tmp; \
	rm -rf rpi_ws281x; \
	git clone https://github.com/jgarff/rpi_ws281x --quiet; \
	cd rpi_ws281x; \
	scons; \
	find . -name '*.a' | sudo cpio -pdm /usr/local/bin; \
	find . -name '*.h' | sudo cpio -pdm /usr/local/include

	
	go build -v cmd/run/main.go 
endif

start:
	./main $(DOMINION_IP) $(DOMINION_PORT) $(DOMAIN_UUID) $(SERVICE_PORT)