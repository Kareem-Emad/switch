install:
	GO114MODULE=on go mod tidy
build: install
	GO114MODULE=on go build -o switch.bin .

run: build
	./switch.bin

# Example command
# make run
# for faktory server
# FAKTORY_PASSWORD=as faktory -b 127.0.0.1:7421 -w 127.0.0.1:7422