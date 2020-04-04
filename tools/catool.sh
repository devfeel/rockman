!/bin/bash
# call this script with an email address (valid or not).
# like:
# ./catool.sh pzrr@qq.com
echo "make server cert"
openssl req -new -nodes -x509 -out server.pem -keyout server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=dotweb.cn/emailAddress=$1"
echo "make client cert"
openssl req -new -nodes -x509 -out client.pem -keyout client.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=dotweb.cn/emailAddress=$1"
