#!/bin/sh

if [ "$PRODUCTION" = "false" ]; then
    echo "Starting backend in the DEVELOPMENT mode..."
    /usr/local/bin/migrate
elif [ "$PRODUCTION" = "true" ]; then
    echo "Starting backend in the PRODUCTION mode..."
else
    echo "PRODUCTION env variable must be true or false"
    exit 1
fi

/usr/local/bin/app

