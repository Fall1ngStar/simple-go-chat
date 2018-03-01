package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
)

func awaitMessage(conn net.Conn) {
	buffer := make([]byte, 10000)
	for {
		size, _ := conn.Read(buffer)
		if size > 0 {
			fmt.Println(string(buffer[:size]))
		}
	}
}


func main() {
	fmt.Print("Enter your username : ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := scanner.Text()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to server !")
	go awaitMessage(conn)

	conn.Write([]byte(username))
	for scanner.Scan(){
		message := scanner.Text()
		fmt.Println("You : ", message)
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println(err)
		}
	}
}
