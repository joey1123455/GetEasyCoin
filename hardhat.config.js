require("@nomicfoundation/hardhat-toolbox");
require("@nomicfoundation/hardhat-verify");
require("@nomiclabs/hardhat-ethers");

const { vars } = require("hardhat/config");
const INFURA_API_KEY = vars.get("INFURA_API_KEY");
const SEPOLIA_PRIVATE_KEY = vars.get("SEPOLIA_PRIVATE_KEY");

const POLYGON_AMOY_PRIVATE_KEY = vars.get("POLYGON_AMOY_PRIVATE_KEY");
const INFURA_POLYGON_AMOY_API_KEY = vars.get("INFURA_POLYGON_AMOY_API_KEY");

const OKLINK_API_KEY = vars.get("OKLINK_API_KEY");
/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.24",
  networks: {
    sepolia: {
      url: `https://sepolia.infura.io/v3/${INFURA_API_KEY}`,
      accounts: [SEPOLIA_PRIVATE_KEY],
    },
    polygonAmoy: {
      url: `https://polygon-amoy.infura.io/v3/${INFURA_POLYGON_AMOY_API_KEY}`,
      accounts: [POLYGON_AMOY_PRIVATE_KEY],
    },
  },
  etherscan: {
    apiKey: {
      polygonAmoy: OKLINK_API_KEY,
    },
    customChains: [
      {
        network: "polygonAmoy",
        chainId: 80002,
        urls: {
          apiURL: "https://www.oklink.com/api/expolrer/v1/contract/verify/async/api/polygonAmoy",
          browserURL: "https://www.oklink.com/polygonAmoy"
        },
      }
    ]
  }
}
