FROM ubuntu:latest

RUN apt-get update && \
		apt-get install -y --no-install-recommends vsftpd db-util gcc pkg-config libwebkit2gtk-4.0-dev fuse && \
		apt-get clean

ENV FTP_USER admin
ENV FTP_PASS admin
ENV PASV_ADDRESS REQUIRED

COPY vsftpd.conf /etc/vsftpd/
COPY vsftpd_virtual /etc/pam.d/
COPY run.sh /usr/sbin/
COPY onedriver /usr/sbin/

RUN chmod +x /usr/sbin/run.sh /usr/sbin/onedriver && \
		mkdir -p /var/run/vsftpd/empty

VOLUME /var/log/vsftpd

EXPOSE 20 21

CMD ["/usr/sbin/run.sh"]
