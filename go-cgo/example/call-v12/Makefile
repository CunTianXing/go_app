default:
	go tool cgo main.go hello_go.go

run:
	@go build -o a.out
	@./a.out
	-@rm a.out

debug:
	nm ./_obj/_cgo_.o

clean:
	-rm -rf _obj
