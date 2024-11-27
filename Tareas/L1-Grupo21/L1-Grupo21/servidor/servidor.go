package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strings"
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

var stopFlag = make(chan bool)


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

func generarTablero(tipo int) [2][2]string {
	tablero := [2][2]string{
		{"A", "B"},
		{"C", "D"},
	}
	i, j := rand.Intn(2), rand.Intn(2)
	if tipo == 1 {
		tablero[i][j] = "1" // Cliente
	} else {
		tablero[i][j] = "2" // Servidor
	}
	return tablero
}

func servidorUDP(tableroCliente [2][2]string) {
	direccionServer, _ := net.ResolveUDPAddr("udp", "localhost:8080")

	conexionUDP, _ := net.ListenUDP("udp", direccionServer)
	defer conexionUDP.Close()

	buffer := make([]byte, 1024)
	for {
		n, dir, _ := conexionUDP.ReadFromUDP(buffer)
		if strings.TrimSpace(string(buffer[0:n])) == "Hola! queremos jugar" {

			fmt.Println("Cliente te envió el mensaje:", string(buffer[0:n]))

			configInicial := ConfigInicial{
				TableroJugador: tableroCliente,
				IP:             "localhost",
				Puerto:         "8080",
			}

			configInicialBytes, _ := json.Marshal(configInicial)
			conexionUDP.WriteToUDP(configInicialBytes, dir)
			break
		}
	}
}

func iniciarTCP(tableroCliente [2][2]string, tableroServidor [2][2]string) {
    dirrecionTCP, _  := net.ResolveTCPAddr("tcp", "localhost:8080")
    tcpListener, _ := net.ListenTCP("tcp", dirrecionTCP)
    defer tcpListener.Close()

	go func() {
        <-stopFlag // Se espera la señal para cerrar el servidor
        tcpListener.Close() // AcceptTCP da un error
    }()

    for {
        tcpConn, err := tcpListener.AcceptTCP()
        if err != nil {
			fmt.Println("Conexión terminada")
            break // Sale del bucle cuando el tcpListener.Close()
        }
        go handleClient(tcpConn, tableroCliente, tableroServidor)
    }
}

func handleClient(conn net.Conn, tableroClienteReal [2][2]string, tableroServidor [2][2]string) {
	defer conn.Close()

	tableroClienteJuego := [2][2]string{
		{"A", "B"},
		{"C", "D"},
	}

	for {
		fmt.Println("Tablero del servidor:")
		imprimirTablero(tableroServidor)
		fmt.Println("Tablero del cliente:")
		imprimirTablero(tableroClienteJuego)

		buffer := make([]byte, 1024)
		n, _ := conn.Read(buffer)
		letraAtaqueCliente := string(buffer[:n])
		fmt.Println("El cliente atacó al servidor en posición:", letraAtaqueCliente)

		k, j := 0, 0
		switch letraAtaqueCliente {
			case "A":
				k, j = 0, 0
			case "B":
				k, j = 0, 1
			case "C":
				k, j = 1, 0
			case "D":
				k, j = 1, 1
			default:
				k, j = 0, 0
		}

		if tableroServidor[k][j] == "2" {
			fmt.Println("¡Servidor ha perdido!")
			tableroFin := TableroActualizado{
				Letra:   "1",
				Tablero: tableroClienteJuego,
			}
			tableroFinBytes, _ := json.Marshal(tableroFin)
			conn.Write(tableroFinBytes) // 1 = Ganó el cliente
			stopFlag <- true
			return
		} else {
			tableroServidor[k][j] = "X"
		}

		fmt.Println("Tablero del servidor:")
		imprimirTablero(tableroServidor)
		fmt.Println("Tablero del cliente:")
		imprimirTablero(tableroClienteJuego)

		var letraAtaqueServidor string
		fmt.Print("Servidor ingrese una letra para atacar al cliente (A, B, C o D): ")
		fmt.Scanln(&letraAtaqueServidor)

		switch letraAtaqueServidor {
			case "A":
				k, j = 0, 0
			case "B":
				k, j = 0, 1
			case "C":
				k, j = 1, 0
			case "D":
				k, j = 1, 1
			default:
				k, j = 0, 0
		}

		if tableroClienteReal[k][j] == "1" {
			fmt.Println("¡Servidor ha ganado!")
			tableroFin := TableroActualizado{
				Letra:   "2",
				Tablero: tableroClienteReal,
			}
			tableroFinBytes, _ := json.Marshal(tableroFin)
			conn.Write(tableroFinBytes)
			stopFlag <- true
			return
		} else {
			tableroClienteJuego[k][j] = "X"
			tableroClienteReal[k][j] = "X"
		}

		tableroActualizado := TableroActualizado{
			Letra:   letraAtaqueServidor,
			Tablero: tableroClienteReal,
		}

		tableroActualizadoBytes, _ := json.Marshal(tableroActualizado)

		conn.Write(tableroActualizadoBytes)
	}
}

func main() {
	fmt.Println("Servidor iniciado...")
	tableroCliente := generarTablero(1)
	tableroServidor := generarTablero(2)

	// Se crea el servidor UDP
	servidorUDP(tableroCliente)

	// Se crea el servidor TCP
	iniciarTCP(tableroCliente, tableroServidor)
}