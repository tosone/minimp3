TARGET = $(shell basename $(abspath $(dir $$PWD))).test

all: clean build

build:
	env GO111MODULE=off CGO_ENABLED=1 go build -v -o $(TARGET) -ldflags "-s -w"

clean:
	$(RM) $(TARGET)

test:
	./$(TARGET)
