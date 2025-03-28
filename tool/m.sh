#!/bin/bash

rootfs="/home/yzr/Documents/mydocker/data/overlay2/a1a932b41bc2c85d855ec739341e30160056c9e95359f6ef012a8dafbbaa2882/merged"
echo "正在创建命名空间并启动容器..."
unshare --fork --pid --mount --user --map-root-user --uts --ipc --net --mount-proc \
  bash -c "
     echo '用户映射:'
     cat /proc/self/uid_map
     cat /proc/self/gid_map
     hostname mycontainer
     id
     echo '正在进入 chroot 环境...'
     chroot \"$rootfs\" /bin/bash -c \"echo '容器内 hostname: \$(hostname)'; exec /bin/bash\"
  "

