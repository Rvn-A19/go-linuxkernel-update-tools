
all: getkver

getkver:
	go build .

clean:
	rm ./getkver
