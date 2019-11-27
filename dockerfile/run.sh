#!/bin/bash

# Create home dir and update vsftpd user db:
echo -e "${FTP_USER}\n${FTP_PASS}" > /etc/vsftpd/virtual_users.txt
db_load -T -t hash -f /etc/vsftpd/virtual_users.txt /etc/vsftpd/virtual_users.db

# Set passive mode parameters:
if [ "$PASV_ADDRESS" = "REQUIRED" ]; then
	echo "Please insert IPv4 address of your host"
	exit 1
fi
echo "pasv_address=${PASV_ADDRESS}" >> /etc/vsftpd/vsftpd.conf
echo "guest_username=root" >> /etc/vsftpd/vsftpd.conf

# Run vsftpd:
sleep 10 && vsftpd /etc/vsftpd/vsftpd.conf &
mkdir -p /home/vsftpd
chown -R ftp:ftp /home/vsftpd
usermod -g root ftp
/usr/sbin/onedriver /home/vsftpd

