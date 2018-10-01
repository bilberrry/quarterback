# Installing latest version
PLATFORM='linux'

if [[ `uname` == 'Linux' ]]; then
   PLATFORM='linux'
elif [[ `uname` == 'Darwin' ]]; then
   PLATFORM='darwin'
fi

curl -s https://api.github.com/repos/bilberrry/quarterback/releases/latest | grep "browser_download_url.*$PLATFORM*" | cut -d : -f 2,3 | tr -d \" | wget -O quarterback.tar.gz -qi -
tar -zxf quarterback.tar.gz && sudo mv quarterback /usr/local/bin/quarterback && rm quarterback.tar.gz
echo "Quarterback successfully installed"