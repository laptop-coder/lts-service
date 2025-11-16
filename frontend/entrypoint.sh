#!/bin/sh

# set the default value of the FRONTEND_PORT if empty
if [ "$FRONTEND_PORT" = '' ]; then
    export FRONTEND_PORT=31022
fi

sed -i '/^#/d' nginx.conf.template # remove comments
sed -i 's/\$\$\$/###/g' nginx.conf.template # replace $$$ by ### (screening)
envsubst < nginx.conf.template > nginx.conf # substitute env variables
sed -i 's/###/$/g' nginx.conf # replace ### by $
cp nginx.conf /etc/nginx/conf.d/default.conf

nginx -g "daemon off;"
