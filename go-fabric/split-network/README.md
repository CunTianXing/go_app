## Build Split Network (BYFN)

#### ./byfn.sh -m up -c mychannel
#### docker exec -it cli bash
#### cd scripts
#### source setclienv.sh
#### ./channel-setup.sh

# Chaincode
### ./install-chaincode.sh 1.0
#### ./instantiate-chaincode.sh 1.0
##dev-peer0.org1.example.com-bigdatacc-1.0
#### ./update-invoke.sh myvar 100 +
#### ./get-invoke.sh myvar


### ./install-chaincode.sh 2.0
#### ./upgrade-chaincode.sh 2.0
##dev-peer0.org1.example.com-bigdatacc-2.0
#### ./get-invoke.sh myvar
#### ./delete-invoke.sh myvar


##dev-peer0.org1.example.com-bigdatacc-3.0
###bigdatacc  export CC_NAME=bigdatacc
