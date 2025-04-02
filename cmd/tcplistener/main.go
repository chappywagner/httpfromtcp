package main

import (
	"fmt"
	"internal/request"
	"io"
	"log"
	"net"
	"strings"

	"github.com/golang-jwt/jwt/v5/request"
)

const port = ":42069"

func main() {
 
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		linesChan := request.RequestFromReader(conn)

		fmt.Println("Request line:\n")
		fmt.Printf("- Method: %s\n-Target: %s/\n-Version: %s",&linesChan.Method,&linesChan.Target,&linesChan.Version)
		//for line := range linesChan {
		//		fmt.Println(line)
		//		}
		//	fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}

}

func getLinesChannel(f io.ReadCloser) (<-chan string){
	
	
	str_chan := make(chan string)
	
	go func(){
		
		defer close(str_chan)
		

		defer f.Close()
		
		current_line:=""

		b:=make([]byte,8)

		for {

		n,err := f.Read(b)
		
		if err==io.EOF{
			if current_line !=""{
				str_chan <- current_line
			}
			
			break	
						
		}
		
		if err!=nil{
			break
		}

		s:=string(b[:n])

		parts := strings.Split(s,"\n")

		current_line += parts[0]

		if (len(parts)>1){
		  str_chan <- current_line
		  current_line=parts[1]
		}
	}
   }()

	return str_chan
}


