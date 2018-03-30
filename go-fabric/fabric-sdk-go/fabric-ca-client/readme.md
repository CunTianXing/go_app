##
#### export FABRIC_CA_CLIENT_HOME=/data/go/src/github.com/CunTianXing/go_app/go-fabric/fabric-sdk-go/fabric-ca-client/admin
#### fabric-ca-client enroll -u https://admin:adminpw@localhost:7054
#### fabric-ca-client register -u https://admin:adminpw@localhost:7054 --id.name admin2 --id.type user --id.affiliation org1.department1 --id.attrs 'hf.Revoker=true,admin=true:ecert'
#### export FABRIC_CA_CLIENT_HOME=/data/go/src/github.com/CunTianXing/go_app/go-fabric/fabric-sdk-go/fabric-ca-client/admin
#### fabric-ca-client register -u https://admin:adminpw@localhost:7054 --id.name peer1 --id.type peer --id.affiliation org1.department1 --id.secret peer1pw
#### export FABRIC_CA_CLIENT_HOME=/data/go/src/github.com/CunTianXing/go_app/go-fabric/fabric-sdk-go/fabric-ca-client/peer1
#### fabric-ca-client enroll -u https://peer1:peer1pw@localhost:7054 -M $FABRIC_CA_CLIENT_HOME/msp
#### fabric-ca-client reenroll -u https://admin:adminpw@localhost:7054
#### fabric-ca-client revoke -u https://admin:adminpw@localhost:7054 -e peer1

#### fabric-ca-client enroll -d --enrollment.profile tls -u https://admin:adminpw@localhost:7054 -M /tmp/tls --csr.hosts ca.org1.example.com
