# Heroku doesn't directly support golang all that well.
# This is a workaround, where I can use travis to build the binary and
# Trick heroku into thinking this is a python project and just run the binary
main: main.go
	GOOS=linux GOARCH=amd64 go build 
	echo "web: ./server \$$PORT" >Procfile
	touch requirements.txt
