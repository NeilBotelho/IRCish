main: main.go
	GOOS=linux GOARCH=amd64 go build 
	mv server bin/main
	echo "web: ./bin/main \$$PORT" >Procfile
	touch requirements.txt
