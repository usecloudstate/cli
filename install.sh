#!/bin/bash
{
    set -e
    SUDO=''
    if [ "$(id -u)" != "0" ]; then
      SUDO='sudo'
      echo "This script requires superuser access."
      echo "You will be prompted for your password by sudo."
      # clear any previous sudo permission
      sudo -k
    fi


    # run inside sudo
    $SUDO bash <<SCRIPT
  set -e

  echoerr() { echo "\$@" 1>&2; }

  if [[ ! ":\$PATH:" == *":/usr/local/bin:"* ]]; then
    echoerr "Your path is missing /usr/local/bin, you need to add this to use this installer."
    exit 1
  fi

  if [ "\$(uname)" == "Darwin" ]; then
    OS=darwin
  elif [ "\$(expr substr \$(uname -s) 1 5)" == "Linux" ]; then
    OS=linux
  else
    echoerr "This installer is only supported on Linux and MacOS"
    exit 1
  fi

  ARCH="\$(uname -m)"
  if [ "\$ARCH" == "x86_64" ]; then
    ARCH=amd64
  elif [[ "\$ARCH" == aarch* ]]; then
    ARCH=arm
  else
    echoerr "unsupported arch: \$ARCH"
    exit 1
  fi

  mkdir -p /usr/local/lib
  cd /usr/local/lib
  rm -rf cloudstate
  rm -rf ~/.local/share/cloudstate/client

  mkdir -p cloudstate/bin
  cd cloudstate/bin

  URL="https://github.com/usecloudstate/cli/releases/download/v1.0.10/cloudstate-cli-\$OS-\$ARCH"

  echo "Installing CLI from \$URL"
  if [ \$(command -v curl) ]; then
    curl -L "\$URL" --output cli
  else
    wget -O- "\$URL" > cli
  fi
  
  # delete old cloudstate bin if exists
  rm -f \$(command -v cloudstate) || true
  rm -f /usr/local/bin/cloudstate
  ln -s /usr/local/lib/cloudstate/bin/cli /usr/local/bin/cloudstate

  chmod +x /usr/local/lib/cloudstate/bin/cli

SCRIPT
  # test the CLI
  LOCATION=$(command -v cloudstate)
  echo "cloudstate installed to $LOCATION"
}