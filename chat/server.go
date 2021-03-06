package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type User struct {
	nick string
	conn net.Conn
}

var (
	users []*User
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":1200")
	if err != nil {
		log.Fatalln(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	nick := getUserNickname(conn)
	newUser := bindNewUser(nick, conn)
	handleMsgsFromUser(newUser)
}

func getUserNickname(conn net.Conn) string {
	var buf = make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	nick := strings.TrimSpace(string(buf[:n]))
	fmt.Printf("%s joined chat\n", nick)
	return nick
}

func bindNewUser(nick string, conn net.Conn) *User {
	user := &User{
		nick: nick,
		conn: conn,
	}
	users = append(users, user)
	return user
}

func handleMsgsFromUser(user *User) {
	var buf = make([]byte, 512)
	for {
		n, err := user.conn.Read(buf)
		if err != nil {
			log.Fatalln(err)
		}
		msg := strings.TrimSpace(string(buf[:n]))
		fmt.Printf("Received msg from %s: %s\n", user.nick, msg)
		go broadcast(msg, user.nick)
	}
}

func broadcast(msg, nick string) {
	for _, user := range users {
		if user.nick == nick {
			continue
		}
		fmt.Println("Broadcasting message to", user.nick)
		user.conn.Write([]byte(msg + "\n"))
	}
}
