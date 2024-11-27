# Importo la libreria socket para manejar conexiones
import socket as skt

# Se declara dirección del servidor y puerto
serverAddr = 'localhost'
serverPort = 63420

# Se crea un socket para hacer el manejo de la conexión
clientSocket = skt.socket(skt.AF_INET, skt.SOCK_DGRAM)

while True:
    # Solicitar al usuario que elija una opción
    print("1. Ingresar un registro")
    print("2. Consultar un registro")
    print("3. Detener el programa")
    option = input("Elige una opción: ")

    if option == "1":
        # Solicitar al usuario que ingrese los detalles del registro
        nombre_dominio = input("Ingresa dominio: ")
        ip = input("Ingresa IP: ")
        ttl = input("Ingresa TTL: ")
        tipo = input("Ingresa tipo: ")

        # Crear el mensaje a enviar
        toSend = f"{option},{nombre_dominio},{ip},{ttl},{tipo}"
        clientSocket.sendto(toSend.encode(), (serverAddr,serverPort))

    elif option == "2":
        # Solicitar al usuario que ingrese el nombre del dominio
        nombre_dominio = input("Ingresa el nombre del dominio para consultar: ")

        # Crear el mensaje a enviar
        toSend = f"{option},{nombre_dominio}"
        clientSocket.sendto(toSend.encode(), (serverAddr,serverPort))

    elif option == "3" or toSend == "STOP":
        # Detener el programa
        toSend = f"{option},STOP"
        clientSocket.sendto(toSend.encode(), (serverAddr,serverPort))
        print("Programa Finalizado")
        break

    else:
        print("Opción no válida. Por favor, elige una opción válida.")

    # Recibir y imprimir la respuesta del servidor
    msg, addr = clientSocket.recvfrom(1024)
    print(msg.decode()) 