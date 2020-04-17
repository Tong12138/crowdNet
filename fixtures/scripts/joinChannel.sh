# 该脚本文件我们命名为joinChannel.sh
# ~/Demo/joinChannel.sh
# 定义一些变量, 下边的代码中会使用
# ====================================================
#						加入通道
# ====================================================
# 引用 utils.sh, createChannel.sh 脚本, 该文件也被创建在scripts目录中
. utils.sh
. createChannel.sh
joinChannel () {
	for org in 1 2; do
	    for peer in 0 1; do
		joinChannelWithRetry $peer $org
		echo "======== peer${peer}.org${org} joined channel '$CHANNEL_NAME' ======== "
		sleep $DELAY
		echo
	    done
	done
}
joinChannelWithRetry() {
  PEER=$1
  ORG=$2
  setGlobals $PEER $ORG

  set -x
  peer channel join -b $CHANNEL_NAME.block >&log.txt
  res=$?
  set +x
  cat log.txt
  if [ $res -ne 0 -a $COUNTER -lt $MAX_RETRY ]; then
    COUNTER=$(expr $COUNTER + 1)
    echo "peer${PEER}.org${ORG} failed to join the channel, Retry after $DELAY seconds"
    sleep $DELAY
    joinChannelWithRetry $PEER $ORG
  else
    COUNTER=1
  fi
  verifyResult $res "After $MAX_RETRY attempts, peer${PEER}.org${ORG} has failed to join channel '$CHANNEL_NAME' "
}
# 函数调用
joinChannel