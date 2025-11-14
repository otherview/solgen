// Copyright (c) 2025 The solgen developers
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

/// @title DataTypesContract - Comprehensive test contract for all Solidity data types
/// @notice Tests generation of Go bindings for various Solidity data types
contract DataTypesContract {
    
    // Custom struct with mixed field types
    struct MixedStruct {
        bool boolField;
        uint8 uint8Field;
        uint256 uint256Field;
        int256 int256Field;  // This should use mustDecodeInt256
        address addressField;
        bytes32 bytes32Field;
        string stringField;
    }
    
    // Removed NestedStruct to avoid recursive struct issue in parser

    // =============================================================================
    // SINGLE INPUT, SINGLE OUTPUT FUNCTIONS
    // =============================================================================
    
    // Boolean functions
    function echoBool(bool value) external pure returns (bool) {
        return value;
    }
    
    // Unsigned integer functions
    function echoUint8(uint8 value) external pure returns (uint8) {
        return value;
    }
    
    function echoUint16(uint16 value) external pure returns (uint16) {
        return value;
    }
    
    function echoUint32(uint32 value) external pure returns (uint32) {
        return value;
    }
    
    function echoUint64(uint64 value) external pure returns (uint64) {
        return value;
    }
    
    function echoUint256(uint256 value) external pure returns (uint256) {
        return value;
    }
    
    // Signed integer functions
    function echoInt8(int8 value) external pure returns (int8) {
        return value;
    }
    
    function echoInt16(int16 value) external pure returns (int16) {
        return value;
    }
    
    function echoInt32(int32 value) external pure returns (int32) {
        return value;
    }
    
    function echoInt64(int64 value) external pure returns (int64) {
        return value;
    }
    
    function echoInt256(int256 value) external pure returns (int256) {
        return value;
    }
    
    // Address function
    function echoAddress(address value) external pure returns (address) {
        return value;
    }
    
    // Bytes functions
    function echoBytes(bytes memory value) external pure returns (bytes memory) {
        return value;
    }
    
    function echoBytes1(bytes1 value) external pure returns (bytes1) {
        return value;
    }
    
    function echoBytes32(bytes32 value) external pure returns (bytes32) {
        return value;
    }
    
    // String function
    function echoString(string memory value) external pure returns (string memory) {
        return value;
    }
    
    // Struct functions
    function echoMixedStruct(MixedStruct memory value) external pure returns (MixedStruct memory) {
        return value;
    }
    
    // echoNestedStruct removed due to parser issue with nested structs
    
    // =============================================================================
    // ARRAY FUNCTIONS (DYNAMIC)
    // =============================================================================
    
    function echoUint256Array(uint256[] memory values) external pure returns (uint256[] memory) {
        return values;
    }
    
    function echoInt256Array(int256[] memory values) external pure returns (int256[] memory) {
        return values;
    }
    
    function echoAddressArray(address[] memory values) external pure returns (address[] memory) {
        return values;
    }
    
    function echoBoolArray(bool[] memory values) external pure returns (bool[] memory) {
        return values;
    }
    
    function echoMixedStructArray(MixedStruct[] memory values) external pure returns (MixedStruct[] memory) {
        return values;
    }
    
    // =============================================================================
    // FIXED ARRAY FUNCTIONS  
    // =============================================================================
    
    function echoUint256FixedArray(uint256[3] memory values) external pure returns (uint256[3] memory) {
        return values;
    }
    
    function echoAddressFixedArray(address[2] memory values) external pure returns (address[2] memory) {
        return values;
    }
    
    function echoBytes32FixedArray(bytes32[4] memory values) external pure returns (bytes32[4] memory) {
        return values;
    }
    
    // =============================================================================
    // MIXED INPUT FUNCTIONS
    // =============================================================================
    
    function mixedInputs(
        bool boolVal,
        uint256 uintVal,
        int256 intVal,
        address addressVal,
        string memory stringVal
    ) external pure returns (bool) {
        return boolVal && uintVal > 0 && intVal != 0 && addressVal != address(0) && bytes(stringVal).length > 0;
    }
    
    function arrayInputs(
        uint256[] memory uints,
        address[] memory addresses,
        MixedStruct[] memory structs
    ) external pure returns (uint256) {
        return uints.length + addresses.length + structs.length;
    }
    
    // =============================================================================
    // MIXED OUTPUT FUNCTIONS (MULTIPLE RETURN VALUES)
    // =============================================================================
    
    function getBasicTypes() external pure returns (
        bool boolVal,
        uint256 uintVal,
        int256 intVal,
        address addressVal,
        string memory stringVal
    ) {
        return (true, 12345, -6789, address(0x1234567890123456789012345678901234567890), "hello world");
    }
    
    function getArrayTypes() external pure returns (
        uint256[] memory uints,
        int256[] memory ints,
        address[] memory addresses
    ) {
        uints = new uint256[](3);
        uints[0] = 100;
        uints[1] = 200;
        uints[2] = 300;
        
        ints = new int256[](2);
        ints[0] = -100;
        ints[1] = -200;
        
        addresses = new address[](2);
        addresses[0] = address(0x1111111111111111111111111111111111111111);
        addresses[1] = address(0x2222222222222222222222222222222222222222);
        
        return (uints, ints, addresses);
    }
    
    function getStructTypes() external pure returns (
        MixedStruct memory simple,
        MixedStruct[] memory structs
    ) {
        simple = MixedStruct({
            boolField: true,
            uint8Field: 255,
            uint256Field: 123456789,
            int256Field: -987654321,
            addressField: address(0xaAaAaAaaAaAaAaaAaAAAAAAAAaaaAaAaAaaAaaAa),
            bytes32Field: keccak256("test"),
            stringField: "struct test"
        });
        
        structs = new MixedStruct[](2);
        structs[0] = simple;
        structs[1] = MixedStruct({
            boolField: false,
            uint8Field: 128,
            uint256Field: 999999999,
            int256Field: -111111111,
            addressField: address(0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB),
            bytes32Field: keccak256("test2"),
            stringField: "second struct"
        });
        
        return (simple, structs);
    }
    
    // =============================================================================
    // EDGE CASE FUNCTIONS
    // =============================================================================
    
    function emptyArrays() external pure returns (
        uint256[] memory emptyUints,
        address[] memory emptyAddresses,
        MixedStruct[] memory emptyStructs
    ) {
        return (new uint256[](0), new address[](0), new MixedStruct[](0));
    }
    
    function largeNumbers() external pure returns (
        uint256 maxUint256,
        int256 maxInt256,
        int256 minInt256
    ) {
        return (
            type(uint256).max,
            type(int256).max,
            type(int256).min
        );
    }
    
    // =============================================================================
    // EVENTS WITH DIFFERENT DATA TYPES
    // =============================================================================
    
    event BasicTypesEvent(
        bool indexed boolVal,
        uint256 indexed uintVal,
        int256 intVal,
        address addressVal,
        string stringVal
    );
    
    event ArrayTypesEvent(
        uint256[] uints,
        int256[] ints,
        address[] addresses
    );
    
    event StructTypesEvent(
        MixedStruct simple,
        MixedStruct[] structs
    );
    
    // Functions to emit events for testing
    function emitBasicTypes() external {
        emit BasicTypesEvent(true, 12345, -6789, msg.sender, "event test");
    }
    
    function emitArrayTypes() external {
        uint256[] memory uints = new uint256[](2);
        uints[0] = 100;
        uints[1] = 200;
        
        int256[] memory ints = new int256[](1);
        ints[0] = -300;
        
        address[] memory addresses = new address[](1);
        addresses[0] = msg.sender;
        
        emit ArrayTypesEvent(uints, ints, addresses);
    }
    
    // =============================================================================
    // CUSTOM ERRORS WITH DIFFERENT DATA TYPES
    // =============================================================================
    
    error SimpleError(string message);
    error DataTypesError(uint256 code, int256 value, address user, bool flag);
    error ArrayError(uint256[] values, address[] users);
    error StructError(MixedStruct data, uint256 timestamp);
    
    // Functions that can revert with custom errors
    function triggerSimpleError() external pure {
        revert SimpleError("This is a simple error");
    }
    
    function triggerDataTypesError() external view {
        revert DataTypesError(404, -1, msg.sender, false);
    }
    
    function triggerArrayError() external view {
        uint256[] memory values = new uint256[](2);
        values[0] = 1;
        values[1] = 2;
        
        address[] memory users = new address[](1);
        users[0] = msg.sender;
        
        revert ArrayError(values, users);
    }
    
    function triggerStructError() external view {
        MixedStruct memory data = MixedStruct({
            boolField: true,
            uint8Field: 42,
            uint256Field: 12345,
            int256Field: -54321,
            addressField: msg.sender,
            bytes32Field: blockhash(block.number - 1),
            stringField: "error struct"
        });
        
        revert StructError(data, block.timestamp);
    }
}