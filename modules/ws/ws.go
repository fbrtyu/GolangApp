package ws

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/micmonay/keybd_event"
)

// Структура для хранения данных в формате JSON и обмена
type Message struct {
	Message string `json:"message"`
}

// Канал передачи данных
var broadcast = make(chan Message)

// Массив подключенных клиентов
var clients = make(map[*websocket.Conn]bool)
var clientMass = make(map[string]*websocket.Conn)

// Проверка возможности соединения по вебсокет
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Обновление соединения
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Регистрация нового клиента
	clients[ws] = true
	clientMass[r.Host] = ws

	// Структуры для данных
	var msg Message
	// Чтение данных от пользователя и передача их в канал broadcast
	for {
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		ws.WriteJSON(msg)
		broadcast <- msg
	}
}

func HandleMessages() {
	// Чтение сообщений из канала данных
	var msg Message
	for {
		msg = <-broadcast
		if len(msg.Message) > 30 {
			_, err := ioutil.ReadFile("text.txt")
			if err != nil {
				fmt.Println("Файл не читается или не существует!\n", err)
			}
			modifiedString := msg.Message
			data := []byte(modifiedString)
			e := ioutil.WriteFile("text.txt", data, 0600)
			if e != nil {
				fmt.Println("Не могу создать файл!\n", e)
			}
			fmt.Println("Данные об элементах сохранены!")
		} else {
			kb, err := keybd_event.NewKeyBonding()
			if err != nil {
				panic(err)
			}

			/* // Запуск обычного блокнота
			path, err := exec.LookPath("Тут можно указать путь к любой программе")
			if err != nil {
				fmt.Println("Файл не найден")
			}
			fmt.Printf("Доступ к файлу %s\n", path)
			cmd := exec.Command("notepad.exe") // Вместо path для примера используется notepad.exe
			cmd.Run() */

			// Виртуальное нажатие клавиши на ПК
			if msg.Message == "[" {

				kb.SetKeys(keybd_event.VK_SP4)
				err = kb.Launching()
				if err != nil {
					panic(err)
				}
			}

			if msg.Message == "]" {

				kb.SetKeys(keybd_event.VK_SP5)
				err = kb.Launching()
				if err != nil {
					panic(err)
				}
			}

			//Запуск программы
			prog := msg.Message[0]
			if string(prog) == "*" {
				msg.Message = strings.Replace(msg.Message, "*", "", -1)
				path, err := exec.LookPath(msg.Message)
				if err != nil {
					fmt.Println("Файл не найден")
				}
				fmt.Println(msg.Message)
				fmt.Printf("Доступ к файлу %s\n", path)
				cmd := exec.Command(path)
				cmd.Run()
			}
		}
	}
}
