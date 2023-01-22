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

L'option ``--directoryOnly`` sert à afficher uniquement les repertoires lors du tree. Elle prends un booléen en paramètre,
par défaut celui-ci est à false.

L'option ``--fullPath`` sert à afficher les chemins complets des fichiers lors du tree. Elle prends un booléen en paramètre,
par défaut celui-ci est à false.

## Fonctionnement du programme
Le programme à une logique simple. Dans un premier temps nous allons établir une connexion avec l'url fournie lors du lancement
de la CLI. Le protocole FTP fonctionne de tel manière qu'il y a une connexion pour les commandes et leur retour, et une pour les données.
Cette seconde connexion et ses informations est établi grâce aux informations renvoyées par la commande ``PASV`` sur la 
connexion principale. Une fois cette seconde connexion établie, nous allons envoyé la commande ``LIST``sur le canal principal
et nous aurons le retour sur la connexion de données. Cette commande nous renvoyant une liste des fichiers à la position
actuelle, il ne nous reste plus qu'à changer de dossier courant pour chaque dossier trouvé avec la commande ``CWD``. Une
nouvelle connexion avec ``PASV`` est necessaire à chaque ``CWD``. 

## Architecture

La base du programme est dans ``cmd/root.go``, grâce a la bibliothèque [Cobra](https://github.com/spf13/cobra), c'est à cet
endroit que sont definies les options attendues, qu'elles sont parsées et que leur valeur recupérées. Grâce à cela nous 
pouvons établir la connexion TCP qui nous servira à communiquer avec le serveur FTP. Ensuite tout ce passe dans le package
``tcpconn``.
Ce package concerne toute les interactions serveur, ce qui permet de séparer les responsabilités.

L'architecture logiciel est plus compliquée en Go, ce langage n'étant pas un langage orienté objet, cependant j'ai fait l'utilisation des 
[pointer receivers](https://go.dev/tour/methods/4), notamment dans le cadre du type Node, qui représente un arbre de données
pour l'arborescence des fichiers. Les méthodes que j'ai ainsi pu établir sur ce type particulier lui étaient liées. Voici
le fichier model.go, avec la definition de la structure Node et ses fonctions liées. Cela nous permet de facilement ajouter 
des fonctionnalités dans le type node si necessaire, sans avoir à faire trop de refactoring.

Voici à quoi ressemble cette classe

![img.png](rsc/img.png)

![img.png](rsc/uml.png)

## UML


## Vidéo de fonctionnement