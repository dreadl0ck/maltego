#!/bin/bash

cd docker/ubuntu || exit 1
echo "[INFO] creating ${UBUNTU_IMAGE} container..."

echo "[INFO] building & tagging container"
docker build -t ${UBUNTU_IMAGE} .
if (( $? != 0 )); then
	echo "[ERROR] failed to build container."
	exit 1
fi

echo "[INFO] pushing ${UBUNTU_IMAGE} to docker hub registry"
docker push ${UBUNTU_IMAGE}
if (( $? != 0 )); then
	echo "[ERROR] failed to push container image. maybe you are not logged in?"
	exit 1
fi

echo "[INFO] done."