// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.24;

/** 
 * @title Game records
 * @dev Getter and setter methods 
 */
contract GameHistory {

    struct GameSession {
        uint gid;     //game ID
        string gtid;  //game ID
        string uid;   //user ID
        string data;  //game data
        uint time;    //game time stamp
    }

    address private immutable owner;

    // event for EVM logging
    event OwnerSet(address indexed oldOwner, address indexed newOwner);

    mapping(uint => GameSession[]) gameHistory; //historical data for each game 
    mapping(string => GameSession[]) userHistory; //historical data for each user

    constructor() {
        owner = msg.sender; // 'msg.sender' is sender of current call, contract deployer for a constructor
        emit OwnerSet(address(0), owner);
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Caller is not the owner");
        _;
    }

    /**
     * @dev Stores the game data for the given game ID, game type ID, user ID, data, and time.
     * @param _gid The unique game ID.
     * @param _gtid The game type ID.
     * @param _uid The user ID.
     * @param _data The game data to be stored.
     * @param _time The timestamp of when the game data is stored.
     */
    function storeGameData(uint _gid, string memory _gtid, string memory _uid, string memory _data, uint _time) public onlyOwner {
        GameSession memory currentGame_ = GameSession(_gid, _gtid, _uid, _data, _time);

        gameHistory[_gid].push(currentGame_);
        userHistory[_uid].push(currentGame_);
    }

    /**
     * @dev Retrieves the game history for a specific game ID.
     * @param _gid The unique game ID for which the game history is to be retrieved.
     * @return An array of GameSession structs representing the game history for the specified game ID.
     */
    function getGameHistory(uint _gid) public view returns (GameSession[] memory) {
        return gameHistory[_gid];
    }

    /**
     * @dev Retrieves the user history for a specific user ID.
     * @param _uid The unique user ID for which the user history is to be retrieved.
     * @return An array of GameSession structs representing the user history for the specified user ID.
     */
    function getUserHistory(string memory _uid) public view returns (GameSession[] memory) {
        return userHistory[_uid];
    }

}