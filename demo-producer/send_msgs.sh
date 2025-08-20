for f in ./msgs/*.json; do
  tr -d '\r' < "$f" | kafka-console-producer.sh --broker-list kafka:9092 --topic orders
  sleep 1
done
