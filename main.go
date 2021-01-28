package main

import (
	"errors"
	"fmt"
	"time"

	//"io"
	"net"
	"os"
	"strconv"
)

var requestMessage, _network, _ip, _port string
var errServ int
var ErrInvalidTypeNetwork = errors.New("invalid type network")
var ErrInvalidPort = errors.New("invalid port number")
var ErrInvalidIPaddress = errors.New("invalid IP address")
var ErrInvalidAnswerServer = errors.New("invalid answer server")
var ErrInvalidServerListen = errors.New("invalid listen server")

const (
	answerServer = "Hello, I am a server."
	readyServer  = "I'm ready!"
)

// выбор протокола сети
//Протокол должен представлять одно из значений: "tcp", "tcp4", "tcp6", "unix", "unixpacket".

func inpNetwork() (string, int) {
	var typNet string
	len := 256
	data := make([]byte, len)
	n, err := os.Stdin.Read(data)
	typNet = string(data[0 : n-1])
	if err != nil || typNet != "tcp" {
		return typNet, 1
	} else {
		return typNet, 0
	}
}

func inpIP() {
	err := 1
	for err == 1 {
		fmt.Print("Введите IP сервера:	")
		//inpIP()
		fmt.Scanf(
			"%s\n",
			&_ip,
		)
		iperr := net.ParseIP(_ip)
		if iperr != nil {
			err = 0
		} else {
			fmt.Println(ErrInvalidIPaddress)
		}
	}
}

// ввод локального адреса
//Локальный адрес может содержать только номер порта, например, ":8080".
//В этом случае приложение будет обслуживать по всем.
func inpPort() (string, int) {
	var (
		webPort string
		res     float64
	)
	fmt.Scanf(
		"%s\n",
		&webPort,
	)
	res, err := strconv.ParseFloat(webPort, 16)
	res = res + 1
	if err != nil || len(webPort) != 4 {
		return ":" + webPort, 1
	} else {
		return ":" + webPort, 0
	}
}

// Accept() (принимает входящее подключение) и Close() (закрывает подключение)
//В случае успешного выполнения функция возвращает объект интерфейса net.Listener,
//который представляет функционал для приема входящих подключений.
//В зависимости от типа используемого протокола возвращаемый объект
//Listener может представлять тип net.TCPListener или net.UnixListener
//(оба этих типа реализуют интерфейс net.Listener).

func _server() {
	errServ = 0
	fmt.Println("Server= ", _network, _port)
	listener, err := net.Listen(_network, _port) // установка тип сети и порта для прослушивания
	if err != nil {
		fmt.Println(ErrInvalidAnswerServer)
		errServ = 1
		return
	}
	defer listener.Close()                                  // включение прослушивание порта
	fmt.Println(answerServer, " network=", _network, _port) //сообщаем что сервер запущен
	fmt.Println(readyServer)
	for {
		conn, err := listener.Accept() // conn <== срез байт из порта
		if err != nil {
			fmt.Println(ErrInvalidServerListen)
			errServ = 1
			return
		}
		go handleConnect(conn) // запуск обработчика запросов сервера
		// нужен таймер отключения
	}
}

// Обработчик и ответчик запросов сервера
func handleConnect(conn net.Conn) {
	defer conn.Close()
	for {
		input := make([]byte, (1024 * 4)) // считываем полученные в запросе срез байт из порта
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			fmt.Println(err)
			break
		}
		requestMessage := string(input[0:n])
		answerText := "На запрос " + requestMessage + "---> Ответ сервера: Все ништяк!"
		conn.Write([]byte(answerText)) //транслируем в порт ответ ввиде среза байт
	}
}

func _beg() {
	fmt.Println("------------------------------------")
	fmt.Println("|  Запуск Go server                |")
	fmt.Println("|  Запускать, не перезапускать!    |")
	fmt.Println("|                                  |")
	fmt.Println("|   (c) jiliaevyp@gmail.com        |")
	fmt.Println("------------------------------------")
}

func main() {
	var komand string
	err := 1
	_beg() // заголовок
	for err == 1 {
		fmt.Print("Введите тип сети:	")
		_network, err = inpNetwork()
		if err == 1 {
			fmt.Println(ErrInvalidTypeNetwork)
		}
	}
	inpIP()
	err = 1
	for err == 1 {
		fmt.Print("Введите номер порта:	")
		_port, err = inpPort()
		if err == 1 {
			fmt.Println(ErrInvalidPort)
		}
	}
	_port = _ip + _port
	go _server() // запуск сервера
	<-time.After(5 * time.Second)
	if errServ == 0 {
		komand = "Y"
		for komand == "Y" || komand == "y" || komand == "Н" || komand == "н" {
			fmt.Println("\n", "Закончить работу сервера? (Y)")
			fmt.Scanf(
				"%s\n",
				&komand,
			)
			if komand == "Y" || komand == "y" || komand == "Н" || komand == "н" {
				fmt.Println("\n", "Рад был с Вами пработать!")
				fmt.Print("Обращайтесь в любое время без колебаний!", "\n", "\n")
				return
			}
		}
	} else {
		fmt.Println(ErrInvalidServerListen)
	}
}
