#!/bin/bash

# 检查是否输入了镜像名称和输出模式

if [ "$#" -ne 2 ]; then
    echo "使用方法: $0 <镜像名称> <输出模式>"
    echo "输出模式可以是 'class' 或 'layer'"
    exit 1
fi

# 获取镜像名称和输出模式

image_name=$1
output_mode=$2

# 执行 docker inspect 命令

output=$(docker inspect $image_name)

# 使用 jq 工具解析 JSON 输出，并获取 RootFS.Layers 数组中的所有值

layers=$(echo $output | jq -r '.[0].RootFS.Layers[]')

# 初始化 chainID 为第一层的 layerID

chainID=$(echo $layers | cut -d' ' -f1)

if [ "$output_mode" = "class" ]; then
    # 输出所有的 layerID
    echo "layerIDs:"
    index=1
    for layerID in $layers; do
        echo "   layer ${index}: $layerID"
        index=$((index+1))
    done

    # 输出所有的 chainID
    echo "chainIDs:"
    chainID=$(echo $layers | cut -d' ' -f1)
    echo -n "   layer 1: "
    echo $chainID
    
    index=2
    for layerID in $(echo $layers | cut -d' ' -f2-); do
        chainID=$(echo -n "$chainID $layerID" | sha256sum | awk '{print $1}')
        chainID="sha256:$chainID"
        echo -n "   layer ${index}: "
        echo $chainID
        index=$((index+1))
    done
    
    # 输出所有的 cache-id 文件的内容
    echo "cacheIDs:"
    chainID=$(echo $layers | cut -d' ' -f1)
    cache_id=$(cat /var/lib/docker/image/overlay2/layerdb/sha256/${chainID:7}/cache-id)
    echo "   layer 1: $cache_id"
    
    index=2
    for layerID in $(echo $layers | cut -d' ' -f2-); do
        chainID=$(echo -n "$chainID $layerID" | sha256sum | awk '{print $1}')
        chainID="sha256:$chainID"
        cache_id=$(cat /var/lib/docker/image/overlay2/layerdb/sha256/${chainID:7}/cache-id)
        echo "   layer ${index}: $cache_id"
        index=$((index+1))
    done

elif [ "$output_mode" = "layer" ]; then
    index=1
    for layerID in $layers; do
        echo "layer ${index}:"
        echo "   layerID: $layerID"
        if [ "$index" -eq 1 ]; then
            chainID=$layerID
        else
            chainID=$(echo -n "$chainID $layerID" | sha256sum | awk '{print $1}')
            chainID="sha256:$chainID"
        fi
        echo "   chainID: $chainID"
        cache_id=$(cat /var/lib/docker/image/overlay2/layerdb/sha256/${chainID:7}/cache-id)
        echo "   cacheID: $cache_id"
        index=$((index+1))
    done
else
    echo "无效的输出模式: $output_mode"
    echo "输出模式可以是 'class' 或 'layer'"
    exit 1
fi