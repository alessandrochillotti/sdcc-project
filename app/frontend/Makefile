frontend: 
	go build -o ./frontend ./frontend.go ./lib_frontend.go

test:
	go build -o ./test ./test.go ./lib_frontend.go

clean:
	rm -f ./frontend ./test

all: frontend test

