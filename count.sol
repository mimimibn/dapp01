// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

contract count {

    uint256 public num;

    function addOne() public returns (uint256){
        num++;
        return num;
    }
    function getCount() public view returns (uint256) {
        return num;
    }
}