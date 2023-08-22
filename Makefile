
APP_NAME=timeandtemp

run: binary
	@temp/$(APP_NAME)

binary:
	@mkdir -p temp
	@CGO_ENABLED=0 go build -o temp/$(APP_NAME) main.go

deploy: binary
	scp temp/$(APP_NAME) optiplex:/tmp
	ssh optiplex sudo /tmp/$(APP_NAME) -action deploy




