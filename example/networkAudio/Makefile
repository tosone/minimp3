Target = $(shell basename $(abspath $(dir $$PWD))).test

all: clean build

build:
	CGO_ENABLED=1 go build -v -o $(Target) -ldflags "-s -w"

clean:
	$(RM) $(Target)

test:
	./$(Target) http://119.90.25.29/m7c.music.126.net/20180712165443/9870883d0594edc455903140c5b325c8/ymusic/aac8/e215/54a8/6c49025c90e6bd211bf61bd91efed99d.mp3
