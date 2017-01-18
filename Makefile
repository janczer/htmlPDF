all:
	rm -f htmlPDF
	go build
	./htmlPDF

clean:
	rm -f htmlPDF

