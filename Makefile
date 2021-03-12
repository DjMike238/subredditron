PROGNAME = subredditron

$(PROGNAME):
	CGO_ENABLED=0 go build -o $(PROGNAME)

all: $(PROGNAME)

.PHONY: all $(PROGNAME) clean

tiny:
	CGO_ENABLED=0 go build -o $(PROGNAME) -ldflags="-s -w"
	upx --brute $(PROGNAME)

clean:
	rm $(PROGNAME)
