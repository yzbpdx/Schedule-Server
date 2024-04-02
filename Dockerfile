FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Shanghai

WORKDIR /root

RUN apt update && mkdir logs

COPY schedule /root
COPY run.sh /root
COPY script.sql /root
COPY html /root/html
COPY conf/config.yaml /root/conf/config.yaml

RUN chmod +x /root/schedule && chmod +x /root/run.sh

EXPOSE 9000
EXPOSE 3306
EXPOSE 6379

CMD sh /root/run.sh