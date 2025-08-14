#!/bin/bash
set -e

CB_USER=${COUCHBASE_ADMINISTRATOR_USERNAME:-admin}
CB_PASS=${COUCHBASE_ADMINISTRATOR_PASSWORD:-password}
CB_HOST=host.docker.internal
BUCKET_NAME=my_bucket

echo "Cluster başlatılıyor..."

couchbase-cli cluster-init -c 127.0.0.1 \
  --cluster-username $CB_USER \
  --cluster-password $CB_PASS \
  --services data,index,query \
  --cluster-ramsize 512 \
  --cluster-index-ramsize 256 \
  --index-storage-setting default \
  --cluster-name my-cluster

echo "Dış hostname ayarlanıyor..."

# Couchbase ayar dosyasını düzenle (örnek)
sed -i "s/127.0.0.1/$CB_HOST/g" /opt/couchbase/etc/couchbase/static_config || true

# Alternatif olarak, external hostname'i API ile ayarlayabilirsin:
curl -u $CB_USER:$CB_PASS -X POST "http://127.0.0.1:8091/node/controller/rename" -d "hostname=$CB_HOST"

echo "Bucket oluşturuluyor..."

couchbase-cli bucket-create -c 127.0.0.1 \
  --username $CB_USER \
  --password $CB_PASS \
  --bucket=$BUCKET_NAME \
  --bucket-type couchbase \
  --bucket-ramsize 256 \
  --wait

echo "Cluster hazır!"
