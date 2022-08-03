FROM confluentinc/cp-kafka:7.0.1
ADD . /app
WORKDIR /app
CMD ["/bin/sh", "./run.sh"]
