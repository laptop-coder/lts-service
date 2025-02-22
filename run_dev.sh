#!/bin/sh
# Script to run this project in development mode.
# Run this script from the project's root directory (where this script is located).
cd ./frontend
npm run dev &
cd ../backend
. ./env/bin/activate
fastapi dev main.py

