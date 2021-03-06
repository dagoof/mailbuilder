package mailbuilder

import (
	"bytes"
	"math/rand"
	"strconv"
)

//Represents a MIME multipart.
type MultiPart struct {
	content               []BodyPart
	headers               map[string]string
	boundary, contentType string
}

//Takes a content type and returns a new *MultiPart with a boundary initialized to a random number.
//No seed, but we only care about consistency within the message so it's OK
func NewMultiPart(ctype string) *MultiPart {
	bound := "--==_Part_" + strconv.Itoa(rand.Int()) + "=="
	return &MultiPart{make([]BodyPart, 0),
		make(map[string]string, 0),
		bound,
		ctype}
}

//Adds a part to the multipart
//Use this to attach a SimplePart or nest MultiParts
func (self *MultiPart) AddPart(part BodyPart) {
	self.content = append(self.content, part)
}

//Add the ContentType header.
//We treat this differently because of slightly different formatting requirements
func (self *MultiPart) SetContentType(ctype string) {
	self.contentType = ctype
}

//Add a header to the multipart. Don't use for content type, use SetContentType
func (self *MultiPart) AddHeader(key, value string) {
	self.headers[key] = value
}

//Builds the body of the multipart.
func (self *MultiPart) Bytes() []byte {
	var b bytes.Buffer
	b.WriteString("Content-Type: " + self.contentType + ";\n\tboundary=\"" + self.boundary + "\"")
	for k, v := range self.headers {
		b.WriteString(k + ": " + v + "\n")
	}
	b.WriteString("\n\nThis is a multi-part message in MIME format.\n\n")
	for _, f := range self.content {
		b.WriteString("--" + self.boundary + "\n")
		b.Write(f.Bytes())
		b.WriteString("\n")
	}
	b.WriteString("--" + self.boundary + "--\n")

	return b.Bytes()
}
