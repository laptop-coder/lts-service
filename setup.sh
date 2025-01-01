#!/bin/sh
mkdir ./backend/storage
cat ./backend/db.sql | sqlite3 ./backend/db.sqlite3
cd frontend
npm i
alacritty -e npm run dev &
cd ../backend
python3 -m venv env
. ./env/bin/activate
pip install -r ./requirements.txt
fastapi dev main.py

