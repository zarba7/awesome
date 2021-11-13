#!/bin/bash

##docker
curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun

##etcd
port=2379
name=sdw_etcd
port2=$(expr $port + 1)

data_dir=/data/"$name"/data

##etcd
set -eux
mkdir -p "$data_dir"
docker run -d --restart=always -p "$port":2379 -p "$port2":2380 \
	--mount type=bind,source="$data_dir",destination=/etcd-data \
	--name "$name" \
	quay.io/coreos/etcd:v3.4.9 \
	/usr/local/bin/etcd \
	--name s1 \
	--data-dir /etcd-data \
	--listen-client-urls http://0.0.0.0:2379 \
	--advertise-client-urls http://0.0.0.0:2379 \
	--listen-peer-urls http://0.0.0.0:2380 \
	--initial-advertise-peer-urls http://0.0.0.0:2380 \
	--initial-cluster s1=http://0.0.0.0:2380 \
	--initial-cluster-token tkn \
	--initial-cluster-state new \
	--log-level info \
	--logger zap \
	--log-outputs stderr \
	--auto-compaction-retention 1 --max-request-bytes 8388608 --quota-backend-bytes 8589934592


docker exec etcd-gcr-v3.4.9 /bin/sh -c "/usr/local/bin/etcd --version"
docker exec etcd-gcr-v3.4.9 /bin/sh -c "/usr/local/bin/etcdctl version"
docker exec etcd-gcr-v3.4.9 /bin/sh -c "/usr/local/bin/etcdctl endpoint health"