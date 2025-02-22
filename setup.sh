#!/bin/sh
# Script to setup this project for development.
# Run this script from the project's root directory (where this script is located)

function info {
  printf "\033[1m$1\033[0m\n"
}

function result {
  if [ $? -eq 0 ]; then
    printf "\033[1;32mDone\033[0m\n"
  else
    printf "\033[1;31mError\033[0m\n"
  fi
}

cd ./frontend

info "Installing Npm requirements..."
npm i > /dev/null
result

cd ../backend

info "Creating storage..."
mkdir -p ./storage/{lost,found} > /dev/null
result

info "Creating database..."
cat ./db.sql | sqlite3 ./db.sqlite3 > /dev/null
result

info "Creating Python virtual environment..."
python3 -m venv env > /dev/null
result

info "Activating Python virtual environment..."
. ./env/bin/activate > /dev/null
result

info "Updating Python package manager..."
pip install --upgrade pip > /dev/null
result

info "Installing Python requirements..."
pip install -r ./requirements.txt > /dev/null
result

cd ../


