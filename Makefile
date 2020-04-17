.PHONY: all dev clean build env-up env-down run destroy

all: clean build env-up run

dev: build run

####BUILD
#### @dep ensure
build:
	@echo "Build ..."
	@go build
	@echo "build done"

####ENV
env-up:
	@echo "Start environment ..."
	@cd fixtures && docker-compose up --force-recreate -d
	@echo "Enviroment up"

env-down:
	@echo "Stop environment ..."
	@cd fixtures && docker-compose down 
	@echo "Enviroment down"

####RUN
run:
	@echo "Start app ..."
	@./Demo

####CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/Demo-* Demo
	@docker volume prune
	@docker rm -f -v `docker ps -a --no-trunc | grep "Demo" | cut -d ' ' -f 1` 2>/dev/null || true 
	@docker rmi `docker images --no-trunc | grep "Demo" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"

destroy: 
	@echo "destroy network ..."
	@cd fixtures && rm -fr channel-artifacts/*
	@cd fixtures && rm -fr crypto-config/*
	@cd fixtures && cryptogen generate --config=./crypto-config.yaml
	@cd fixtures && configtxgen -profile crowdOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
	@cd fixtures && configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mychannel
	@cd fixtures && configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/OrgReqMSPanchors.tx -channelID mychannel -asOrg OrgReq
	@cd fixtures && configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/OrgWorkMSPanchors.tx -channelID mychannel -asOrg OrgWork

# 	@cd fixtures && configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/OrgValidMSPanchors.tx -channelID mychannel -asOrg OrgValid
updatechannel:
	@echo "update the channel"
	@cd fixtures && rm -fr channel-artifacts/*
	@cd fixtures && configtxgen -profile crowdOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
	@cd fixtures && configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mychannel
	@cd fixtures && configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/OrgReqMSPanchors.tx -channelID mychannel -asOrg OrgReq
# 	@cd fixtures && configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/OrgWorkMSPanchors.tx -channelID mychannel -asOrg OrgWork

newchannel:
# 	ifdef channelname
	@echo "create a new channel ..."
	@cd fixtures && configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./channel-artifacts/${channelname}.tx -channelID ${channelname}
# 	@cd fixtures && configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/testOrgReqMSPanchors.tx -channelID testchannel -asOrg OrgReq
# 	@cd fixtures && configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/testOrgWorkMSPanchors.tx -channelID testchannel -asOrg OrgWork
#     endif 
