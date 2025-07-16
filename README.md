# Go Peer-to-Peer File Synchronization

> Lien du repo : https://github.com/Axel77g/go-peer-to-peer

### Vue d'ensemble

Cette application est un **système de synchronisation de fichiers de pair à pair** qui permet à plusieurs pairs de synchroniser le contenu d'un répertoire partagé. Sa conception est similaire à celle de Git, où les changements de fichiers sont suivis comme des événements dans un historique, ce qui permet une synchronisation efficace au sein d'un réseau distribué.

# Lancement rapide
### Prérequis

- Deux ordinateurs avec **Go** installé (version 1.18 ou supérieure recommandée).
- Assurez-vous que les pare-feux autorisent les connexions **UDP** et **TCP** sur le réseau local.
  - **Attention :** Sur Windows, vérifiez que le pare-feu ne bloque pas l’application. Sur macOS, assurez-vous que le coupe-feu est désactivé ou que l’application est autorisée.
  - Le réseau local ne doit pas être bloqué par des politiques personnalisées (exemple : *Campus Erard*).

### Lancement de l’application

Sur chaque ordinateur, ouvrez un terminal dans le dossier du projet et lancez :

```bash
go run main.go
```

Les deux ordinateurs devraient se détecter automatiquement.  
Vous pouvez ensuite modifier, créer ou supprimer des fichiers dans le dossier [`./shared`](./shared/) (généré automatiquement au lancement).  
La synchronisation des événements sera visible dans [`./events.jsonl`](./events.jsonl).




-----

### Concepts clés

#### Synchronisation de fichiers basée sur les événements

Le système suit les modifications de fichiers (création, modification, suppression) comme des événements qui sont stockés dans un historique. Cette approche offre une piste d'audit de toutes les modifications, permet la résolution des conflits lors de la fusion des changements de plusieurs pairs et rend possible la reconstruction de l'état du répertoire à n'importe quel moment.

#### Architecture réseau

L'application utilise une double approche de protocole pour la communication entre les pairs :

  * **UDP** pour la découverte des pairs sur le réseau local.
  * **TCP** pour le transfert fiable des événements de fichiers et la synchronisation.

-----

### Composants clés

#### Système d'événements de fichiers

Les événements sont stockés dans un fichier au format **JSONL (JSON Lines)**.

  * **Collection d'événements :** Stocke les événements dans un journal en mode "append-only".
  * **Types d'événements :** Créer, Mettre à jour et Supprimer pour les fichiers du répertoire partagé.
  * **Itérateur d'événements :** Permet une traversée efficace de l'historique des événements.



#### Couche de communication réseau

Cette couche gère toute la communication de pair à pair.

  * **Abstraction de transport :** Une interface unifiée pour tous les canaux de communication.
  * **Gestion des pairs :** Permet de suivre les pairs connectés et leurs canaux de communication disponibles.
  * **Diffusion d'événements :** Distribue les événements de fichiers à tous les pairs connectés.

#### Observateur de fichiers

Ce composant surveille le répertoire partagé pour détecter les modifications et génère les événements de fichiers correspondants (création, modification ou suppression).

-----

### Architecture

#### Flux des événements

1.  Un observateur de système de fichiers détecte un changement dans le répertoire partagé.
2.  Un événement est généré et ajouté à la collection d'événements locale.
3.  L'événement est ensuite diffusé à tous les pairs connectés.
4.  Les pairs récepteurs fusionnent les événements entrants avec leur collection d'événements locale.
5.  Les changements sont propagés à travers le réseau.

#### Découverte et communication réseau

  * Des **messages de diffusion UDP** sont utilisés pour découvrir d'autres pairs sur le réseau local.
  * Des **connexions TCP** sont établies pour un transfert de données fiable.
  * Les pairs échangent des collections d'événements pour synchroniser l'état des fichiers.

#### Modèle de gestionnaire

Les gestionnaires implémentent la logique métier pour différents protocoles de communication :

  * **Gestionnaire de découverte UDP :** Gère la découverte des pairs et les connexions initiales.
  * **Gestionnaire de contrôleur TCP :** Gère la synchronisation des événements et le transfert des données de fichiers.

-----

### Détails d'implémentation

#### Collection d'événements de fichiers

Les événements sont stockés dans un **fichier JSONL**, ce qui permet des opérations d'ajout simples, une sérialisation/désérialisation facile et un stockage et un transfert efficaces.

#### Gestion des pairs

  * Chaque pair est identifié par son adresse IP.
  * Un seul pair peut avoir plusieurs canaux de transport.
  * Les pairs sont gérés dans un registre thread-safe.


### Autre informations

- L'ajout de comprésion peut être fait rapidement crace a notre compresser.go et l'utilisation de TransportChannel
- L'ajout de l'encryption peut être ajouter de la même manière
- Test unitaires sur le Compresseur / Décompresseur, L'encryption et la collection JSONL events


-----

### Développements futurs

  * **Reconstruction de répertoire** à partir des événements.
  * **Transfert de contenu de fichiers** pour une synchronisation effective.
  * **Résolution des conflits** pour les modifications simultanées.
  * **Visualisation du répertoire à distance.**

-----

### Philosophie de conception

L'application est conçue selon plusieurs principes clés :

  * **Abstraction :** Les couches réseau et de stockage sont abstraites derrière des interfaces.
  * **Modularité :** Les composants sont découplés et peuvent être remplacés indépendamment (exemple l'event collection peut ne pas dépendre du JSONL, il peut même y avoir une cohabitation).
  * **Indépendance du protocole :** La logique métier est séparée des mécanismes de transport.
  * **Event-Sourcing :** Tous les changements sont suivis comme des événements immuables dans un journal en mode "append-only".

Cette architecture offre de la flexibilité pour de futures extensions et modifications tout en maintenant une séparation claire des préoccupations.
