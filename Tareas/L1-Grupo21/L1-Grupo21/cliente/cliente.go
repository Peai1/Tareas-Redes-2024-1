package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type ConfigInicial struct {
	TableroJugador [2][2]string // A, B, C, D
	IP             string
	Puerto         string
}

type TableroActualizado struct {
	Letra   string
	Tablero [2][2]string
}

func imprimirTablero(tablero [2][2]string) {
	fmt.Println("-------------")
	for i := 0; i < 2; i++ {
		fmt.Print("| ")
		for j := 0; j < 2; j++ {
			fmt.Print(tablero[i][j], " | ")
		}
		fmt.Println()
		fmt.Println("-------------")
	}
}

func conexionUDP() ConfigInicial {
	udpAddr, _ := net.ResolveUDPAddr("udp", "localhost:8080")

	// Crea un socket UDP
	conn, _ := net.DialUDP("udp", nil, udpAddr)

	defer conn.Close()

	message := "Hola! queremos jugar"
	_, _ = conn.Write([]byte(message))

	// Esperar mensaje del servidor
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)

	// Decodificar mensaje del servidor
	var juego ConfigInicial
	json.Unmarshal(buffer[:n], &juego)

	// imprime todos los campos de juego
	fmt.Println("Datos recibidos por conexión UDP:")
	fmt.Println("Tablero del cliente:")
	imprimirTablero(juego.TableroJugador)
	fmt.Println("IP:", juego.IP)
	fmt.Println("Puerto:", juego.Puerto)

	return juego
}

func conexionTCP(juego ConfigInicial) {
	dirrecionTCP := juego.IP + ":" + juego.Puerto
	conn, err := net.Dial("tcp", dirrecionTCP)
	if err != nil {
		fmt.Println("Error al conectarse al servidor TCP")
	}

	tableroServidor := [2][2]string{
		{"A", "B"},
		{"C", "D"},
	}

	defer conn.Close()

	fmt.Println("Empieza el juego.")
	fmt.Println("------------------")

	for {
		fmt.Println("Tablero del cliente:")
		imprimirTablero(juego.TableroJugador)
		fmt.Println("Tablero del servidor:")
		imprimirTablero(tableroServidor)

		// Solicita al usuario que ingrese una letra A, B, C o D
		var letra string
		fmt.Print("Cliente ingrese una letra para atacar al servidor (A, B, C o D): ")
		fmt.Scanln(&letra)

		switch letra {
		case "A":
			tableroServidor[0][0] = "X"
		case "B":
			tableroServidor[0][1] = "X"
		case "C":
			tableroServidor[1][0] = "X"
		case "D":
			tableroServidor[1][1] = "X"
		}

		conn.Write([]byte(letra))

		buffer := make([]byte, 1024)
		n, _ := conn.Read(buffer)

		var tableroActualizado TableroActualizado
		json.Unmarshal(buffer[:n], &tableroActualizado)

		if tableroActualizado.Letra == "2" {
			fmt.Println("¡El servidor ganó, haz perdido!")
			fmt.Println("Tablero final:")
			imprimirTablero(tableroActualizado.Tablero)
			return
		}
		if tableroActualizado.Letra == "1" {
			fmt.Println("¡El cliente ganó!")
			fmt.Println("Tablero final:")
			imprimirTablero(tableroActualizado.Tablero)
			return
		}

		fmt.Println("El servidor te atacó en la letra:", tableroActualizado.Letra)

		juego.TableroJugador = tableroActualizado.Tablero
	}
}

func main() {
	juego := conexionUDP()

	conexionTCP(juego)
	fmt.Println("Fin del juego, conexión con servidor terminada.")
}
