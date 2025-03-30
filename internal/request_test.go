package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error){
	
	request,err:=io.ReadAll(reader)
	
	if err!=nil{
		return nil, err
	}
	
	content:= string(request)
	
	lines:= strings.Split(content,"\r\n")
	
	if len(lines)==0{
		return nil, errors.New("no data found")
	}

	request_line := strings.Split(lines[0]," ")

	if len(request_line) !=3{
		return nil, errors.New("invalid request")
	}
	valid:=true
	for _,s := range request_line[0]{
		if unicode.IsUpper(s)!= true{
			valid = false
			break
		}
	}
	if !valid{
		return nil, errors.New("invalid method type")
	}
	method:=request_line[0]
	
	target:=request_line[1]

	if request_line[2]!="HTTP/1.1"{
		return nil, errors.New("invalid http version, HTTP 1.1 only supported")
	}
	http_version:= strings.Replace(request_line[2],"HTTP/","",1)

	req := &Request{
		RequestLine: RequestLine{ 
			HttpVersion: http_version,
			Method: method,
			RequestTarget: target},
	}
	fmt.Println(lines[0])
 	return req,nil
}

func TestRequestLineParse(t *testing.T) {
	
	
	// Test: Good GET Request line
r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "GET", r.RequestLine.Method)
assert.Equal(t, "/", r.RequestLine.RequestTarget)
assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

// Test: Good GET Request line with path
r, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "GET", r.RequestLine.Method)
assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

// Test: Invalid number of parts in request line
_, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
require.Error(t, err)
}
