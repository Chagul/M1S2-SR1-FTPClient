# Client FTP
**Auteur** : Aurélien Plancke M1 E-Services

Ce programme tree-FTP à été réalisé en GO dans le cadre des cours de Système Répartis en Master 1 E-Services.
Il a pour but de se comporter en tant que client FTP et se connecter à un serveur, afin d'en afficher le contenu, sous 
une forme semblable au résultat que produit la commande linux tree.

## Installation

### Prerequis
Afin d'executer ce programme, il faut tout d'abord installer [le langage Go](https://go.dev/doc/install)

Une fois cette étape effectuée, il suffit de se rendre à la racine de ce dépôt et d'exécuter la commande 
``
go build
``

L'executable sera ainsi créé sous le nom de tree-ftp. Afin de l'utiliser il suffit de tapper ``./tree-ftp``

## Fonctionnement de la CLI
Ce programme étant en Command LIne, il y a de nombreux paramètres pouvant être précisés.

Tout d'abord les paramètres de bases obligatoires sont --addressServer et --port, afin de les utiliser, il suffit de taper 

``./tree-ftp --addressServer="l'adresse_du_serveur" --port="port_du_serveur"``

Le port par défaut est le 21 pour le
protocole FTP, mais la possibilité est laissée à l'utilisateur de le préciser. Par défaut, si le port n'est pas précisé,
la valeur 21 sera donc prise.

Dans le cas d'une connexion qui requiert une authentification, il est possible de préciser les flags ``--user`` et ``
--password`.
Si l'un des deux est précisé, l'autre est obligatoire. Sans les préciser, les valeurs "anonymous" seront utilisées pour
l'user et le mot de passe.

L'option ``--maxDepth`` peut être précisé afin d'arrêter le listing à une profondeur de l'arborescence du serveur. Par défaut ou 
si celle-ci n'est pas précisée, tree-ftp parcoureras toute l'arborescence.

## Architecture

## UML


## Vidéo de fonctionnement