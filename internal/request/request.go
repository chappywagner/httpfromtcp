package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const bufferSize = 8
const initialized = 0
const done = 1

type Request struct {
	state int
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(data []byte)(*RequestLine, int, error){
// this function gets the requestline from the entire byte array
// and returns a requestline struct
	
	new_line_idx := bytes.Index(data, []byte("\r\n"))
	
	if new_line_idx == -1{
		return nil, 0, nil
	}
	
	request_text := string(data[:new_line_idx])
	
	byte_len := len(request_text) + 2

	requestLine, err := requestLineFromString(request_text)
	
	if err!=nil{
		return nil, byte_len, err
	}

	return requestLine,byte_len,nil

}

func requestLineFromString(str string)(*RequestLine,error){

	request_line := strings.Split(str," ")

	if len(request_line) !=3{
		return nil, errors.New("invalid request")
	}

	valid := true

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

	return &RequestLine{ 
			HttpVersion: http_version,
			Method: method,
			RequestTarget: target,
			},nil
		
		
}

func RequestFromReader(reader io.Reader) (*Request, error){
	
	buf := make([]byte, bufferSize)
	
	readToIndex := 0

	parser := Request{state: initialized,}

	for parser.state != done{
		
		if readToIndex >= len(buf){

			new_buf := make([]byte,len(buf)*2)
			copy(new_buf,buf)
			buf = new_buf
		
		}
	
		// read starting at current read index all the way to the end
		nbytesRead,err:= reader.Read(buf[readToIndex:])

		// if we hit the end of the file
		if err == io.EOF{
		
			parser.state = done
			break
	
		}

		// now parse the buffer starting at current readToIndex up to readToIndex + bytes read

		nbytesParsed, err := parser.parse(buf[:readToIndex])
		
		if err != nil{

			return nil, err

		}

		// increment readToIndex by bytes read
		readToIndex += nbytesRead

		// shift data to front of the array, removing parsed data
		copy(buf,buf[nbytesParsed:readToIndex])

		// adjust read to index by number of bytes actually parsed
		readToIndex -= nbytesParsed

	}

 	return &parser, nil

}

func (r *Request) parse(data []byte) (int, error){
	
	if r.state == initialized {
	
		requestLine, bytes, err := parseRequestLine(data)
	
		if err != nil {
			return 0, err
		}

		if bytes == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine

		r.state = done
		
		return bytes, nil
	
	}else if r.state == done {
			return 0, errors.New("error: cannot read data in done state")
		}else{
			return 0, errors.New("unknown parser state")
		}
	
}

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
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

	// Test: Good GET Request line
	r, err = RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
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

	// Test: Good GET Request line
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}

	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	reader = &chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}

	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

}

