# 该脚本文件我们命名为createChannel.sh
# ~/Demo/createChannel.sh
# 定义一些变量, 下边的代码中会使用
CHANNEL_NAME="$1"	# 脚本文件执行时接收的第1个参数
DELAY="$2"			# 脚本文件执行时接收的第2个参数
LANGUAGE="$3"		# 脚本文件执行时接收的第3个参数
: ${CHANNEL_NAME:="mychannel"}	# 通道名，如果值为空, 变量CHANNEL_NAME值设置为mychannel
: ${DELAY:="3"}					# 延时时长，如果值为空, 变量DELAY值设置为3, 操作失败重试时使用
: ${LANGUAGE:="golang"}			# 链码语言，如果值为空, 变量LANGUAGE值设置为golang
LANGUAGE=`echo "$LANGUAGE" | tr [:upper:] [:lower:]` # 语言转换为大写，再转换为小写
COUNTER=1			# 计数器
MAX_RETRY=5			# 操作失败，重试的次数
# 链码文件路径
CC_SRC_PATH="github.com/chaincode/chaincode_example02/go/" # go链码
if [ "$LANGUAGE" = "node" ]; then						   # node.js链码
	CC_SRC_PATH="/opt/gopath/src/github.com/chaincode/chaincode_example02/node/"
fi
echo "Channel name : "$CHANNEL_NAME		# 打印通道的名字

# ====================================================
#						创建通道
# ====================================================
# 引用 utils.sh 脚本, 该文件也被创建在scripts目录中
. utils.sh
createChannel() {
	setGlobals 0 1
	if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
        set -x
        # 创建通道, 不使用tls加密, 结果重定向到文件log.txt中
		peer channel create -o orderer.crowd.com:7050 -c $CHANNEL_NAME -f ../channel-artifacts/channel.tx >&log.txt
		res=$?	# 保存 peer channel create 命令的返回值
        set +x
	else
		set -x
		# 创建通道, 使用tls加密, 结果重定向到文件log.txt中
		peer channel create -o orderer.crowd.com:7050 -c $CHANNEL_NAME -f ../channel-artifacts/channel.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
		res=$?	# 保存 peer channel create 命令的返回值
		set +x
	fi
	cat log.txt	# 将log文件内容输出到终端
	verifyResult $res "Channel creation failed"	# 验证操作是否成功
	echo "===================== Channel '$CHANNEL_NAME' created ===================== "
	echo
}
# 调用创建通道的函数
createChannel