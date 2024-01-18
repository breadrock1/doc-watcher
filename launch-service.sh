docker build -t fs-notifier:latest .
docker volume create --driver local fsnotifiervol
docker run -d -v fsnotifiervol:/indexer fs-notifier:latest
