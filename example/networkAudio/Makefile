Target = $(shell basename $(abspath $(dir $$PWD))).test

all: clean build

build:
	CGO_ENABLED=1 go build -v -o ${Target} -ldflags "-s -w"

clean:
	-rm ${Target}

test:
	./${Target} http://119.90.25.30/ip.h5.ra01.sycdn.kuwo.cn/e3244f78b1b42da5a4e74075b5870cce/5ab8b559/resource/n1/320/62/90/2295644473.mp3
