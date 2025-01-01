#!/bin/sh
# Script to setup this project for development.
# Run this script from the project's root directory (where this script is located)
cd ./frontend
npm i
cd ../backend
mkdir -p ./storage/{lost,found}
cat ./db.sql | sqlite3 ./db.sqlite3
python3 -m venv env
. ./env/bin/activate
pip install -r ./requirements.txt
clear
cd ../

