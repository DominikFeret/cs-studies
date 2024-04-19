package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"server/internal/database"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	db         *database.Queries
}

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable not set")
	}

	databaseURL := os.Getenv("DB_URL")
	if databaseURL == "" {
		log.Fatal("DB_URL environment variable not set")
	}

	dbConn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}
	defer dbConn.Close()

	server := NewServer(":"+portString, database.New(dbConn))

	server.Start()

}

func NewServer(listenAddr string, db *database.Queries) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		db:         db,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.Serve()
	<-s.quitch

	return nil
}

func (s *Server) Stop() {
	close(s.quitch)
	s.ln.Close()
}

func (s *Server) Serve() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.quitch:
				return
			default:
				log.Println("Error accepting connection: ", err)
			}
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading from connection: ", err)
		return
	}

	action := string(buf[:n])
	switch action {
	case "create":
		s.handleCreate(conn)
	case "login":
		s.handleLogin(conn)
	default:
		s.handleInvalidAction(conn)
	}

}

func (s *Server) handleCreate(conn net.Conn) {
	log.Println("Create user")

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading from connection: ", err)
		return
	}
	username := string(buf[:n])
	log.Println("Received username: ", username)

	_, err = s.db.GetUserByName(context.Background(), username)
	// ignore sql.ErrNoRows, as it means the name is not occupied and all's good
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error getting user: ", err)
		return
	}
	// all errors have been handled, so if err is nil, the name is already occupied
	if err == nil {
		conn.Write([]byte("exists"))
		return
	}
	conn.Write([]byte("ok"))

	n, err = conn.Read(buf)
	if err != nil {
		log.Println("Error reading from connection: ", err)
		return
	}
	pwd := string(buf[:n])
	log.Println("Received password: ", pwd)

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
		Password:  pwd,
	}

	_, err = s.db.CreateUser(context.Background(), user)
	if err != nil {
		log.Println("Error creating user: ", err)
		return
	}
	conn.Write([]byte("ok"))
}

func (s *Server) handleLogin(conn net.Conn) {
	log.Println("Login user")

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading from connection: ", err)
		return
	}
	conn.Write([]byte("ok"))

	username := string(buf[:n])
	log.Println("Received username: ", username)

	user, err := s.db.GetUserByName(context.Background(), username)
	if err == nil {
		conn.Write([]byte(user.Password))
	} else if err.Error() == "sql: no rows in result set" {
		conn.Write([]byte("notfound"))
		return
	}
	if err != nil {
		log.Println("Error getting user: ", err)
		return
	}
}

func (s *Server) handleInvalidAction(conn net.Conn) {
	log.Println("Invalid action")
	conn.Write([]byte("invalid"))
}
