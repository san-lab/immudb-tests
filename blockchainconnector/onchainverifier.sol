// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract HashVerification {
    struct StateCheck {
        bytes32 submittedHash;
        bytes32 submittedPreimage;
        bool verified;
        uint256 blockNumber;
    }

    mapping(address => mapping(address => StateCheck[])) public stateChecks;

    modifier onlyOriginatorBank(address _originatorBank) {
        require(msg.sender == _originatorBank, "Only the originator bank can use this function");
        _;
    }

    modifier onlyRecipientBank(address _recipientBank) {
        require(msg.sender == _recipientBank, "Only the recipient bank can use this function");
        _;
    }

    modifier preimageNotSubmitted(
        address _originatorBank,
        address _recipientBank,
        uint256 _blockNumber
    ) {
        StateCheck[] storage senderStateChecks = stateChecks[_originatorBank][_recipientBank];

        require(senderStateChecks.length > 0, "No submissions found for the given parties");

        bool found = false;
        for (uint256 i = 0; i < senderStateChecks.length; i++) {
            if (
                senderStateChecks[i].blockNumber == _blockNumber &&
                senderStateChecks[i].verified
            ) {
                found = true;
                break;
            }
        }

        require(!found, "Preimage already submitted and verified for the given blockNumber");
        _;
    }


    constructor() {}

    function submitHash(
        address _originatorBank,
        address _recipientBank,
        bytes32 _hash
    ) external onlyOriginatorBank(_originatorBank) {
        uint256 currentBlock = block.number;

        stateChecks[_originatorBank][_recipientBank].push(
            StateCheck({
                submittedHash: _hash,
                submittedPreimage: bytes32(0),
                verified: false,
                blockNumber: currentBlock
            })
        );
    }

    function submitPreimage(
        address _originatorBank,
        address _recipientBank,
        bytes32 _preimage,
        uint256 _blockNumber
    ) external onlyRecipientBank(_recipientBank) preimageNotSubmitted(_originatorBank, _recipientBank, _blockNumber){

        StateCheck[] storage senderStateChecks = stateChecks[_originatorBank][_recipientBank];

        for (uint256 i = 0; i < senderStateChecks.length; i++) {
            if (!senderStateChecks[i].verified && senderStateChecks[i].blockNumber == _blockNumber) {
                senderStateChecks[i].submittedPreimage = _preimage;
                senderStateChecks[i].verified =
                    keccak256(abi.encodePacked(_preimage)) ==
                    senderStateChecks[i].submittedHash;
                break;
            }
        }
    }

    function getPendingSubmissions(
        address _originatorBank,
        address _recipientBank
    ) external view onlyRecipientBank(_recipientBank) returns (uint256[] memory) {

        StateCheck[] storage senderStateChecks = stateChecks[_originatorBank][_recipientBank];
        uint256[] memory pendingBlocks = new uint256[](senderStateChecks.length);

        uint256 pendingCount = 0;
        for (uint256 i = 0; i < senderStateChecks.length; i++) {
            if (!senderStateChecks[i].verified) {
                pendingBlocks[pendingCount] = senderStateChecks[i].blockNumber;
                pendingCount++;
            }
        }

        // Resize the array to the number of pending blocks
        assembly {
            mstore(pendingBlocks, pendingCount)
        }

        return pendingBlocks;
    }

    function getStateCheckByIndex(
        address _originatorBank,
        address _recipientBank,
        uint256 index
    )   external view returns (
            bytes32 submittedHash,
            bytes32 submittedPreimage,
            bool verified,
            uint256 blockNumber
        ) 
    {
        StateCheck storage stateCheck = stateChecks[_originatorBank][_recipientBank][index];
        return (
            stateCheck.submittedHash,
            stateCheck.submittedPreimage,
            stateCheck.verified,
            stateCheck.blockNumber
        );
    }

    function getStateCheckByBlockNumber(
        address _originatorBank,
        address _recipientBank,
        uint256 _blockNumber
    )   external view returns (
            bytes32 submittedHash,
            bytes32 submittedPreimage,
            bool verified,
            uint256 blockNumber
        )
    {
        StateCheck[] storage senderStateChecks = stateChecks[_originatorBank][_recipientBank];

        for (uint256 i = 0; i < senderStateChecks.length; i++) {
            if (senderStateChecks[i].blockNumber == _blockNumber) {
                return (
                    senderStateChecks[i].submittedHash,
                    senderStateChecks[i].submittedPreimage,
                    senderStateChecks[i].verified,
                    senderStateChecks[i].blockNumber
                );
            }
        }

        revert("No state check found for the given blockNumber");
    }
}

