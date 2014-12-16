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



### Mise à jour 
Mise à jour est identique à l'installation. Il faut verififer la configuraiton de la version en explotation.





