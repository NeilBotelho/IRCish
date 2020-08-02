main: main.go
	GOOS=linux GOARCH=amd64 go build 
	echo "web: ./server \$$PORT" >Procfile
	touch requirements.txt
