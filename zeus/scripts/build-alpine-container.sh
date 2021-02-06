#!/bin/bash

# copy golang trx server binary
rm -f docker/alpine/trx
cp bin/trx docker/alpine

cd docker/alpine || exit 1
echo "[INFO] creating ${ALPINE_IMAGE} container..."

echo "[INFO] building & tagging container"
docker build -t ${ALPINE_IMAGE} .
if (( $? != 0 )); then
	echo "[ERROR] failed to build container."
	exit 1
fi

echo "[INFO] pushing ${ALPINE_IMAGE} to docker hub registry"
docker push ${ALPINE_IMAGE}
if (( $? != 0 )); then
	echo "[ERROR] failed to push container image. maybe you are not logged in?"
	exit 1
fi

echo "[INFO] done."