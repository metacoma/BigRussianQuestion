#!/bin/sh
cd /image_generator
IMAGE="`find /images -type f | sort -R | head -n1`"
cat "$IMAGE" | /image_generator/image_generator.sh "$*"
