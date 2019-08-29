#!/bin/bash

config_file="$1"
kernels_dir="$2"
last_local_version="$3"

if [ ! "$config_file" ]; then
  echo "No path to config file"
  exit 3
fi

if [ ! "$kernels_dir" ]; then
  echo "No path to kernels dir"
  exit 4
fi

cd "$kernels_dir"
ver="`basename $(pwd)`"
tar xvf linux-$ver.tar.xz
cd linux-$ver
cp $config_file . -v
make -j6
