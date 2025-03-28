#!/bin/bash

# 定义变量
NS_NAME="demo_ns"
VETH_HOST="veth0"
VETH_NS="veth1"
HOST_IP="192.168.100.1"
NS_IP="192.168.100.2"
HOST_PORT=8080
NS_PORT=80

# 清理旧配置（防止冲突）
ip netns del $NS_NAME 2>/dev/null
ip link del $VETH_HOST 2>/dev/null
iptables -F
iptables -t nat -F

# 1. 创建网络命名空间
ip netns add $NS_NAME
echo "[1] 创建命名空间: $NS_NAME"

# 2. 创建veth设备对并配置IP
ip link add $VETH_HOST type veth peer name $VETH_NS
ip link set $VETH_NS netns $NS_NAME
ip addr add $HOST_IP/24 dev $VETH_HOST
ip link set $VETH_HOST up
echo "[2] 宿主机端veth配置完成: $VETH_HOST ($HOST_IP)"

# 命名空间内配置网络
ip netns exec $NS_NAME ip addr add $NS_IP/24 dev $VETH_NS
ip netns exec $NS_NAME ip link set $VETH_NS up
ip netns exec $NS_NAME ip link set lo up
ip netns exec $NS_NAME ip route add default via $HOST_IP
echo "[3] 命名空间内网络配置完成: $VETH_NS ($NS_IP)"

# 3. 启用IP转发和NAT规则
sysctl -w net.ipv4.ip_forward=1 >/dev/null
iptables -t nat -A POSTROUTING -s 192.168.100.0/24 -j MASQUERADE
echo "[4] 启用IP转发和SNAT规则"

# 4. 添加端口转发规则
# DNAT规则（外部访问宿主机端口）
iptables -t nat -A PREROUTING -i ens33 -p tcp --dport $HOST_PORT -j DNAT --to-destination $NS_IP:$NS_PORT
# DNAT规则（宿主机本地访问）
iptables -t nat -A OUTPUT -p tcp --dport $HOST_PORT -j DNAT --to-destination $NS_IP:$NS_PORT
# FORWARD规则（允许双向流量）
iptables -A FORWARD -i ens33 -o veth0 -p tcp --dport $NS_PORT -j ACCEPT
iptables -A FORWARD -i veth0 -o ens33 -m state --state ESTABLISHED,RELATED -j ACCEPT
echo "[5] 添加DNAT和FORWARD规则"

# 5. 在命名空间内启动HTTP服务
echo "[6] 在命名空间内启动HTTP服务..."
ip netns exec $NS_NAME python3 -m http.server $NS_PORT --bind 0.0.0.0 >/dev/null 2>&1 &
sleep 2

# 6. 验证服务监听状态
echo "[7] 检查命名空间内端口监听状态:"
ip netns exec $NS_NAME ss -tulnp | grep ":$NS_PORT"
if [ $? -ne 0 ]; then
    echo "错误: 命名空间内端口 $NS_PORT 未监听！"
    exit 1
fi

# 7. 测试端口转发效果
echo "[8] 测试端口转发："
echo "  - 在宿主机执行: curl http://127.0.0.1:$HOST_PORT"
curl -s http://127.0.0.1:$HOST_PORT | grep "Directory listing"
if [ $? -eq 0 ]; then
    echo "测试成功！宿主机 $HOST_PORT 端口已映射到命名空间内的 HTTP 服务。"
else
    echo "测试失败！请执行以下命令进一步排查："
    echo "1. 检查 iptables 规则: iptables-save"
    echo "2. 检查命名空间内服务日志: ip netns exec $NS_NAME ps aux | grep python3"
    echo "3. 抓包分析: tcpdump -i ens33 port $HOST_PORT -nnv"
fi

# 8. 清理（可选）
# kill %1
# ip netns del $NS_NAME
# ip link del $VETH_HOST
# iptables -F
# iptables -t nat -