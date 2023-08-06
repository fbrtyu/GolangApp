package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"main/modules/ws"
	"net"
	"net/http"
	"os"
	"strings"
)

// Хранение IPv4
var ip string

// Структура для хранения данных в формате JSON
type Message struct {
	Message string `json:"message"`
}

// Функция получения IPv4
func GetIp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write([]byte(ip))
	fmt.Println("Подключение выполнено!")
}

// Функция получения новых данных из JSON
func GetElem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	file_data, err := ioutil.ReadFile("text.txt")
	if err != nil {
		fmt.Println("Файл не читается или не существует!\n", err)
	}
	w.Write([]byte(file_data))
	fmt.Println("Данные об элементах отправлены!")
}

func main() {
	fmt.Println("Server start! localhost:8080")

	// Получения адреса ПК
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip = ipv4.String() + ":8080"
			fmt.Println("В браузере телефона/компьютера перейдите по адресу:", ip)
		}
	}

	// Изменение IPv4 на фронтенде, чтобы иметь возможность подключиться к серверу с телефона
	file_data, err := ioutil.ReadFile("public/static/js/main.76ac0ff9.js")
	var ipp = ip
	modifiedString := strings.Replace(string(file_data), "000.000.000.000:8080", ipp, -1)
	data := []byte(modifiedString)
	e := ioutil.WriteFile("public/static/js/main.76ac0ff9.js", data, 0600)
	if e != nil {
		fmt.Println("Не могу создать файл!\n", e)
	}

	if err != nil {
		fmt.Println("Файл не читается или не существует!\n", err)
	}

	http.HandleFunc("/getip", GetIp)
	http.HandleFunc("/getelem", GetElem)

	// Создание файлового сервера
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// Конфигурация вебсокет соединения
	http.HandleFunc("/ws", ws.HandleConnections)

	// Запуск горутины для работы с данными
	go ws.HandleMessages()

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
