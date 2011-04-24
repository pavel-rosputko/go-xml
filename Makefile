include $(GOROOT)/src/Make.inc

TARG = github.com/pavelrosputko/go-xml

GOFILES= \
	parser.go \
	fragment.go \
	reader.go \
	iterator.go \
	cursor.go \
	writer.go \
	builder.go \
	string-stack.go \
	desc.go
	
include $(GOROOT)/src/Make.pkg

