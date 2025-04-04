package headers

import (
	"bytes"
	"errors"
	"fmt"
)

type Headers map[string]string 
const CRLF = "\r\n"

func NewHeaders() Headers{
	return make( map[string]string)
}
func (h Headers) Parse(data [] byte) (n int, done bool, err error){
// based upon tests, a valid header will not start with any spacing
// and must end with \r\n\r\n
	fmt.Printf("Parsing header data: %s...\r\n",string(data))
	
	nbytes := 0
	crlfidx := bytes.Index(data,[]byte(CRLF))
	
	if crlfidx == 0{
		return 2,true, nil
	
	}

	endIdx := bytes.Index(data, []byte(CRLF))
    if endIdx == -1 {
        // Not enough data yet
        return 0, false, nil
    }

	idx := bytes.Index(data,[]byte(":"))
	fmt.Printf("  **** Index of : is %v\n",idx)
	if idx == -1 || idx < 4{
		return 0,false, errors.New("header missing colon separator")
	}	
	
	if data[idx - 1] == ' '{
		return 0,false, errors.New("space before colon separator not permitted")
	}

		


		
		key := string(bytes.TrimSpace(data[:idx]))
		val := string(bytes.TrimSpace(data[idx+1:endIdx]))

		nbytes += len(key) + len(val)

		h[key]=val
		//newdata:=make([]byte,len(data)-len(CRLF))
		//copy(newdata,data[endidx:])
		//data = newdata

		fmt.Printf("new data value is %s\r\nbytes read: %v",string(data),nbytes)

		return endIdx +2,false, nil
	}

	


