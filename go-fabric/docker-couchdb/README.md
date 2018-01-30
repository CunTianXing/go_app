一.配置CouchDB+Fabric环境


Fabric可能会遇到的问题
虽然区块链是一个只能插入和查询的数据库，但是我们的业务数据是存放在State Database中的，如果我们直接修改了CouchDB的数据，那么接下来的查询和事务是直接基于修改后的CouchDB的，并不会去检查区块链中的记录，所以理论上是可以通过直接改CouchDB来实现业务数据的修改。我们以官方的Marble为例，看看修改CouchDB会怎么样？


ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
OR
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/cacerts/ca.example.com-cert.pem

peer chaincode install -n marbles02 -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/marbles02

peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C mychannel -n marbles02 -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.member','Org2MSP.member')"

peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C mychannel -n marbles02 -c '{"Args":["initMarble","marble2","red","50","tom"]}'

peer chaincode query -C mychannel -n marbles02 -c '{"Args":["readMarble","marble2"]}'

我们可以看到通过curl直接查询CouchDB中的数据：

curl http://127.0.0.1:5984/mychannel/marbles02%00marble2

{"_id":"marbles02\u0000marble2","_rev":"1-a1844f47b9ed94294b430c9a9a6f543b","chaincodeid":"marbles02","data":{"docType":"marble","name":"marble2","color":"red","size":50,"owner":"tom"},"version":"6:0"}

如果我们要修改其中的数据，把颜色改成green，大小改成10，那么我们可以运行：

curl -X PUT http://127.0.0.1:5984/mychannel/marbles02%00marble2 -d '{"_id":"marbles02\u0000marble2","_rev":"1-a1844f47b9ed94294b430c9a9a6f543b","chaincodeid":"marbles02","data":{"docType":"marble","name":"marble2","color":"green","size":10,"owner":"tom"},"version":"6:0"}'
系统返回结果：
{"ok":true,"id":"marbles02\u0000marble2","rev":"2-6ffc6652cfc707f8352a5f06c3ce1ce6"}

我们在4个CouchDB中都运行这个命令，把4个数据库的数据都改了。
接下来我们通过ChainCode来查询，看看会怎么样。
peer chaincode query -C mychannel -n marbles02 -c '{"Args":["readMarble","marble2"]}'
返回结果：
Query Result: {"color":"green","docType":"marble","name":"marble2","owner":"tom","size":10}
可以看到数据已经变成新的值，那么接下来运行其他的Transaction会怎么样？我们试一试转账操作：
peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C mychannel -n marbles02 -c '{"Args":["transferMarble","marble2","jerry"]}'
系统返回成功，我们再查一下呢
peer chaincode query -C mychannel -n marbles02 -c '{"Args":["readMarble","marble2"]}'
Query Result: {"color":"green","docType":"marble","name":"marble2","owner":"jerry","size":10}
所以我们对CouchDB数据库的更改都是有效的，在Fabric看来似乎并不知道我们改了CouchDB的内容。
