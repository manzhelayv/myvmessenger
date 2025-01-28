package servers

import (
	"chat/models"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"time"
)

var PORT = ":1235"
var Users = make(map[string]*User, 0)

const PING = "pingWS123443321"
const PONG = "pongWS123443321"

type User struct {
	UserFrom   string
	UserTo     string
	Message    []byte
	Initiated  bool
	UChannel   chan []byte
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
	Connection *websocket.Conn
}

type Kafka struct {
	UserFrom string
	Message  []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ConnectionManager struct {
	name      string
	initiated bool
}

type MessageWithUserTo struct {
	Message string
	userTo  string
}

func (u *User) Listen() {
	u.InfoLog.Println("Прослушивание для пользователя", u.UserFrom)

	for {
		select {
		case msg := <-u.UChannel:
			if u.UserFrom != u.UserTo {
				err := u.Connection.WriteMessage(1, msg)
				if err != nil {
					u.ErrorLog.Println(err)
				}
			}
		}
	}
}

func (cM *ConnectionManager) Listen(ws *websocket.Conn, userFrom string, h *HttpServer) {
	user := User{}

	user = User{UserFrom: userFrom, Initiated: false}
	Users[userFrom] = &user

	uChan := make(chan []byte, 100)

	user.Initiated = true
	user.UChannel = uChan
	user.Connection = ws
	user.InfoLog = h.InfoLog
	user.ErrorLog = h.ErrorLog

	go cM.messageReady(ws, &user, h)
}

func (h *HttpServer) pingPong(user *User, status string) {
	userStatusResponse, err := h.manChat.GetStatusUser(user.UserFrom)
	if err != nil {
		h.ErrorLog.Println(err)
	}

	message, err := json.Marshal(`{"Message":"` + PONG + `","UserTo":"` + user.UserFrom + `","Status":"` + status + `"}`)
	if err != nil {
		h.ErrorLog.Println(err)
	}

	if userStatusResponse != nil {
		for userTo, _ := range userStatusResponse {
			if _, ok := Users[userTo]; ok {
				Users[userTo].UChannel <- message
			}
		}
	}

	Users[user.UserFrom].UChannel <- message
}

func (cM *ConnectionManager) messageReady(ws *websocket.Conn, user *User, h *HttpServer) {
	go user.Listen()

	ctx := context.Background()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			h.ErrorLog.Println("ReadMessage websocket", err)
			break
		}
		log.Println("message", string(message))
		newMessage := &models.MessageWS{}
		if err := json.Unmarshal(message, newMessage); err != nil {
			continue
		}

		if newMessage.Message == PING {
			h.pingPong(Users[user.UserFrom], newMessage.Status)
			continue
		}

		kafkaMessage := Kafka{
			UserFrom: user.UserFrom,
			Message:  message,
		}
		log.Println("message UserFrom", string(user.UserFrom))
		userByte, err := json.Marshal(&kafkaMessage)
		if err != nil {
			h.ErrorLog.Println("Marshal for kafka", err)
			continue
		}

		err = h.kafkaProducer.WriteMessages(ctx, kafka.Message{
			Key:   []byte(time.Now().Format(time.RFC3339Nano)),
			Value: userByte,
		})
		if err != nil {
			h.ErrorLog.Println("WriteMessages производитель", err)
		}
	}
}

func Initiate() *ConnectionManager {
	return &ConnectionManager{
		name:      "Чат сервер",
		initiated: false,
	}
}

func (h *HttpServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	connManage := Initiate()

	userFrom := r.URL.Query().Get("user_from")

	h.InfoLog.Println("Подключение из:", r.Host)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.InfoLog.Println("upgrader.Upgrade:", err)
		return
	}

	go connManage.Listen(ws, userFrom, h)
}

func (h *HttpServer) StartWS() {
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         PORT,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	mux.Handle("/", http.HandlerFunc(h.wsHandler))

	h.InfoLog.Println("Прослушивание TCP порта", PORT)
	log.Println("Прослушивание TCP порта", PORT)

	err := s.ListenAndServe()
	if err != nil {
		h.ErrorLog.Println(err)
		return
	}
}
