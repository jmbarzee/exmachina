cI use this code on raspberry pis. I've tried to document the setup here, both for myself and for you.  


## SSH over USB
[(source)](https://www.thepolyglotdeveloper.com/2016/06/connect-raspberry-pi-zero-usb-cable-ssh/)

## Building a piZero
I've left these instructions for myself, I struggled to follow the complete headless start due to problems with network interfaces and wpa_supplicant not playing nicely. When in doubt, read more about the commands that you are copy pasting from some fourm with 2 upvotes on the post.

1. Burn [Raspbain Stretch Lite](https://downloads.raspberrypi.org/raspbian_lite_latest) to a microSD with [Etcher](https://etcher.io/)
2. Disconect mircoSD (safely) and then reconnect to view contents
3. Enable UART by adding a line to `config.txt`    [(source)](https://learn.adafruit.com/raspberry-pi-zero-creation/text-file-editing)
    ```# Enable UART
    enable_uart=1```
4. Disconnect mircoSD (safely) and place in piZero
5. Plug console cable into computer (not piZero)
5. Install drivers for console cable on your computer [(source)](https://learn.adafruit.com/adafruits-raspberry-pi-lesson-5-using-a-console-cable/software-installation-mac)
6. From the terminal `ls /dev/cu.*` [(source)](https://learn.adafruit.com/adafruits-raspberry-pi-lesson-5-using-a-console-cable/test-and-configure)
7. Look for `cu.SLAB_USBtoUART` or something similar
8. `screen /dev/cu.SLAB_USBtoUART 115200` it takes a while to get used to 
9. Connect console cable [(source1)](https://learn.adafruit.com/adafruits-raspberry-pi-lesson-5-using-a-console-cable/test-and-configure) [(source2)](https://learn.adafruit.com/adafruits-raspberry-pi-lesson-4-gpio-setup/the-gpio-connector)
    ```Red       Disconnected
    Black    GND
    White    TXD
    Green    RXD
10. Power up the piZero. Green light should turn on and then flicker intermentantly. Screen should be reading output from serial pins of piZero. Wait for output to stop and bash line appear.
11. search results of `sudo iwlist wlan0 scan` for the name of your network. This demonstrates that the piZero can see networks [(source)](https://www.raspberrypi.org/documentation/configuration/wireless/wireless-cli.md)
12. `sudo vi /etc/wpa_supplicant/wpa_supplicant.conf` and append
    ```network={
        ssid="networkName"
        psk="networkPassword"
    }
13. `wpa_cli -i wlan0 reconfigure` (try running this twice, sometimes the command times out on the first call)
    1. If your output looks bad (doesn't just read `OK`) then something is screwed.
    2. first, check `/etc/wpa_supplicant/wpa_supplicant.conf` it should only have 6 or 7 lines
        ```ctrl_interface=DIR=/var/run/qpa_supplicant GROUP=netdev
        update_config=1

        network={
            ssid="networkName"
            psk="networkPassword"
        }
    3. try `ip link set wlan0 up` (down will do the oposite) [(source)](https://wiki.archlinux.org/index.php/Network_configuration#Network_interfaces)
    4. try `wpa_supplicant -B -i wlan0 -c /etc/wpa_supplicant/wpa_supplicant.conf`
    5. If `wpa_supplicant` started successfully, you should be able to redo #13 on this list without any bugs
14. `ifconfig wlan0` the piZero should have an address assigned to it and be connected to your network
15. Enable ssh by `sudo raspi-config` >> `Interfacing Options` >> `SSH` >> `Yes` >> `OK` >> `Finish`
        or
    ```sudo systemctl enable ssh
    sudo systemctl start ssh
16. `ssh pi@raspberrypi.local` or `ssh pi@PI.IP.AD.DR` (password is `raspberry` by default)
17. `passwd` and change password to prevent your pi becoming a zombie


## No Password Login
1. `ssh-keygen -t rsa` ENTER to every field [(source)](https://stackoverflow.com/questions/12202587/automatically-enter-ssh-password-with-script)
2. `ssh-copy-id pi@raspberrypi.local` or `ssh-copy-id pi@PI.IP.AD.DR`
3. `ssh pi@raspberrypi.local` or `ssh pi@PI.IP.AD.DR` (password is `raspberry` by default)

## Cross compiling Golang on the Pi
[(source1)](https://www.thepolyglotdeveloper.com/2017/04/cross-compiling-golang-applications-raspberry-pi/)
[(source2)](https://stackoverflow.com/questions/32309030/go-1-5-cross-compile-using-cgo-on-os-x-to-linux-and-windows)
[(source3)](https://hub.docker.com/r/alexellis2/go-armhf/tags/)

# Docker 

## Docker on rpi
1. `curl -sSL https://get.docker.com | sh` [(source)](https://www.raspberrypi.org/blog/docker-comes-to-raspberry-pi/)

## Cross Archetecture Docker
[(source)](https://docs.docker.com/docker-for-mac/multi-arch/)