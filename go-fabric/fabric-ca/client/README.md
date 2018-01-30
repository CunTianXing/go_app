export FABRIC_CA_CLIENT_HOME=$GOPATH/src/github.com/CunTianXing/go_app/go-fabric/fabric-ca/client
fabric-ca-client enroll -u http://admin:adminpw@localhost:7054

fabric-ca-client register --id.name xingcuntian --id.type user --id.affiliation org1.department1 --id.attrs 'hf.Revoker=true,foo=bar'

fabric-ca-client enroll -u http://xingcuntian:XPQDlZBdblus@localhost:7054 -M $FABRIC_CA_CLIENT_HOME/xingcuntian

docker exec -it cli bash

peer chaincode install -n fabricca -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/fabricca

ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

peer chaincode instantiate -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C mychannel -n fabricca -v 1.0 -c '{"Args":[]}'

peer chaincode query -C mychannel -n fabricca -c '{"Args":["cert"]}'

cd $GOPATH/src/github.com/CunTianXing/go_app/go-fabric/docker-clis/crypto-config/peerOrganizations/org1.example.com/users
mkdir xingcuntian/msp

cp ./xingcuntian/ xingcuntian/msp –R

cp -rf signcerts  admincerts


CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/xingcuntian/msp

peer chaincode query -C mychannel -n fabricca -c '{"Args":["cert"]}'
