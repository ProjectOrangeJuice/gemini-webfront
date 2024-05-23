package main

import "time"

type chatFile struct {
	When     time.Time
	Title    string
	Messages []chat
}

type chat struct {
	Who     string
	Content string
}

type chats struct {
	Token string
	Title string
	When  time.Time
}
