SERVER=simplecloud
USER=god
APPNAME=multik

go build -o $APPNAME
upx $APPNAME

ssh $USER@$SERVER "mkdir -p /home/$USER/sites/$APPNAME"
rsync -avzhp --exclude sites/ --exclude log/ --exclude scripts/install.sh . $USER@$SERVER:/home/$USER/sites/$APPNAME
ssh $USER@$SERVER << EOF
  cd /home/$USER/sites/multik
  sudo mv scripts/multik.conf /etc/init
  sudo mv scripts/multik_caddy.conf /etc/caddy/sites
  sudo stop caddy
  sudo start caddy
  sudo stop multik
  sudo start multik
EOF