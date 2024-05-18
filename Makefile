host=ubuntu-1
#home=marchome
# cannot for the life of me figure out how to change dns on ubuntu
home=192.168.50.89

build:
	go build -o main main.go server.go

run:
	go run main.go server.go

on:
	curl -vvv http://$(host):1323/controls/dimmer/255

off:
	curl -vvv http://$(host):1323/controls/dimmer/0

ssh:
	ssh -A $(host) -l marc -t "cd play-go; setup; zsh --login"

get:
	rsync -avz --exclude .git --exclude node_modules --exclude main -e ssh marchome@$(home):develop/mholzen/play-go/ .

pull: get

push:
	rsync -avz --exclude .git --exclude node_modules --exclude main -e ssh . marc@$(host):play-go

status:
	sudo systemctl status play-go.service

stop:
	sudo systemctl stop play-go-watcher.service
	sudo systemctl stop play-go.service

start:
	sudo systemctl start play-go.service
	sudo systemctl start play-go-watcher.service

log:
	journalctl -u play-go.service -f

live:
	go run live.go