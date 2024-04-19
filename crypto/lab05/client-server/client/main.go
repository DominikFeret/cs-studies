package main

import (
	"errors"
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func main() {
	err := connect()
	for err != nil && err.Error() != "EOF" {
		time.Sleep(1 * time.Second)
		err = connect()
	}
}

func connect() error {
	conn, err := net.Dial("tcp", "localhost:"+"8080")
	if err != nil {
		fmt.Println("Couldn't connect to server:", err)
		return err
	}
	defer conn.Close()

	var action int
	buf := make([]byte, 1024)

	// in case of invalid input, the user will be prompted again
	// input clearup is implemented on the server side
	fmt.Print("Choose an action:\n1. Create account\n2. Login\n3. Exit\n")
	fmt.Scanln(&action)
	switch action {
	case 1:
		err = createAcc(conn)
	case 2:
		err = logIn(conn)
	case 3:
		fmt.Println("Exiting...")
		return nil
	default:
		fmt.Println("Invalid action")
		return errors.New("invalid action")
	}
	if err != nil {
		return err
	}

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return err
	}
	fmt.Println(string(buf[:n]))

	return nil
}

func createAcc(conn net.Conn) error {
	conn.Write([]byte("create"))

	var username, password string

	fmt.Println("Username:")
	fmt.Scanln(&username)
	if username == "" {
		fmt.Println("Username can't be empty")
		return errors.New("username can't be empty")
	}
	conn.Write([]byte(username))
	resp := readResponse(conn)
	if resp == "exists" {
		fmt.Println("Username already exists")
		return errors.New("username already exists")
	}
	if resp == "empty" {
		fmt.Println("Username can't be empty")
		return errors.New("username can't be empty")
	}

	fmt.Println("Password:")
	password = getPwd()
	if password == "" {
		fmt.Println("Password can't be empty")
		return errors.New("password can't be empty")
	}
	conn.Write([]byte(password))

	r := readResponse(conn)
	if r == "ok" {
		fmt.Println("Account created successfully!")
	}

	return nil
}

func logIn(conn net.Conn) error {
	conn.Write([]byte("login"))

	var username, password string

	fmt.Println("Username:")
	fmt.Scanln(&username)
	conn.Write([]byte(username))
	resp := readResponse(conn)
	if resp == "notfound" {
		fmt.Println("Username not found")
		return errors.New("username not found")
	}

	fmt.Println("Password:")
	password = getPwdPlain()
	resp = readResponse(conn)
	if resp == "notfound" {
		fmt.Println("Username not found")
		return errors.New("username not found")
	}
	if bcrypt.CompareHashAndPassword([]byte(resp), []byte(password)) != nil {
		fmt.Println("Invalid password")
		return errors.New("invalid password")
	}
	fmt.Println("Logged in successfully!")

	return nil
}

func readResponse(conn net.Conn) string {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return ""
	}

	return string(buf[:n])
}

func getPwdPlain() string {
	pwd, err := term.ReadPassword(0)
	if err != nil {
		fmt.Println("Error reading password:", err)
	}
	return string(pwd)
}

func getPwd() string {
	pwd, err := term.ReadPassword(0)
	if err != nil {
		fmt.Println("Error reading password:", err)
	}
	return string(hashAndSalt(pwd))
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
	}
	return string(hash)
}
