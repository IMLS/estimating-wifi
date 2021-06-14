# imls-ppa

    curl -s --compressed "https://jadudm.github.io/imls-ppa/KEY.gpg" | sudo apt-key add -
    sudo curl -s --compressed -o /etc/apt/sources.list.d/imls-ppa.list "https://jadudm.github.io/imls-ppa/contents.list"
    sudo apt update
    sudo apt install session-counter
