# 该脚本文件我们命名为updateChannel.sh
# ~/Demo/updateChannel.sh
# 引用 utils.sh, createChannel.sh 脚本, 该文件也被创建在scripts目录中
. utils.sh
. createChannel.sh
updateAnchorPeers() {
  PEER=$1
  ORG=$2
  setGlobals $PEER $ORG
  if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
    set -x
    # 更新锚节点不使用tls加密, 得到的结果保存到log.txt中
    peer channel update -o orderer.crowd.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx >&log.txt
    res=$?
    set +x
  else
    set -x
    # 更新锚节点使用tls加密, 得到的结果保存到log.txt中
    peer channel update -o orderer.crowd.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
    res=$?
    set +x
  fi
  cat log.txt
  verifyResult $res "Anchor peer update failed"
  echo " Anchor peers updated for org '$CORE_PEER_LOCALMSPID' on channel '$CHANNEL_NAME' "
  sleep $DELAY	# 休眠
  echo
}
# 函数调用
echo "Updating anchor peers for orggo..."
updateAnchorPeers 0 1
echo "Updating anchor peers for orgcpp..."
updateAnchorPeers 0 2