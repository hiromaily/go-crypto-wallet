// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var address_pb = require('./address_pb.js');
var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js');

function serialize_google_protobuf_Empty(arg) {
  if (!(arg instanceof google_protobuf_empty_pb.Empty)) {
    throw new Error('Expected argument of type google.protobuf.Empty');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_google_protobuf_Empty(buffer_arg) {
  return google_protobuf_empty_pb.Empty.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_address_RequestIsValidAddress(arg) {
  if (!(arg instanceof address_pb.RequestIsValidAddress)) {
    throw new Error('Expected argument of type rippleapi.address.RequestIsValidAddress');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_address_RequestIsValidAddress(buffer_arg) {
  return address_pb.RequestIsValidAddress.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_address_ResponseGenerateAddress(arg) {
  if (!(arg instanceof address_pb.ResponseGenerateAddress)) {
    throw new Error('Expected argument of type rippleapi.address.ResponseGenerateAddress');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_address_ResponseGenerateAddress(buffer_arg) {
  return address_pb.ResponseGenerateAddress.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_address_ResponseGenerateXAddress(arg) {
  if (!(arg instanceof address_pb.ResponseGenerateXAddress)) {
    throw new Error('Expected argument of type rippleapi.address.ResponseGenerateXAddress');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_address_ResponseGenerateXAddress(buffer_arg) {
  return address_pb.ResponseGenerateXAddress.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_address_ResponseIsValidAddress(arg) {
  if (!(arg instanceof address_pb.ResponseIsValidAddress)) {
    throw new Error('Expected argument of type rippleapi.address.ResponseIsValidAddress');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_address_ResponseIsValidAddress(buffer_arg) {
  return address_pb.ResponseIsValidAddress.deserializeBinary(new Uint8Array(buffer_arg));
}


// RippleAddressAPI 
var RippleAddressAPIService = exports.RippleAddressAPIService = {
  // https://xrpl.org/rippleapi-reference.html#generateaddress
generateAddress: {
    path: '/rippleapi.address.RippleAddressAPI/GenerateAddress',
    requestStream: false,
    responseStream: false,
    requestType: google_protobuf_empty_pb.Empty,
    responseType: address_pb.ResponseGenerateAddress,
    requestSerialize: serialize_google_protobuf_Empty,
    requestDeserialize: deserialize_google_protobuf_Empty,
    responseSerialize: serialize_rippleapi_address_ResponseGenerateAddress,
    responseDeserialize: deserialize_rippleapi_address_ResponseGenerateAddress,
  },
  // https://xrpl.org/rippleapi-reference.html#generatexaddress
generateXAddress: {
    path: '/rippleapi.address.RippleAddressAPI/GenerateXAddress',
    requestStream: false,
    responseStream: false,
    requestType: google_protobuf_empty_pb.Empty,
    responseType: address_pb.ResponseGenerateXAddress,
    requestSerialize: serialize_google_protobuf_Empty,
    requestDeserialize: deserialize_google_protobuf_Empty,
    responseSerialize: serialize_rippleapi_address_ResponseGenerateXAddress,
    responseDeserialize: deserialize_rippleapi_address_ResponseGenerateXAddress,
  },
  // https://xrpl.org/rippleapi-reference.html#isvalidaddress
isValidAddress: {
    path: '/rippleapi.address.RippleAddressAPI/IsValidAddress',
    requestStream: false,
    responseStream: false,
    requestType: address_pb.RequestIsValidAddress,
    responseType: address_pb.ResponseIsValidAddress,
    requestSerialize: serialize_rippleapi_address_RequestIsValidAddress,
    requestDeserialize: deserialize_rippleapi_address_RequestIsValidAddress,
    responseSerialize: serialize_rippleapi_address_ResponseIsValidAddress,
    responseDeserialize: deserialize_rippleapi_address_ResponseIsValidAddress,
  },
};

exports.RippleAddressAPIClient = grpc.makeGenericClientConstructor(RippleAddressAPIService);
