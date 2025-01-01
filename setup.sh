#!/bin/sh
# Run this script from the project's root directory (where this script is located)
mkdir -p ./backend/storage/{lost,found}
cat ./backend/db.sql | sqlite3 ./backend/db.sqlite3
cd frontend
npm i
alacritty -e npm run dev &
cd ../backend
python3 -m venv env
. ./env/bin/activate
pip install -r ./requirements.txt
fastapi dev main.py

