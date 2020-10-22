
#How to set up a raspberry Pi Zero as a npBar

*commands are added as that part of the process is reviewed*



Install Golang
`export GOLANG="$(curl https://golang.org/dl/|grep armv6l|grep -v beta|head -1|awk -F\> {'print $3'}|awk -F\< {'print $1'})"`
`wget https://golang.org/dl/$GOLANG`
`sudo tar -C /usr/local -xzf $GOLANG`
`rm $GOLANG`
`unset GOLANG`

Download Domain & Services
`cd /usr/local/`
`sudo git clone https://github.com/jmbarzee/dominion`
`sudo chown -R $(whoami) /usr/local/dominion`
`cd dominion`
`rm -rf services`
`sudo git clone https://github.com/jmbarzee/services`
`sudo /usr/local/go/bin/go build cmd/domain/main.go`