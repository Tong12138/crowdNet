# 该脚本文件我们命名为installChainCode.sh
# ~/Demo/installChainCode.sh
# 引用 utils.sh, createChannel.sh 脚本, 该文件也被创建在scripts目录中
. utils.sh
. createChannel.sh
# 安装链码
installChaincode() 
{
  PEER=$1
  ORG=$2
  setGlobals $PEER $ORG
  VERSION=${3:-1.0}
  set -x
  peer chaincode install -n mycc -v ${VERSION} -l ${LANGUAGE} -p ${CC_SRC_PATH} >&log.txt
  res=$?
  set +x
  cat log.txt
  verifyResult $res "Chaincode installation on peer${PEER}.org${ORG} has failed"
  echo "===================== Chaincode is installed on peer${PEER}.org${ORG} ===================== "
  echo
}

# 链码初始化
instantiateChaincode() 
{
  PEER=$1
  ORG=$2
  setGlobals $PEER $ORG
  VERSION=${3:-1.0}

  # while 'peer chaincode' command can get the orderer endpoint from the peer
  # (if join was successful), let's supply it directly as we know it using
  # the "-o" option
  if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
    set -x
    peer chaincode instantiate -o orderer.crowd.com:7050 -C $CHANNEL_NAME -n mycc -l ${LANGUAGE} -v ${VERSION} -c '{"Args":["init","a","100","b","200"]}' -P "AND ('OrgReqMSP.peer','OrgWorkMSP.peer')" >&log.txt
    res=$?
    set +x
  else
    set -x
    peer chaincode instantiate -o orderer.crowd.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n mycc -l ${LANGUAGE} -v 1.0 -c '{"Args":["init","a","100","b","200"]}' -P "AND ('OrgReqMSP.peer','OrgWorkMSP.peer')" >&log.txt
    res=$?
    set +x
  fi
  cat log.txt
  verifyResult $res "Chaincode instantiation on peer${PEER}.org${ORG} on channel '$CHANNEL_NAME' failed"
  echo " Chaincode is instantiated on peer${PEER}.org${ORG} on channel '$CHANNEL_NAME'"
  echo
}

# 函数调用
# go组织的peer0, peer1安装链码
echo "Install chaincode on peer0.orgreq..."
#nstallChaincode 0 1
echo "Install chaincode on peer0.orgreq..."
installChaincode 0 2
# cpp组织的peer0, peer1安装链码
# go组织的peer0进行初始化
echo "Instantiating chaincode on peer0.org2..."
instantiateChaincode 0 2