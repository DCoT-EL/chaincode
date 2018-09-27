# Chaincode

####For install a HyperLedger Fabric network see this link [configuration-network-fabric](https://github.com/ascatox/configuration-network-fabric)
###Download
```bash
$ cd configuration-network-fabric/chaincode/go
$ git clone https://github.com/DCoT-EL/chaincode.git
```

###Istallation

For first installation:

```bash
$ docker exec -it cli bash

$ peer chaincode install -p github.com/hyperledger/fabric/examples/chaincode/go/dcot-chaincode -n dcot-chaincode -v 1.0

$ peer chaincode instantiate -n dcot-chaincode -c '{"Args":["a","10"]}' -C ledgerchannel -v 1.0
```
At every modification of the chaincode , you must use a  *`upgrade`* command :

```bash
$ docker exec -it cli bash

$ peer chaincode install -p github.com/hyperledger/fabric/examples/chaincode/go/dcot-chaincode -n dcot-chaincode -v [version upgrade]

$ peer chaincode upgrade -n dcot-chaincode -c '{"Args":["a","10"]}' -C ledgerchannel -v [version upgrade] 
```

If there are no errors:
```bash
$ exit
```
Check if everything has gone correctly 
```bash
$ docker ps  -a
$ docker logs dev-peer0.org1.example.com-dcot-chaincode-1.0 
```



You must see a login like this: `[dcot-chaincode-log] Info -> INFO 001 Initializing Chain of Custody
`







*PS: Commands tested with Ubuntu 16.04*
