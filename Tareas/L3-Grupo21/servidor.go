package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type RegistroDNS struct {
	ip   string
	TTL  string
	tipo string
}

func agregarRegistro(registrodns map[string]RegistroDNS, nombre_dominio string, ip string, ttl string, tipo string) {
	registro := RegistroDNS{ip, ttl, tipo}
	registrodns[nombre_dominio] = registro
}

func obtenerRegistro(registrodns map[string]RegistroDNS, nombre_dominio string) (RegistroDNS, bool) {
	registro, ok := registrodns[nombre_dominio]
	return registro, ok
}

func main() {
	// Definir la dirección del servidor
	registrodns := make(map[string]RegistroDNS)

	registrodns["example.com"] = RegistroDNS{"93.184.216.34", "3600", "A"}

	serverAddr := "localhost:63420"

	// Crear un servidor UDP
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// Escuchar en la dirección definida
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Buffer para leer los datos entrantes
	buffer := make([]byte, 1024)

	for {
		// Leer los datos entrantes
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		// Imprimir los datos recibidos
		partes := strings.Split(string(buffer[:n]), ",")

		opcion := partes[0]

		if opcion == "3" { // Salir
			fmt.Println("Saliendo del servidor")
			break
		}

		if opcion == "2" { // Consulta de IP
			nombre_dominio := partes[1]
			record, ok := obtenerRegistro(registrodns, nombre_dominio)
			response := []byte("Error")
			if ok {
				fmt.Printf("IP: %s\n", record.ip)
				response = []byte("IP Encontrada " + record.ip)
			} else {
				fmt.Printf("Nombre del dominio: %s no encontrado\n", nombre_dominio)
				response = []byte("Nombre del dominio no encontrado")
			}
			// Enviar el mensaje de respuesta
			_, err = conn.WriteToUDP(response, addr)
			if err != nil {
				fmt.Println("Error: ", err)
				continue
			}
		}

		if opcion == "1" { // Agregar registro
			nombre_dominio := partes[1]
			ip := partes[2]
			ttl := partes[3]
			tipo := partes[4]
			agregarRegistro(registrodns, nombre_dominio, ip, ttl, tipo)
			response := []byte("Registro agregado")
			_, err = conn.WriteToUDP(response, addr)
			if err != nil {
				fmt.Println("Error: ", err)
				continue
			}
			fmt.Printf("Registro agregado: Nombre Dominio: %s, IP: %s, TTL: %s, Tipo: %s\n", nombre_dominio, ip, ttl, tipo)
		}

	}
}
