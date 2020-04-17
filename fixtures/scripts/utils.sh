# 该脚本文件我们命名为 utils.sh
# ~/Demo/utils.sh
# 定义一些变量, 下边的代码中会使用 
# orderer 节点的证书
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/crowd.com/orderers/orderer.crowd.com/msp/tlscacerts/tlsca.crowd.com-cert.pem
# Go组织peer0节点证书
PEER0_ORGREQ_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgreq.crowd.com/peers/peer0.orgreq.crowd.com/tls/ca.crt
# CPP组织peer0节点证书
PEER0_ORGWORK_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgwork.crowd.com/peers/peer0.orgwork.crowd.com/tls/ca.crt

PEER0_ORGVALID_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgvalid.crowd.com/peers/peer0.orgvalid.crowd.com/tls/ca.crt

# 根据传递进行来的参数读取各个组织中各节点信息
# 第一次参数表示节点编号, 第二个参数表示组织编号/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/signcerts
setGlobals() {
  PEER=$1
  ORG=$2
  if [ $ORG -eq 1 ]; then
    CORE_PEER_LOCALMSPID="OrgReqMSP"		# Go组织的ID, 在configtx.yaml中设置的
    # 客户端cli想要连接哪个节点, 就必须要设置该节点的证书路径
    CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORGREQ_CA	# cli可以连接Go组织的peer0节点
    # Go组织的管理员证书
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgreq.crowd.com/users/Admin@orgreq.crowd.com/msp
    if [ $PEER -eq 0 ]; then
      CORE_PEER_ADDRESS=peer0.orgreq.crowd.com:7051	# go组织peer0地址
    else
      CORE_PEER_ADDRESS=peer1.orgreq.crowd.com:7051	# go组织peer1地址
    fi
  elif [ $ORG -eq 2 ]; then
    CORE_PEER_LOCALMSPID="OrgWorkMSP"	# CPP组织ID
    CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORGWORK_CA	# cli可以连接CPP组织的peer0节点
    # CPP组织的管理员证书
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgwork.crowd.com/users/Admin@orgwork.crowd.com/msp
    if [ $PEER -eq 0 ]; then
      CORE_PEER_ADDRESS=peer0.orgwork.crowd.com:7051	# cpp组织peer0地址
    else
      CORE_PEER_ADDRESS=peer1.orgwork.crowd.com:7051	# cpp组织peer0地址
    fi

  elif [ $ORG -eq 3 ]; then
    CORE_PEER_LOCALMSPID="OrgValidMSP" # CPP组织ID
    CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORGVALID_CA # cli可以连接CPP组织的peer0节点
    # CPP组织的管理员证书
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgvalid.crowd.com/users/Admin@orgvalid.crowd.com/msp
    if [ $PEER -eq 0 ]; then
      CORE_PEER_ADDRESS=peer0.orgvalid.crowd.com:7051  # cpp组织peer0地址
    elif [ $PEER  -eq 1 ]; then
      CORE_PEER_ADDRESS=peer1.orgvalid.crowd.com:7051  # cpp组织peer0地址
    elif [ $PEER -eq 2 ]; then
      CORE_PEER_ADDRESS=peer2.orgvalid.crowd.com:7051
    elif [ $PEER -eq 3 ]; then
      CORE_PEER_ADDRESS=peer3.orgvalid.crowd.com:7051
    else 
      CORE_PEER_ADDRESS=peer4.orgvalid.crowd.com:7051
      #statements
      #statements
      #statements
    fi
  else
    echo "================== ERROR !!! ORG Unknown =================="
  fi
  if [ "$VERBOSE" == "true" ]; then
    env | grep CORE
  fi
}

# 该函数对命令执行的结果进行验证, 如果验证失败, 直接退出, 不再执行后续操作流程
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
    echo "========= ERROR !!! FAILED ==========="
    echo
    exit 1
  fi
}