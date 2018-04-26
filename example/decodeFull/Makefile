Target = $(shell basename $(abspath $(dir $$PWD))).test

all: clean build

build:
	CGO_ENABLED=1 go build -v -o ${Target} -ldflags "-s -w"

clean:
	-rm ${Target}

test:
	./${Target}
