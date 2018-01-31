CouchDB安装
下面我们来说一说这个CouchDB。

CouchDB是一个完全局域RESTful API的键值数据库，也就是说我们不需要任何客户端，只需要通过HTTP请求就可以操作数据库了。LevelDB是Peer的本地数据库，那么肯定是和Peer一对一的关系，那么CouchDB是个网络数据库，应该和Peer是什么样一个关系呢？在生产环境中，我们会为每个组织部署节点，而且为了高可用，可能会在一个组织中部署多个Peer。同样我们在一个组织中也部署多个CouchDB，每个Peer对应一个CouchDB。

HyperLedger在Docker Hub上也发布了CouchDB的镜像，为了能够深入研究CouchDB和Fabric的集成，我们就采用官方发布的CouchDB来做。

docker pull klaemo/couchdb
【注意，如果我们是docker pull couchdb，那么只能获得1.6版本的CouchDB，而要获得最新的2.0版，就需要用上面这个镜像。】

可以获得官方的CouchDB镜像。CouchDB在启动的时候需要指定一个本地文件夹映射成CouchDB的数据存储文件夹，所以我们可以在当前用户的目录下创建一个文件夹用于存放数据。

mkdir couchdb
下载完成后，我们只需要执行以下命令即可启用一个CouchDB的实例：

docker run -p 5984:5984 -d --name my-couchdb -e COUCHDB_USER=admin -e COUCHDB_PASSWORD=password -v ~/couchdb:/opt/couchdb/data klaemo/couchdb
启动后我们打开浏览器，访问Linux的IP的5984端口的URL，比如我的Linux是http://127.0.0.1，那么URL是：
http://127.0.0.1:5984/_utils
这个时候我们就可以看到CouchDB的Web管理界面了。输入用户名admin密码password即可进入。
现在是一个空数据库，我们将CouchDB和Peer结合起来后再看会是什么样的效果。
配置CouchDB+Fabric环境
先删除刚才创建的CouchDB容器：
docker rm -f my-couchdb
首先我们是4个Peer+1Orderer的模式，所以我们先创建4个CouchDB数据库：
cd ~
mkdir couchdb0
mkdir couchdb1
mkdir couchdb2
mkdir couchdb3
docker run -p 5984:5984 -d --name couchdb0 -e COUCHDB_USER=admin -e COUCHDB_PASSWORD=password -v ~/couchdb0:/opt/couchdb/data klaemo/couchdb
docker run -p 6984:5984 -d --name couchdb1 -e COUCHDB_USER=admin -e COUCHDB_PASSWORD=password -v ~/couchdb1:/opt/couchdb/data klaemo/couchdb
docker run -p 7984:5984 -d --name couchdb2 -e COUCHDB_USER=admin -e COUCHDB_PASSWORD=password -v ~/couchdb2:/opt/couchdb/data klaemo/couchdb
docker run -p 8984:5984 -d --name couchdb3 -e COUCHDB_USER=admin -e COUCHDB_PASSWORD=password -v ~/couchdb3:/opt/couchdb/data klaemo/couchdb
然后我们需要启动Fabric了

官方已经提供了多个Docker-compose文件，如果我们使用的是./network_setup.sh up命令，那么启用的就是docker-compose-cli.yaml这个文件。如果要基于这个yaml文件启用CouchDB的Peer，则打开该文件，并编辑其中的Peer节点，改为如下的形式：

peer0.org1.example.com:
  container_name: peer0.org1.example.com
  environment:
    - CORE_LEDGER_STATE_STATEDATABASE=CouchDB
    - CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=127.0.0.1:5984
    - CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=admin
    - CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=password
  extends:
    file:  base/docker-compose-base.yaml
    service: peer0.org1.example.com

这里的192.168.100.129:5984是我映射CouchDB后的Linux的IP地址和IP。然后是设置用户名和密码。把4个Peer的配置都改好后，保存，我们试着启用Fabric：

./network_setup.sh up
等Fabric启动完成并运行了ChainCode测试后，我们刷新http://127.0.0.1:5984/_utils ，可以看到以Channel名字创建的Database，另外还有几个是系统数据库。

点进mychannel数据库，我们可以看到其中的数据内容。点击“Mango Query”可以编写查询，默认提供的查询可以点击Run Query按钮查询所有的数据结果：
CouchDB的直接查询
接下来我们使用Linux的curl来查询CouchDB数据库。
比如我们要看看mychannel数据库下有哪些数据：
curl http://127.0.0.1:5984/mychannel/_all_docs
可以看到我运行了一些ChainCode后的State DATABASE结果：
{"total_rows":7,"offset":0,"rows":[
{"id":"devincc\u0000a","key":"devincc\u0000a","value":{"rev":"2-a979bf6c2716ecae6d106999f833a59c"}},
{"id":"devincc\u0000b","key":"devincc\u0000b","value":{"rev":"2-ad1c549305fd277097180405f96bdcd8"}},
{"id":"lscc\u0000devincc","key":"lscc\u0000devincc","value":{"rev":"1-05d2cd0b344c4dd8a8d1a3ffd7332544"}},
{"id":"lscc\u0000mycc","key":"lscc\u0000mycc","value":{"rev":"1-2cba0344b1610b9d9254bbafbda5e9b1"}},
{"id":"mycc\u0000a","key":"mycc\u0000a","value":{"rev":"2-588a45b289359afa9dc6e5e7866eaf97"}},
{"id":"mycc\u0000b","key":"mycc\u0000b","value":{"rev":"2-54e6639a858b0f91298c9a354484513a"}},
{"id":"statedb_savepoint","key":"statedb_savepoint","value":{"rev":"10-6ccde2a55c71d7d6a70d9333d119fc8e"}}
 ]}

如果我们要查询其中的一条数据，只需要用/ChannelId/id来查询，比如查询：statedb_savepoint

curl http://127.0.0.1:5984/mychannel/statedb_savepoint
返回的结果：
{"_id":"statedb_savepoint","_rev":"10-6ccde2a55c71d7d6a70d9333d119fc8e","BlockNum":4,"TxNum":0,"UpdateSeq":"19-g1AAAAEzeJzLYWBg4MhgTmHgzcvPy09JdcjLz8gvLskBCjMlMiTJ____PyuRAYeCJAUgmWQPVsOCS40DSE08WA0jLjUJIDX1eO3KYwGSDA1ACqhsPiF1CyDq9mclsuJVdwCi7j4h8x5A1AHdx5kFAI6sYwk"}

麻烦的是业务数据是“ChainCodeName\u0000数据”这样的格式的ID，而如果我们要通过这个ID查询，那么就根本找不到啊！
curl http://127.0.0.1:5984/mychannel/mycc\u0000a
{"error":"not_found","reason":"missing"}

正确的做法是把\u0000替换为%00，也就是说我们的查询应该是：

curl http://127.0.0.1:5984/mychannel/mycc%00a
正确返回结果：
{"_id":"mycc\u0000a","_rev":"2-588a45b289359afa9dc6e5e7866eaf97","chaincodeid":"mycc","version":"4:0","_attachments":{"valueBytes":{"content_type":"application/octet-stream","revpos":2,"digest":"md5-hhOYXsSeuPdXrmQ56Hm7Kg==","length":2,"stub":true}}}

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
