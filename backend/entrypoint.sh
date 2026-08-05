#!/bin/sh

chown -R appuser:appuser /backend
exec su -c /usr/local/bin/app appuser
