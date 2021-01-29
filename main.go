package main

import (
	"errors"
	"fmt"
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
	defaultNet   = "tcp"
	defaultIp    = "192.168.1.101"
	defaultPort  = "8181"
)

// заголовок
func _beg() {
	fmt.Println("-------------------------------------")
	fmt.Println("|        	  'Go' сервер            |")
	fmt.Println("|     Отвечаем незамедлительно!     |")
	fmt.Println("|                                   |")
	fmt.Println("|   (c) jiliaevyp@gmail.com         |")
	fmt.Println("-------------------------------------")
}

// установка конфигурации сервера
func _config() (string, string) {
	var err int
	var ip, port, net string
	err = 1
	for err == 1 {
		net, err = inpNetwork() // ввод сети сервера
		fmt.Println("Тип сети:   ", net, "\n")
	}
	err = 1
	for err == 1 {
		ip, err = inpIP() // ввод IP сервера
		fmt.Println("IP адрес сервера:   ", ip, "\n")
	}
	err = 1
	for err == 1 {
		port, err = inpPort() // ввод порта сервера
		fmt.Println("Порт сервера:   ", port, "\n")
	}
	return net, ip + port
}

// ввод типа протокола сети
func inpNetwork() (string, int) {
	var typNet string
	var err int
	len := 256
	err = 1
	for err == 1 {
		fmt.Print("Тип сети по умолчанию:	", defaultNet, "\n", " Нажмите (Y) для изменения ")
		yes := yesNo()
		if yes == 0 {
			typNet = defaultNet
			return typNet, 0
		} else {
			fmt.Print("Введите тип сети:	")
			data := make([]byte, len)
			n, err := os.Stdin.Read(data)
			typNet = string(data[0 : n-1])
			if err != nil || typNet != "tcp" {
				fmt.Println(ErrInvalidTypeNetwork)
				return typNet, 1
			}
			return typNet, 0
		}
	}
	return typNet, 0
}

// ввод ip адреса сервера
func inpIP() (string, int) {
	data := ""
	err := 1
	for err == 1 {
		fmt.Print("IP адрес сервера по умолчанию:	", defaultIp, "\n", "Для изменения нажмите 'Y' ")
		yes := yesNo()
		if yes != 1 {
			data = defaultIp
			err = 0
		} else {
			fmt.Print("Введите IP адрес сервера:	")
			fmt.Scanf(
				"%s\n",
				&data,
			)
			iperr := net.ParseIP(data)
			if iperr == nil {
				fmt.Println(ErrInvalidIPaddress)
				return data, 1
			} else {
				err = 0
			}
		}
	}
	return data, err
}

//ввод номера порта
func inpPort() (string, int) {
	var (
		webPort string
	)
	err := 1
	for err == 1 {
		fmt.Print("Порт по умолчанию:	", defaultPort, "\n", "Для изменения нажмите 'Y' ")
		yes := yesNo()
		if yes != 1 {
			webPort = defaultPort
			err = 0
		} else {
			fmt.Print("Введите порт:	")
			fmt.Scanf(
				"%s\n",
				&webPort,
			)
			res, err1 := strconv.ParseFloat(webPort, 16)
			res = res + 1
			err = 0
			if err1 != nil {
				fmt.Println(ErrInvalidPort)
				return ":" + webPort, 1
			}
		}
	}
	return ":" + webPort, 0
}

// проверка на ввод  'Y = 1
func yesNo() int {
	var yesNo string
	len := 4
	data := make([]byte, len)
	n, err := os.Stdin.Read(data)
	yesNo = string(data[0 : n-1])
	if err == nil && (yesNo == "Y" || yesNo == "y" || yesNo == "Н" || yesNo == "н") {
		return 1
	} else {
		return 0
	}
}

// Accept() (принимает входящее подключение) и Close() (закрывает подключение)
//В случае успешного выполнения функция возвращает объект интерфейса net.Listener,
//который представляет функционал для приема входящих подключений.
//В зависимости от типа используемого протокола возвращаемый объект
//Listener может представлять тип net.TCPListener или net.UnixListener
//(оба этих типа ре   ализуют интерфейс net.Listener).

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
		fmt.Println("Получен запрос: " + requestMessage)
		answerText := "Все ништяк!"
		conn.Write([]byte(answerText)) //транслируем в порт ответ ввиде среза байт
		fmt.Println("Отправлен ответ: ", answerText)
	}
}

func main() {
	var yes int
	_beg()                      // заголовок
	_network, _port = _config() // конфигурация сервера
	go _server()                // запуск сервера
	//<-time.After(5 * time.Second)
	if errServ != 0 { // ошибка сервера
		fmt.Println(ErrInvalidServerListen)
		fmt.Println("\n", "Сервер остановлен")
		return
	} else {
		fmt.Println("\n", "Закончить работу сервера? (Y)")
		yes = yesNo()
		if yes == 1 {
			fmt.Println("\n", "Сервер остановлен")
			fmt.Println("\n", "Рад был на Вас пработать!")
			fmt.Print("Обращайтесь в любое время без колебаний!", "\n", "\n")
			return
		}
	}
}
