package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	clients     = make(map[net.Conn]string)
	messages    []string
	clientCount int
)

// Mutex pour protéger l'accès concurrent à la liste de connections
var clientsMutex sync.Mutex

func main() {
	port := "8989"
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()
	fmt.Println("Ecouteur TCP en attente de connections sur le port", port)
	// Boucle infinie pour accepter les connections entrantes
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Vérifie si le nombre de clients connectés atteint déjà 10
		if clientCount >= 10 {
			fmt.Fprint(conn, "Le serveur est plein, veuillez réessayer plus tard.\n")
			conn.Close()
			continue
		}
		clientCount++
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	// Fermeture de la connection lorsque la fonction se termine
	defer func() {
		// Suppression de la connection client de la liste
		clientsMutex.Lock()
		clientCount--
		delete(clients, conn)
		clientsMutex.Unlock()
		conn.Close()
	}()
	// Création d'un lecteur pour lire les données envoyées par le client
	reader := bufio.NewReader(conn)
	// Ouvre le fichier texte à envoyer
	file, err := os.Open("welcome.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// Envoie le contenu du fichier au client
	_, err = io.Copy(conn, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Demande au client de choisir un pseudo
	var pseudo string
	for {
		fmt.Fprint(conn, "Entrez votre pseudo :")
		pseudo, _ = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		pseudo = strings.TrimSpace(pseudo)
		// Vérifie si le pseudo est vide
		if pseudo != "" {
			break
		}
	}
	// Ajout de la connection client et de son pseudo à la liste
	clientsMutex.Lock()
	clients[conn] = pseudo
	clientsMutex.Unlock()
	// Envoi d'un message de bienvenue à tous les clients
	sendMessageToOtherClients(conn, pseudo+" a rejoint le chat\n")
	// Envoi tout les messages précédent la connection du client
	for _, message := range messages {
		fmt.Fprint(conn, message)
	}
	messages = append(messages, pseudo+" a rejoint le chat\n")
	for {
		message := ""
		now := time.Now()
		message, err = reader.ReadString('\n')
		if message != "\n" {
			message = ("[" + now.Format("2006-01-02 15:04:05") + "][" + pseudo + "]:" + message)
			if err != nil {
				if errors.Is(err, io.EOF) {
					messages = append(messages, "Connexion fermée par "+pseudo+"\n")
					sendMessageToAllClients("Connexion fermée par " + pseudo + "\n")
				} else {
					fmt.Println(err)
				}
				return
			}
			fmt.Fprintf(conn, "\033[1A")
			fmt.Fprintf(conn, "\033[2K")
			sendMessageToAllClients(message)
			if err != nil {
				fmt.Println(err)
				return
			}
			messages = append(messages, message)
		} else {
			fmt.Fprintf(conn, "["+now.Format("2006-01-02 15:04:05")+"]["+pseudo+"]:")
		}
	}
}
func sendMessageToAllClients(message string) {
	// Envoi du message à toutes les connections client
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for conn := range clients {
		_, err := fmt.Fprintf(conn, "%s", message)
		if err != nil {
			fmt.Println(err)
		}
	}
}
func sendMessageToOtherClients(conn net.Conn, message string) {
	// Envoi du message aux autres connections client
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for c := range clients {
		// Vérifie si la connection courante est la connection du nouveau client
		if c == conn {
			continue
		}
		_, err := fmt.Fprintf(c, "%s", message)
		if err != nil {
			fmt.Println(err)
		}
	}
}
