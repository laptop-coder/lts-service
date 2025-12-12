#!/bin/sh

BASE_DIR="$HOME/lts-service"
DATA_DIR="data"
PATH_TO_BACKUPS="$HOME/backups"
tar -C "$BASE_DIR" -cf "$PATH_TO_BACKUPS/backup_lts_$(date '+%Y-%m-%d_%H-%M-%S').tar" "$DATA_DIR"

