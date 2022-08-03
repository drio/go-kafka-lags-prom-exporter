#!/bin/sh

bs="kafka-test-01.tuftscloud.com:9094,kafka-test-02.tuftscloud.com:9094,kafka-test-03.tuftscloud.com:9094"
config=client-ssl.properties.sasl_plain.aws

# GROUP  TOPIC    PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG     CONSUMER-ID     HOST   CLIENT-ID
kafka-consumer-groups \
	--bootstrap-server $bs \
	--command-config $config  \
	--describe \
	--all-groups | \
	grep -v GROUP | \
  grep -v "^$" | \
	# GROUP  TOPIC    PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG 
	awk '{print $1","$2","$3","$4","$5","$6}'
