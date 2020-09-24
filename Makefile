install:
	GO114MODULE=on go mod tidy

build: install
	GO114MODULE=on go build -o ingestion-engine.bin .

run: build
	./ingestion-engine.bin \
	-execution-group-id=$(execution-group-id) \
	-video-path=$(video-path) \

# Example command
# orch: make execution-group-id=exec_group_1 video-path=/home/sayed/series/test/file.txt model-path=/home/sayed/series/test/file.txt model-config-path=/home/sayed/series/test/file.txt start-idx=5 frame-count=50 -code-path ./path/to/code_driectory -video-token=mm.mp4 run
# FAKTORY_PASSWORD=as faktory -b 127.0.0.1:7421 -w 127.0.0.1:7422