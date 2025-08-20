FROM bitnami/kafka:3.5

WORKDIR /producer

COPY msgs/ ./msgs/
COPY send_msgs.sh .

CMD ["bash", "./send_msgs.sh"]

