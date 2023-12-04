#!/bin/sh

if [ ! -f "/var/www/html/packages.json" ]; then
	satis build satis.json /var/www/html
fi

satis-hook
