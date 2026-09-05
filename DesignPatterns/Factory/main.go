package main

import "fmt"

type Notification interface {
	Send(string)
}

type Sms struct {}

func (e Sms) Send(message string) {
	fmt.Println("SMS message: ", message)
}

type Email struct {}

func (e Email) Send(message string) {
	fmt.Println("Email message: ", message)
}

func NewNotification(typ string) Notification {
	switch typ {
	case "email":
		return Email{}
	case "sms":
		return Sms{}
	default:
		return Email{}
	}
}

func main() {
	n1 := NewNotification("email")
	n2 := NewNotification("sms")
	n3 := NewNotification("")

	n1.Send("hi")
	n2.Send("hi")
	n3.Send("hi")
}