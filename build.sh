#!/bin/bash

xgo --targets=windows/*,darwin/amd64,linux/amd64,linux/386,linux/arm --dest=cache ./

osarchs=(windows_amd64 windows_386 darwin_amd64 linux_amd64 linux_386 linux_arm)
files=(samaritan-windows-4.0-amd64.exe samaritan-windows-4.0-386.exe samaritan-darwin-10.6-amd64 samaritan-linux-amd64 samaritan-linux-386 samaritan-linux-arm-5)
saves=(samaritan_windows_amd64.exe samaritan_windows_386.exe samaritan_darwin_amd64 samaritan_linux_amd64 samaritan_linux_386 samaritan_linux_arm)

unzip web/dist.zip -d web

for i in 0 1 2 3 4 5; do
  mkdir cache/samaritan_${osarchs[${i}]}
  mkdir cache/samaritan_${osarchs[${i}]}/web
  mkdir cache/samaritan_${osarchs[${i}]}/custom
  cp LICENSE cache/samaritan_${osarchs[${i}]}/LICENSE
  cp -r plugin cache/samaritan_${osarchs[${i}]}/plugin
  cp README.md cache/samaritan_${osarchs[${i}]}/README.md
  cp -r web/dist cache/samaritan_${osarchs[${i}]}/web/dist
  cp config.ini cache/samaritan_${osarchs[${i}]}/custom/config.ini
  mv cache/${files[${i}]} cache/samaritan_${osarchs[${i}]}/${saves[${i}]}
done

zip -r ./cache.zip ./cache/

rm -rf web/dist cache
