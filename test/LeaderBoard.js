const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("GameHistory", function () {
 let GameHistory, gameHistory, owner, addr1;
 let ownerAddress, addr1Address;

 beforeEach(async function () {
    GameHistory = await ethers.getContractFactory("GameHistory");
    [owner, addr1] = await ethers.getSigners();

    ownerAddress = await owner.getAddress();
    addr1Address = await addr1.getAddress();

    // Deploy the contract
    gameHistory = await GameHistory.deploy();
    // await gameHistory.deployTransaction.wait();
 });

 describe("Game Data Storage", function () {
    it("Should store game data correctly", async function () {
      await gameHistory.storeGameData(1, "gtid1", "uid1", "data1", 123456789);
      const gameHistory1 = await gameHistory.getGameHistory(1);
      expect(gameHistory1.length).to.equal(1);
      expect(gameHistory1[0].gid).to.equal(1);
      expect(gameHistory1[0].gtid).to.equal("gtid1");
      expect(gameHistory1[0].uid).to.equal("uid1");
      expect(gameHistory1[0].data).to.equal("data1");
      expect(gameHistory1[0].time).to.equal(123456789);

      const userHistory1 = await gameHistory.getUserHistory("uid1");
      expect(userHistory1.length).to.equal(1);
      expect(userHistory1[0].gid).to.equal(1);
      expect(userHistory1[0].gtid).to.equal("gtid1");
      expect(userHistory1[0].uid).to.equal("uid1");
      expect(userHistory1[0].data).to.equal("data1");
      expect(userHistory1[0].time).to.equal(123456789);
    });

    it("Should not allow non-owner to store game data", async function () {
      await expect(gameHistory.connect(addr1).storeGameData(2, "gtid2", "uid2", "data2", 123456790)).to.be.revertedWith("Caller is not the owner");
    });
 });

 describe("Game History Retrieval", function () {
    it("Should return the correct game history", async function () {
      await gameHistory.storeGameData(1, "gtid1", "uid1", "data1", 123456789);
      const gameHistory1 = await gameHistory.getGameHistory(1);
      expect(gameHistory1.length).to.equal(1);
      expect(gameHistory1[0].gid).to.equal(1);
      expect(gameHistory1[0].gtid).to.equal("gtid1");
      expect(gameHistory1[0].uid).to.equal("uid1");
      expect(gameHistory1[0].data).to.equal("data1");
      expect(gameHistory1[0].time).to.equal(123456789);
    });
 });

 describe("User History Retrieval", function () {
    it("Should return the correct user history", async function () {
      await gameHistory.storeGameData(1, "gtid1", "uid1", "data1", 123456789);
      const userHistory1 = await gameHistory.getUserHistory("uid1");
      expect(userHistory1.length).to.equal(1);
      expect(userHistory1[0].gid).to.equal(1);
      expect(userHistory1[0].gtid).to.equal("gtid1");
      expect(userHistory1[0].uid).to.equal("uid1");
      expect(userHistory1[0].data).to.equal("data1");
      expect(userHistory1[0].time).to.equal(123456789);
    });
 });
});

describe("TokenSender", function () {
  let TokenSender, token, user, owner;

  beforeEach(async function () {
    // Deploy the ERC20 token
    const Token = await ethers.getContractFactory("GLDToken");
    const amount = ethers.utils.parseEther("1000");
    token = await Token.deploy(amount);
    await token.deployed();

    // Mint some tokens to the contract
    // const amount = ethers.utils.parseEther("1000");
    await token.mint(address(this), amount);

    // Deploy the TokenSender contract
    TokenSender = await ethers.getContractFactory("GameHistory");
    [owner, user] = await ethers.getSigners();
    const tokenAddressEGC = token.address;
    const contract = await TokenSender.deploy(tokenAddressEGC);
    await contract.deployed();
  });

  it("Should send tokens from the contract to a user", async function () {
    const amount = ethers.utils.parseEther("500");
    await expect(contract.sendEgc(user.address, amount))
    .to.emit(token, "Transfer")
    .withArgs(owner.address, user.address, amount);
  });
});
