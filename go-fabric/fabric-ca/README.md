1.修改docker-compose文件，增加CA容器
我们就以给org1这个组织增加CA容器为例，打开e2e_cli文件夹中的docker-compose-cli.yaml ，增加以下内容：
ca0:
  image: hyperledger/fabric-ca
  environment:
    - FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server
    - FABRIC_CA_SERVER_CA_NAME=ca0
    - FABRIC_CA_SERVER_TLS_ENABLED=false
  ports:
    - "7054:7054"
  command: sh -c 'fabric-ca-server start --ca.certfile /etc/hyperledger/fabric-ca-server-config/ca.org1.example.com-cert.pem --ca.keyfile /etc/hyperledger/fabric-ca-server-config/${PRIVATE_KEY} -b admin:adminpw -d'
  volumes:
    - ./crypto-config/peerOrganizations/org1.example.com/ca/:/etc/hyperledger/fabric-ca-server-config
  container_name: ca0

  这里我们注意到，Fabric CA Server启动的时候，带了3个重要的参数：ca.certfile 指定了CA的根证书，ca.keyfile 指定了接下来给新用户签发证书时的私钥，这里我们使用变量${PRIVATE_KEY}代替，这是因为每次network_setup的时候，私钥的名字是不一样的，所以需要从启动脚本中传入。另外就是-b参数，指定了CA Client连接CA Server时使用的用户名密码

  2.修改network_setup.sh启动脚本，将CA容器启动的参数带入
接下来我们需要修改network_setup.sh文件，因为前面我们使用了变量${PRIVATE_KEY}，所以这里我们需要读取变量并带入docker-compose 启动的时候。具体修改如下：

function networkUp () {
    if [ -f "./crypto-config" ]; then
       echo "crypto-config directory already exists."
    else
       #Generate all the artifacts that includes org certs, orderer genesis block,
      # channel configuration transaction
       source generateArtifacts.sh $CH_NAME
    fi
folder="crypto-config/peerOrganizations/org1.example.com/ca"
privName=""
for file_a in ${folder}/*
do
    temp_file=`basename $file_a`

    if [ ${temp_file##*.} != "pem" ];then
       privName=$temp_file
    fi
done
    echo $privName
    if [ "${IF_COUCHDB}" == "couchdb" ]; then
      CHANNEL_NAME=$CH_NAME TIMEOUT=$CLI_TIMEOUT docker-compose -f $COMPOSE_FILE -f $COMPOSE_FILE_COUCH up -d 2>&1
    else
      CHANNEL_NAME=$CH_NAME TIMEOUT=$CLI_TIMEOUT PRIVATE_KEY=$privName docker-compose -f $COMPOSE_FILE up -d 2>&1
    fi
    if [ $? -ne 0 ]; then
        echo "ERROR !!!! Unable to pull the images "
        exit 1
     fi
    docker logs -f cli
}

这里脚本的逻辑很简单，就是去crypto-config/peerOrganizations/org1.example.com/ca这个文件夹中去遍历文件，找到私钥文件的文件名，并把文件名赋值给privName，然后在docker-compse的启动时，指定到PRIVATE_KEY即可。

3.使用CA Client生成新用户
只需要经过前面2步，我们给Org1设置的CA Server就算完成了。

3.1启动Fabric网络
运行

./network_setup.sh up

启动整个Fabric网络。接下来需要使用CA Client来生成新用户。我们需要以下几步：

3.2下载并安装Fabric CA Client
官方提供的CA Client需要依赖于libtool这个库，所以需要先安装这个库，运行命令：

sudo apt install libtool libltdl-dev
然后执行以下命令安装Fabric CA Client：
go get -u github.com/hyperledger/fabric-ca/cmd/...
该命令执行完毕后，我们应该在~/go/bin下面看到生成的2个文件：
fabric-ca-client  fabric-ca-server

3.3注册认证管理员
我们首先需要以管理员身份使用CA Client连接到CA Server，并生成相应的文件。
export FABRIC_CA_CLIENT_HOME=$HOME/ca
fabric-ca-client enroll -u http://admin:adminpw@localhost:7054

这个时候我们可以去$HOME/ca目录，看到CA Client创建了一个fabric-ca-client-config.yaml文件和一个msp文件夹。config可以去修改一些组织信息之类的。
3.4注册新用户

接下来我们想新建一个叫devin的用户，那么需要先执行这个命令：

fabric-ca-client register --id.name devin --id.type user --id.affiliation org1.department1 --id.attrs 'hf.Revoker=true,foo=bar'
系统会返回一个该用户的密码：
2017/09/05 22:20:41 [INFO] User provided config file: /home/studyzy/ca/fabric-ca-client-config.yaml
2017/09/05 22:20:41 [INFO] Configuration file location: /home/studyzy/ca/fabric-ca-client-config.yaml
Password: GOuMzkcGgGzq
我们拿到这个密码以后就可以再次使用enroll命令，给devin这个用户生成msp的私钥和证书：

fabric-ca-client enroll -u http://devin:GOuMzkcGgGzq@localhost:7054 -M $FABRIC_CA_CLIENT_HOME/devinmsp
现在新用户devin的私钥和证书就在$HOME/ca/devinmsp目录下，我们可以使用tree命令查看一下：

4.编写ChainCode验证当前用户
由于官方提供的example02并没有关于当前用户的信息的代码，所以我们需要编写自己的ChainCode。

这里我们主要是用到ChainCode接口提供的GetCreator方法，具体完整的ChainCode如下：

package main

import (
   "github.com/hyperledger/fabric/core/chaincode/shim"
   pb "github.com/hyperledger/fabric/protos/peer"
   "fmt"
   "encoding/pem"
   "crypto/x509"
   "bytes"
)

type SimpleChaincode struct {
}

func main() {
   err := shim.Start(new(SimpleChaincode))
   if err != nil {
      fmt.Printf("Error starting Simple chaincode: %s", err)
   }
}
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
   return shim.Success(nil)
}


func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
   function, args := stub.GetFunctionAndParameters()
   fmt.Println("invoke is running " + function)
   if function == "cert" {//自定义函数名称
      return t.testCertificate(stub, args)//定义调用的函数
   }
   return shim.Error("Received unknown function invocation")
}
func (t *SimpleChaincode) testCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response{
   creatorByte,_:= stub.GetCreator()
   certStart := bytes.IndexAny(creatorByte, "-----")// Devin:I don't know why sometimes -----BEGIN is invalid, so I use -----
   if certStart == -1 {
      fmt.Errorf("No certificate found")
   }
   certText := creatorByte[certStart:]
   bl, _ := pem.Decode(certText)
   if bl == nil {
      fmt.Errorf("Could not decode the PEM structure")
   }
   fmt.Println(string(certText))
   cert, err := x509.ParseCertificate(bl.Bytes)
   if err != nil {
      fmt.Errorf("ParseCertificate failed")
   }
   fmt.Println(cert)
   uname:=cert.Subject.CommonName
   fmt.Println("Name:"+uname)
   return shim.Success([]byte("Called testCertificate "+uname))
}
