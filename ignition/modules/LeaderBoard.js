const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

const LeaderBoardModule = buildModule("LeaderBoardModule", (m) => {
  const leaderBoard = m.contract("GameHistory");

  return { leaderBoard };
});

module.exports = LeaderBoardModule;