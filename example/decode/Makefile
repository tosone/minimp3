Target = $(shell basename $(abspath $(dir $$PWD))).test

all: clean build

build:
	GO111MODULE=off CGO_ENABLED=1 go build -v -o $(Target) -ldflags "-s -w"

clean:
	$(RM) $(Target)

test:
	./$(Target)
