package abis

var SFC = `[
  {
    "inputs": [],
    "name": "currentEpoch",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "_epoch",
        "type": "uint256"
      }
    ],
    "name": "getEpochSnapshot",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "endTime",
        "type": "uint256"
      },
      {
        "internalType": "uint256",
        "name": "epochFee",
        "type": "uint256"
      },
      {
        "internalType": "uint256",
        "name": "totalBaseRewardWeight",
        "type": "uint256"
      },
      {
        "internalType": "uint256",
        "name": "totalTxRewardWeight",
        "type": "uint256"
      },
      {
        "internalType": "uint256",
        "name": "baseRewardPerSecond",
        "type": "uint256"
      },
      {
        "internalType": "uint256",
        "name": "totalStake",
        "type": "uint256"
      },
      {
        "internalType": "uint256",
        "name": "totalSupply",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  }
]`
