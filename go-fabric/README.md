type ChaincodeStubInterface
type ChaincodeStubInterface interface {
    // GetArgs返回用于Chaincode Init和Invoke的参数作为字节数组的数组。
    GetArgs() [][]byte

    // GetStringArgs将链式代码Init和Invoke的参数作为字符串数组返回。 如果客户端传递用作字符串的参数，则只使用GetStringArgs。
    GetStringArgs() []string

    // GetFunctionAndParameters将第一个参数作为函数名称，其余参数作为字符串数组中的参数返回。 如果客户端传递用作字符串的参数，则只使用GetFunctionAndParameters。
    GetFunctionAndParameters() (string, []string)

    // GetArgsSlice将链式代码Init和Invoke的参数作为字节数组返回
    GetArgsSlice() ([]byte, error)

    // GetTxID返回交易提议的tx_id (see ChannelHeader in protos/common/common.proto)
    GetTxID() string

    // InvokeChaincode使用相同的事务上下文在本地调用指定的链式代码`Invoke`; 也就是说，chaincode调用chaincode不会创建一个新的事务消息。 如果被调用的链接代码位于同一个通道上，则只需将被调用的链接代码读取集和写入集添加到调用事务中即可。 如果被调用的链式码在不同的频道上，则只有Response被返回给调用的链式码。 来自被调用的链接代码的任何PutState调用都不会对分类帐产生任何影响; 也就是说，不同渠道上的被调用的链式代码将不会将其读取集和写入集应用于事务。 只有调用链代码的读集和写集将被应用于事务。 有效地，不同渠道上的被调用的链式代码是一个“查询”，它在随后的提交阶段不参与状态验证检查。 如果`channel`为空，则假定主叫方的频道。
    InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response

    // GetState从分类账中返回指定`key`的值。 请注意，GetState不从没有提交到分类帐的写集读取数据。 换句话说，GetState不考虑PutState所修改的未被提交的数据。 如果密钥不存在于状态数据库中，则返回（nil，nil）。
    GetState(key string) ([]byte, error)

    // PutState将指定的`key`和`value`作为数据写入提议放入事务的写入集中。 直到事务验证并成功提交，PutState才会生效。 简单键不能是空字符串，也不能以空字符（0x00）开头，以避免与复合键（它们在内部以0x00作为复合键命名空间前缀）发生范围查询冲突。
    PutState(key string, value []byte) error

    // DelState在交易提议的writeset中记录指定的要删除的“key”。 事务验证并成功提交时，`key`及其值将从分类帐中删除。
    DelState(key string) error

    // GetStateByRange返回分类账中一组键的范围迭代器。 迭代器可用于迭代startKey（包含）和endKey（不包括）之间的所有键。 密钥由迭代器以词汇顺序返回。 请注意，startKey和endKey可以是空字符串，这意味着无限的范围查询开始或结束。 完成后调用Close（）返回的StateQueryIteratorInterface对象。 查询在验证阶段重新执行，以确保结果集自交易认可（检测到幻像读取）后未发生变化。
    GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error)

    // GetStateByPartialCompositeKey根据给定的部分组合键查询分类帐中的状态。 这个函数返回一个迭代器，它可以用来遍历所有组合键的前缀匹配给定的部分组合键。 objectType和属性应该只有有效的utf8字符串，不应该包含U + 0000（零字节）和U + 10FFFF（最大和未分配的代码点）。 请参阅相关函数SplitCompositeKey和CreateCompositeKey。 完成后，在返回的StateQueryIteratorInterface对象上调用Close（）。在验证阶段重新执行查询，以确保结果集自交易认可（检测到幻像读取）后未发生更改。
    GetStateByPartialCompositeKey(objectType string, keys []string) (StateQueryIteratorInterface, error)

    // CreateCompositeKey将给定的`attributes`组合成一个组合键。 objectType和属性应该只有有效的utf8字符串，不应该包含U + 0000（零字节）和U + 10FFFF（最大和未分配的代码点）。 生成的组合键可以用作PutState（）中的键。
    CreateCompositeKey(objectType string, attributes []string) (string, error)

    // SplitCompositeKey将指定的键分割成组合键所形成的属性。 在范围查询或部分复合键查询期间找到的复合键因此可以分解为它们的复合部分。
    SplitCompositeKey(compositeKey string) (string, []string, error)

    // GetQueryResult对状态数据库执行“丰富”查询。 它仅支持支持丰富查询的状态数据库，例如，CouchDB。 查询字符串是基础状态数据库的本地语法。 一个迭代器被返回，它可以被用来迭代（下一个）在查询结果集上。 查询不在验证阶段重新执行，幻像读取未检测到。 也就是说，其他提交的事务可能会添加，更新或删除影响结果集的键，而这在验证/提交时不会被检测到。 因此，易受此影响的应用程序不应使用GetQueryResult作为更新分类账的交易的一部分，并应将使用限制在只读链式操作。
    GetQueryResult(query string) (StateQueryIteratorInterface, error)

    // GetHistoryForKey在一段时间内返回键值的历史记录。 对于每个历史密钥更新，都会返回历史值和关联的事务标识和时间戳。 时间戳是客户端在提案标题中提供的时间戳。 GetHistoryForKey需要对等配置core.ledger.history.enableHistoryDatabase为true。在验证阶段不会重新执行查询，不会检测到幻像读取。 也就是说，其他提交的事务可能已经同时更新了密钥，影响了结果集，而这在验证/提交时不会被检测到。 因此，易受此影响的应用程序不应使用GetHistoryForKey作为更新分类账的交易的一部分，而应限制使用只读链代码操作。
    GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error)

    // GetCreator返回`SignedProposal`的`SignatureHeader.Creator`（例如标识）。 这是提交交易的代理（或用户）的身份。
    GetCreator() ([]byte, error)

    // GetTransient返回`ChaincodeProposalPayload.Transient`字段。 它是包含可能用于实现某种形式的应用程序级别机密性的数据（例如密码材料）的映射。 这个字段的内容，按照“ChaincodeProposalPayload”的规定，应该总是从交易中被省略，并从分类帐中排除。
    GetTransient() (map[string][]byte, error)

    // GetBinding返回事务绑定
    GetBinding() ([]byte, error)

    // GetSignedProposal返回SignedProposal对象，其中包含事务提议的所有数据元素部分
    GetSignedProposal() (*pb.SignedProposal, error)

    // GetTxTimestamp返回创建事务时的时间戳。 这是从事务ChannelHeader中获取的，因此它将指示客户端的时间戳，并在所有代言人中具有相同的值。
    GetTxTimestamp() (*timestamp.Timestamp, error)

    // SetEvent允许chaincode在事务提议上提出一个事件。 如果事务被验证并成功提交，事件将被传递给当前的事件监听器。
    SetEvent(name string, payload []byte) error
}

docker tag hyperledger/fabric-orderer:x86_64-1.0.0 hyperledger/fabric-orderer:latests
create database fabricca default character set utf8mb4 collate utf8mb4_unicode_ci;
 GRANT ALL PRIVILEGES ON  *.* TO xingcuntian@'%' IDENTIFIED BY 'xingcuntian';
FLUSH PRIVILEGES;


docker exec -it cli bash

安装并初始化我们的ChainCode：

peer chaincode install -n fsmtest -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/fsmtest
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer chaincode instantiate -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C mychannel -n fsmtest -v 1.0 -c '{"Args":[]}'

安装完毕后，我们可以起草一个报销单EXP1：
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile $ORDERER_CA -C mychannel -n fsmtest -c '{"Args":["Draft","EXP1"]}'

现在状态是Draft，然后我们试一试提交报销单EXP1：
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile $ORDERER_CA -C mychannel -n fsmtest -c '{"Args":["Submit","EXP1"]}'

看到状态已经改为Submitted了。接下来我们进一步一级审批通过，二级审批通过，都是执行相同的命令：
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile $ORDERER_CA -C mychannel -n fsmtest -c '{"Args":["Approve","EXP1"]}'

peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile $ORDERER_CA -C mychannel -n fsmtest -c '{"Args":["Approve","EXP1"]}'
这个时候，状态已经是Complete了，如果我们再次调用Approve函数会怎么样？因为我们在状态机中并没有定义这么一个流转事件，所以肯定是报错，无法正常执行的：
总的来说，在Fabric的ChainCode开发中，引入第三方的库可以方便我们编写更强大的链上代码。而这个FSM虽然简单，但是也可以很好的将状态流转的逻辑进行集中，避免了在状态流转时编写大量的Ugly的代码，让我们在每个函数中更专注于业务逻辑，而不是麻烦的状态转移。最后直接粘贴出我的完整ChainCode 源码，方便大家直接使用。
