package speed

import (
//	"fmt"
	"strings"
	"testing"
	"xml"

	gxml "g/xml"
)

const testXML = "<message from='user1@server1' to='user2@server2' type='chat'>" +
	"<subject>Let's talk about hapiness</subject>" +
	"<body>How do you feel about this?</body>" +
	"<thread>asfasdfadfasdfasdfasfasdfsadf</thread>" +
	"</message>"

func Test(t *testing.T) {
}

func BenchmarkGXML(bm *testing.B) {
	for i := 0; i < bm.N; i++ {
		r := gxml.NewReader(strings.NewReader(testXML))
		f := r.ReadElement()
		if f == nil { panic("") }
		// fmt.Println(f)
		// return
	}
}

type Message struct {
//	XMLName	xml.Name	"message"
	From	string	"attr"
	To	string	"attr"
	Type	string	"attr"
	Subject	string
	Thread	string
	Body	string
}

func BenchmarkXML(bm *testing.B) {
	for i := 0; i < bm.N; i++ {
		var m Message
		e := xml.Unmarshal(strings.NewReader(testXML), &m)
		if e != nil { panic(e) }
		// fmt.Println(m)
		// return
	}
}
