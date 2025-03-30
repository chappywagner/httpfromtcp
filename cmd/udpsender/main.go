package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main(){
	n,err:=	net.ResolveUDPAddr("udp","localhost:42069")
	
	if err !=nil{
		return
	}

	uconn,err:=net.DialUDP("udp",nil,n)

	if err!=nil{
		log.Fatalf("error connecting due to %s",err.Error())
	}

	defer uconn.Close()

	reader:=bufio.NewReader(os.Stdin)

	for{
		
		fmt.Println(">")
		
		str, err:= reader.ReadString('\n')

		if err!=nil{
			log.Fatalf("error %s\n",err.Error())
		}

		n, err:= uconn.Write([]byte(str))

		if err!= nil{
			log.Fatalf("error %s\n trying to write %v bytes",err.Error(),n)
		}

	}
}