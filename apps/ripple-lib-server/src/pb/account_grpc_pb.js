// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var account_pb = require('./account_pb.js');

function serialize_rippleapi_account_RequestGetAccountInfo(arg) {
  if (!(arg instanceof account_pb.RequestGetAccountInfo)) {
    throw new Error('Expected argument of type rippleapi.account.RequestGetAccountInfo');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_account_RequestGetAccountInfo(buffer_arg) {
  return account_pb.RequestGetAccountInfo.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_account_ResponseGetAccountInfo(arg) {
  if (!(arg instanceof account_pb.ResponseGetAccountInfo)) {
    throw new Error('Expected argument of type rippleapi.account.ResponseGetAccountInfo');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_account_ResponseGetAccountInfo(buffer_arg) {
  return account_pb.ResponseGetAccountInfo.deserializeBinary(new Uint8Array(buffer_arg));
}


// RippleAccountAPI 
var RippleAccountAPIService = exports.RippleAccountAPIService = {
  // https://xrpl.org/rippleapi-reference.html#getaccountinfo
getAccountInfo: {
    path: '/rippleapi.account.RippleAccountAPI/GetAccountInfo',
    requestStream: false,
    responseStream: false,
    requestType: account_pb.RequestGetAccountInfo,
    responseType: account_pb.ResponseGetAccountInfo,
    requestSerialize: serialize_rippleapi_account_RequestGetAccountInfo,
    requestDeserialize: deserialize_rippleapi_account_RequestGetAccountInfo,
    responseSerialize: serialize_rippleapi_account_ResponseGetAccountInfo,
    responseDeserialize: deserialize_rippleapi_account_ResponseGetAccountInfo,
  },
};

exports.RippleAccountAPIClient = grpc.makeGenericClientConstructor(RippleAccountAPIService);
