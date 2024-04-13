# Leader Board for Get Easy Coin

This project contains both the smart contract and the go server running the leaderboard.

## To test the contract

The smart contract was tested using hardhat and deployed on Sepolia Ethereum for tests. To test the contract using hardhat.

### Current test coverage
------------------|----------|----------|----------|----------|----------------|
File              |  % Stmts | % Branch |  % Funcs |  % Lines |Uncovered Lines |
------------------|----------|----------|----------|----------|----------------|
 contracts/       |      100 |      100 |      100 |      100 |                |
  LeaderBoard.sol |      100 |      100 |      100 |      100 |                |
------------------|----------|----------|----------|----------|----------------|
All files         |      100 |      100 |      100 |      100 |                |
------------------|----------|----------|----------|----------|----------------|

### Test using js evm
* Install the required packages
```shell
npm install
```
* Test and see coverage report
```shell
npx hardhat coverage
```

### Test on Sepolia Ethereum using Infura nodes
* Set hard hat variables
```shell
npx hardhat vars set INFURA_API_KEY
# input your infura api key
npx hardhat vars set SEPOLIA_PRIVATE_KEY
# input your Sepolia private key
```

* Deploy the contract
```shell
npx hardhat ignition deploy ./ignition/modules/LeaderBoard.js --network sepolia
```

For deployment on another chain, the contract can be found in ./contracts/LeaderBoard.sol


## The Go Server
The server was built using Go version 1.22.1


### Set up the project
* Install [Go compiler](https://go.dev/doc/install)

* Create an env file following the format of env.example
`
INFURA_APIKEY=
PRIVATE_KEY=
CONTRACT_ADDRESS=
NODE_URL=
MODE=release
API_VERSION=
API_HOST=
PORT=
ORIGIN=
`

### Test the project
To test the go packages 
```shell
go test ./...
```

### Run the server
```shell
go run ./...
```

### Docs
Swagger docs can be found at /api/swagger/index.html
