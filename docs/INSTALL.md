----------
amon 
----------

 
### Installation 
Copier le ficheir amon_[version].tar.gz sur la machine cible dans le repertoire /opt
Exécuter tar xvzf revor_[version].tar.gz

Si le lien symbolique  /opt/amon existe supprimer le et créer un nouveau lien symbolique

ln -s /opt/amon_[version] /opt/amon

Verifier la configuraiton de la version precedente.
Fichiers de configuraitons sont dans le répértoire /opt/amon/config

Gestion de cluster est prise en compte par crm_mon via la ligne de commande : crm_mon  --daemonize --as-xml /opt/amon/tmp/crm_mon.xml.
Cette commande doit/peut être ajouter dans le ficheir /etc/rc.local pour être prise en comptre au démarrage de l'application.

### Configuraiton
Fichier de la configuraiton config.json contient les patramètres suivantes :
	ListenAddr  l'adresse et le port d'écoute de serveur web
	DbMySqlHost      l'adresse et le port de serveur http
	DbMySqlUser      le nom d'utilisateur mysql 
	DbMySqlPassword  le mot de passe d'utilisateur mysql
	DbMySqlName      le nom de la base de données
	MongoHost        l'adresse de serveur de base de données mongodb
	AsteriskAddr     l'adresse de serveur AMI d'asterisk
	AsteriskPort     le port de serveur AMI d'asterisk
	AsteriskUser     l'utilsateur AMI 
	AsteriskPassword le mot de passe d'utilisateur AMI
	CrmMonFile       le chemin de fichier des information de cluster


### Mise à jour 
Mise à jour est identique à l'installation. Il faut verififer la configuraiton de la version en explotation.





