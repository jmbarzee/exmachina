
UNAME := $(shell uname)

ifneq ($(UNAME), Linux)
build:
	$(error Must build on a raspberry pi)
else 
build:

	apt-get install scons

	cd /tmp; \
	rm -rf rpi_ws281x; \
	git clone https://github.com/jgarff/rpi_ws281x --quiet; \
	cd rpi_ws281x; \
	scons; \
	find . -name '*.a' | cpio -pdm /usr/local/bin; \
	find . -name '*.h' | cpio -pdm /usr/local/include

	
	go build cmd/run/main.go 
endif

get_scons_linux:

start:
	./main $(DOMINION_IP) $(DOMINION_PORT) $(DOMAIN_UUID) $(SERVICE_PORT)