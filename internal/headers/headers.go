package headers

import (
	"bytes"
	"errors"
	"fmt"
)

type Headers map[string]string 
const CRLF = "\r\n\r\n"

func NewHeaders() Headers{
	return make( map[string]string)
}
func (h Headers) Parse(data [] byte) (n int, done bool, err error){
// based upon tests, a valid header will not start with any spacing
// and must end with \r\n\r\n
	nbytes := 0
	for !done{
		
	if len(data)==0{
		return 0, true, errors.New("no data")
	}

	idx := bytes.Index(data,[]byte(":"))
	fmt.Printf("  **** Index of : is %v\n",idx)
	if idx == -1 || idx < 4{
		return 0,false, errors.New("improperly formatted header")
	}	
	
	if data[idx - 1] == ' '{
		return 0,false, errors.New("improperly formatted header")
	}

		//
		endidx := bytes.Index(data,[]byte(CRLF))

		if endidx == -1{
			return 0, false, nil
		}


		if endidx == 0{
			done = true
			break
		}


		sepidx := bytes.Index(data,[]byte(": "))
		key := string(bytes.Trim(data[:sepidx]," "))
		val := string(bytes.Trim(data[sepidx:endidx]," "))
		nbytes += len(key) + len(val)

		h[key]=val
		newdata:=make([]byte,len(data)-len(CRLF))
		copy(newdata,data[endidx:endidx+len(CRLF)])
		data = newdata

		fmt.Printf("new data value is %s\r\nbytes read: %v",string(data),nbytes)

	}

	return nbytes,done,nil

}
