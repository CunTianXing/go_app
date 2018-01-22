#监听
连接到对等节点以接收块和链接代码事件（如果有链式代码事件正在发送）。 目前只有这个例子在环境中禁用TLS。
#构建
```sh
1. go build

2. ./block-listener -events-address=<peer-address> -events-from-chaincode=<chaincode-id> -events-mspdir=<msp-directory> -events-mspid=<msp-id>
```
请注意，如果没有，则将使用fabric/sampleconfig下的默认MSP
提供MSP参数。

# Example with the docker_cli example
In order to use the block listener with the e2e_cli example, make sure that TLS
has been disabled by setting CORE_PEER_TLS_ENABLED=***false*** and ORDERER_GENERAL_TLS_ENABLED=***false*** in
``docker-compose-cli.yaml``, ``base/docker-compose-base.yaml`` and
``base/peer-base.yaml``.

Once the "All in one" command:
```sh
./network_setup.sh up
```
has completed, attach the event client to peer peer0.org1.example.com by doing
the following (assuming you are running block-listener in the host environment):
```sh
./block-listener -events-address=127.0.0.1:7053 -events-mspdir=$GOPATH/src/github.com/CunTianXing/go_app/go-fabric/docker-cli/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp -events-mspid=Org1MSP

```

The event client should output "Event Address: 127.0.0.1:7053" and wait for
events.

Exec into the cli container:

```sh
docker exec -it cli bash
```
Setup the environment variables for peer0.org1.example.com
```sh
CORE_PEER_MSPCONFIGPATH=$GOPATH/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
CORE_PEER_ADDRESS=peer0.org1.example.com:7051
CORE_PEER_LOCALMSPID="Org1MSP"
```

Create an invoke transaction:

```sh
peer chaincode invoke -o orderer.example.com:7050 -C $CHANNEL_NAME -n mycc -c '{"Args":["invoke","a","b","10"]}'
```
Now you should see the block content displayed in the terminal running the block
listener.
