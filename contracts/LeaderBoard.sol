// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

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

    struct Payment {
        address sender; //sender's address
        uint amount;    //total amount sent
        uint time;      //timestamp of the payment
    }

    address private immutable owner;
    uint256 private immutable oneEGC;
    address public tokenAddressEGC;
    // address public tokenAddressEGC = address(0xa2630a1178bA5774395Aa3a21AfDD4c3E654a612);
    

    // event for EVM logging
    event OwnerSet(address indexed oldOwner, address indexed newOwner);
    event ReceivedLessThanTarget(address indexed sender, uint amount);
    event Received(address indexed sender, uint amount);

    mapping(uint => GameSession[]) gameHistory; //historical data for each game 
    mapping(string => GameSession[]) userHistory; //historical data for each user
    mapping(address => Payment[]) payments; //historical data for each payment
    mapping(address => uint) totalPaid; //total amount paid by each sender

    // constructor(address _tokenAddressEGC, uint _oneEGC) {
    //     owner = msg.sender; // 'msg.sender' is sender of current call, contract deployer for a constructor
    //     emit OwnerSet(address(0), owner);
    //     tokenAddressEGC = _tokenAddressEGC;
    //     oneEGC = _oneEGC;
    // }

    constructor(address _address) {
        owner = msg.sender; // 'msg.sender' is sender of current call, contract deployer for a constructor
        emit OwnerSet(address(0), owner);
        tokenAddressEGC = _address;
        oneEGC = 2 gwei;
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

    /**
     * @dev Sends a specified amount of EGC tokens to a given user address.
     * @param _user The address of the user to send the tokens to.
     * @param _amount The amount of EGC tokens to send.
     * @notice Requires that the contract has enough tokens to send.
     */
    function sendEgc(address _user, uint256 _amount) private {
        ERC20 token = ERC20(tokenAddressEGC);
        token.transfer(_user, _amount);
    }

    /**
     * @dev Retrieves the total amount paid by a specific user.
     * @param _user The address of the user.
     * @return The total amount paid by the user.
     */
    function userTotal(address _user) public view returns (uint) {
        return totalPaid[_user];
    }

    /**
     * @dev Retrieves the stake history for a specific user.
     * @param _user The address of the user.
     * @return An array of Payment structs representing the stake history for the specified user.
     */
    function userStakeHistory(address _user) public view returns (Payment[] memory) {
        return payments[_user];
    }

    /**
     * @dev Buys EGC tokens using Ether.
     *      Sends the received Ether to the EGC token contract to buy the tokens.
     *      Keeps track of the payments and total amount paid by the sender.
     * @notice The amount of EGC tokens bought is equal to the amount of Ether sent.
     * @notice Requires that the amount of Ether sent is greater than the minimum amount (one EGC token).
     */
    function buyEgc() public payable {
        if (msg.value < oneEGC) {
            emit ReceivedLessThanTarget(msg.sender, msg.value);
            revert("Received amount is less than the target amount");
        }
        emit Received(msg.sender, msg.value);

        payments[msg.sender].push(Payment({
            sender: msg.sender,
            amount: msg.value,
            time: block.timestamp
        }));
        totalPaid[msg.sender] += msg.value;
        sendEgc(msg.sender, msg.value);
    }


    /**
     * @dev Fallback function to receive Ether.
     * This function is called when the contract receives Ether without a function being explicitly called.
     */
    receive() external payable {
        emit Received(msg.sender, msg.value);
    }
}