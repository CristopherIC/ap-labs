// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

//!+broadcaster
type client chan<- string // an outgoing message channel

var (
	users = make(map[string]string) 
	alias         = make(map[string]net.Conn) 
	connectionLog = make(map[string]string)   
	entering      = make(chan client)
	leaving       = make(chan client)
	messages      = make(chan string) 
	admin         = ""
	usrCount      = 0
)

func broadcaster() {
	clients := make(map[client]bool) 
	remitter := ""
	message := ""
	normal := false
	for {
		normal = false
		select {
		case msg := <-messages:
			values := []rune(msg)
			a := strings.Index(msg, ">")
			if a >= 0 {
				message = string(values[a+2:])
				remitter = string(values[:a])
				normal = true
			}
			//Sub-commands cases
			msgInfo := strings.SplitN(message, " ", 2)
			if len(msgInfo[0]) > 0 {
				if msgInfo[0][0] == '/' {
					switch msgInfo[0] {

					case "/users":
						usersList(remitter)

					case "/msg":
						if len(msgInfo) < 2 {
							break
						}
						messageData := strings.SplitN(msgInfo[1], " ", 2)
						if len(msgInfo) < 2 {
							break
						}
						privateMsg(remitter, messageData[0], messageData[1])

					case "/time":
						timeInfo(remitter)

					case "/user":
						if len(msgInfo) < 2 {
							break
						}
						userInfo(remitter, msgInfo[1])

					case "/kick":
						if len(msgInfo) < 2 {
							break
						}
						if alias[remitter].RemoteAddr().String() == admin {
							kickUser(msgInfo[1])

							for cli := range clients {
								cli <- "irc-server> "+ "[" + msgInfo[1] + "] was kicked from channel for bad language policy violation"
							}
						}
					}

				} else {
					for cli := range clients {

						if normal {
							cli <- remitter + " > " + message
						} else {
							cli <- msg
						}
					}

				}

			} else {
				// Broadcast message to all
				for cli := range clients {

					if normal {
						cli <- remitter + " > " + message
					} else {
						cli <- msg
					}
				}

			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}
// /users
func usersList(to string) { 
	for usr := range alias {
		privateMsg("irc-server", to, usr+" - connected since "+connectionLog[usr]+"\n")
	}

}

// /time
func timeInfo(to string) { 
	location, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		panic(err)
	}
	time := time.Now().In(location)
	privateMsg("irc-server", to, "Local Time: "+time.Location().String()+time.Format(" 15:04"))

}

// /user <username>
func userInfo(to, usr string) { 
	privateMsg("irc-server", to, "username: "+usr+" IP: "+strings.Split(alias[usr].LocalAddr().String(), ":")[0]+" Connected since: "+connectionLog[usr])

}

// /kick <username>
func kickUser(usr string) { 
	a := alias[usr]
	if a != nil {
		alias[usr].Close()
		fmt.Printf("irc-server> [%s] was kicked\n", usr)
		delete(users, alias[usr].RemoteAddr().String())
		delete(alias, usr)

	}
}

// /msg <user> <msg>
func privateMsg(from, to, message string) { 
	ch := make(chan string)
	go clientWriter(alias[to], ch)
	ch <- from + " > " + message
	close(ch)
	return

}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)
	time := time.Now()

	who := conn.RemoteAddr().String()
	input := bufio.NewScanner(conn)
	input.Scan()
	userName := input.Text()
	usrCount++
	users[who] = userName
	alias[userName] = conn
	connectionLog[userName] = time.Format("2006-01-02 15:04:05")
	fmt.Printf("irc-server > New connected user [%s]\n", users[who])
	ch <- "irc-server > Welcome to the Simple IRC Server"
	ch <- "irc-server > Your user [" + users[who] + "] is successfully logged"
	//Make the first user admin
	if admin == "" {
		admin = who
		ch <- "irc-server > Congrats, you were the first user"
		ch <- "irc-server > You're the new IRC Server ADMIN"
		fmt.Printf("irc-server > [%s] was promoted as the channel ADMIN\n", users[who])
	}
	//Msg for New Users
	messages <- "New Connected user [" + users[who] +"]"
	entering <- ch

	for input.Scan() {
		messages <- users[who] + "> " + input.Text()

	}

	//Leave
	leaving <- ch
	fmt.Printf("irc-server > [%s] left\n", users[who])

	a := alias[users[who]]
	if a != nil {
		delete(alias, users[who])
		delete(users, who)

	}

	if who == admin {
		for nam := range alias {
			privateMsg("irc-server", nam, "You're the new IRC Server ADMIN")
			fmt.Printf("irc-server > [%s] was promoted as the channel ADMIN\n", users[alias[nam].RemoteAddr().String()])
			admin = alias[nam].RemoteAddr().String()

		}
	}

	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", os.Args[2]+":"+os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("irc-server > Simple IRC Server started at %s\n", listener.Addr())

	go broadcaster()
	fmt.Print("irc-server > Ready for receiving new clients\n")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}

}

//!-main