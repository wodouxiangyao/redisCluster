#!bin/sh
redis-server redis.conf --cluster-announce-port $redisPort --cluster-announce-bus-port 1$redisPort  --cluster-announce-ip $redisHost
