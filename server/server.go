package main

import (
	"net"
	"fmt"
	"bytes"
)

type Client struct {
	conn     net.Conn
	username string
}

var clients = &Clients{
	make([]*Client, 0),
}

type Clients struct {
	clients []*Client
}

func (c *Clients) SendAll(message string, from *Client) {
	var buffer bytes.Buffer
	buffer.WriteString(from.username)
	buffer.WriteString(" : ")
	buffer.WriteString(message)
	for _, client := range c.clients {
		if client != from {
			client.conn.Write(buffer.Bytes())
		}
	}
}

func (c *Clients) AddClient(newClient *Client) {
	c.clients = append(c.clients, newClient)
}

func (c *Clients) RemoveClient(toRemove *Client) {
	for i, client := range c.clients {
		if toRemove == client {
			c.clients[i] = c.clients[len(c.clients)-1]
			c.clients[len(c.clients)-1] = nil
			c.clients = c.clients[:len(c.clients)-1]
		}
	}
}

func (c *Client) AwaitMessages() {
	buffer := make([]byte, 10000)
	for {
		size, err := c.conn.Read(buffer)
		if err != nil {
			fmt.Println(c.username, " disconnected")
			clients.RemoveClient(c)
			return
		}
		if size > 0 {
			message := string(buffer[:size])
			fmt.Println(c.username, " : ", message)
			clients.SendAll(message, c)
		}
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("A client has connected")
	buffer := make([]byte, 1000)
	size, _ := conn.Read(buffer)
	username := string(buffer[:size])
	client := &Client{
		conn,
		username,
	}
	clients.AddClient(client)
	go client.AwaitMessages()
}

func main() {
	fmt.Println("Starting server on port 8080")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(conn)
	}
}
