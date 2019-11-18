#!/bin/bash

config_file="$2"
kernels_dir="$1"
last_local_version="$3"
echo -e "Config file: $config_file\nKernels dir: $kernels_dir\nLast Local Version: $last_local_version"
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
tar xf linux-$ver.tar.xz
cd linux-$ver
cp $config_file ./.config -v
let threads=`cat /proc/cpuinfo | grep processor | wc -l`+1
make -j$threads
