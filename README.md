# Net-Cat

Netcat est un utilitaire de réseau puissant qui peut être utilisé pour écouter et se connecter à des ports réseau, envoyer et recevoir des données via TCP et UDP, et même créer des tunnels réseau. Il est souvent utilisé pour les diagnostics de réseau, les tests de connectivité, la redirection de ports et l'automatisation de tâches réseau.

## Installation 

Netcat est généralement déjà installé sur les systèmes Linux et macOS, mais vous pouvez le télécharger et le compiler manuellement si nécessaire. Si vous utilisez Windows, vous pouvez utiliser une version de portage de netcat telle que ncat ou socat.

### Utilisation

# 1. Lancer le port / ou si vous souhaitez vous pouvez lui donner un numéro de port
```
go run .

go run . <numéro du port>

```

# 2. Ouvrir une nouvelle page dans votre terminal pour lancer le serveur
```
nc localhost <numéro du port>

```

# 3. Chater avec d'autre personne (une fois le serveur lancer sur votre nouvelle page vous devez lancer le serveur avec votre adressIP)

```
nc <adressIP> <numéro du port>

```
Cette commande est à transmettre aux autres personne avec qui vous voulez parler