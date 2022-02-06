// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var transaction_pb = require('./transaction_pb.js');
var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js');
var gogo_protobuf_gogoproto_gogo_pb = require('./gogo/protobuf/gogoproto/gogo_pb.js');

function serialize_google_protobuf_Empty(arg) {
  if (!(arg instanceof google_protobuf_empty_pb.Empty)) {
    throw new Error('Expected argument of type google.protobuf.Empty');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_google_protobuf_Empty(buffer_arg) {
  return google_protobuf_empty_pb.Empty.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_transaction_RequestCombineTransaction(arg) {
  if (!(arg instanceof transaction_pb.RequestCombineTransaction)) {
    throw new Error('Expected argument of type rippleapi.transaction.RequestCombineTransaction');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_transaction_RequestCombineTransaction(buffer_arg) {
  return transaction_pb.RequestCombineTransaction.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_transaction_RequestGetTransaction(arg) {
  if (!(arg instanceof transaction_pb.RequestGetTransaction)) {
    throw new Error('Expected argument of type rippleapi.transaction.RequestGetTransaction');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_transaction_RequestGetTransaction(buffer_arg) {
  return transaction_pb.RequestGetTransaction.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_transaction_RequestPrepareTransaction(arg) {
  if (!(arg instanceof transaction_pb.RequestPrepareTransaction)) {
    throw new Error('Expected argument of type rippleapi.transaction.RequestPrepareTransaction');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_transaction_RequestPrepareTransaction(buffer_arg) {
  return transaction_pb.RequestPrepareTransaction.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_transaction_RequestSignTransaction(arg) {
  if (!(arg instanceof transaction_pb.RequestSignTransaction)) {
    throw new Error('Expected argument of type rippleapi.transaction.RequestSignTransaction');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_transaction_RequestSignTransaction(buffer_arg) {
  return transaction_pb.RequestSignTransaction.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_transaction_RequestSubmitTransaction(arg) {
  if (!(arg instanceof transaction_pb.RequestSubmitTransaction)) {
    throw new Error('Expected argument of type rippleapi.transaction.RequestSubmitTransaction');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_transaction_RequestSubmitTransaction(buffer_arg) {
  return transaction_pb.RequestSubmitTransaction.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_transaction_ResponseCombineTransaction(arg) {
  if (!(arg instanceof transaction_pb.ResponseCombineTransaction)) {
    throw new Error('Expected argument of type rippleapi.transaction.ResponseCombineTransaction');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_transaction_ResponseCombineTransaction(buffer_arg) {
  return transaction_pb.ResponseCombineTransaction.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_transaction_ResponseGetTransaction(arg) {
  if (!(arg instanceof transaction_pb.ResponseGetTransaction)) {
    throw new Error('Expected argument of type rippleapi.transaction.ResponseGetTransaction');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_transaction_ResponseGetTransaction(buffer_arg) {
  return transaction_pb.ResponseGetTransaction.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_transaction_ResponsePrepareTransaction(arg) {
  if (!(arg instanceof transaction_pb.ResponsePrepareTransaction)) {
    throw new Error('Expected argument of type rippleapi.transaction.ResponsePrepareTransaction');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_transaction_ResponsePrepareTransaction(buffer_arg) {
  return transaction_pb.ResponsePrepareTransaction.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_transaction_ResponseSignTransaction(arg) {
  if (!(arg instanceof transaction_pb.ResponseSignTransaction)) {
    throw new Error('Expected argument of type rippleapi.transaction.ResponseSignTransaction');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_transaction_ResponseSignTransaction(buffer_arg) {
  return transaction_pb.ResponseSignTransaction.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_transaction_ResponseSubmitTransaction(arg) {
  if (!(arg instanceof transaction_pb.ResponseSubmitTransaction)) {
    throw new Error('Expected argument of type rippleapi.transaction.ResponseSubmitTransaction');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_transaction_ResponseSubmitTransaction(buffer_arg) {
  return transaction_pb.ResponseSubmitTransaction.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_rippleapi_transaction_ResponseWaitValidation(arg) {
  if (!(arg instanceof transaction_pb.ResponseWaitValidation)) {
    throw new Error('Expected argument of type rippleapi.transaction.ResponseWaitValidation');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_rippleapi_transaction_ResponseWaitValidation(buffer_arg) {
  return transaction_pb.ResponseWaitValidation.deserializeBinary(new Uint8Array(buffer_arg));
}


// RippleTransactionAPI 
var RippleTransactionAPIService = exports.RippleTransactionAPIService = {
  // https://xrpl.org/rippleapi-reference.html#preparetransaction
prepareTransaction: {
    path: '/rippleapi.transaction.RippleTransactionAPI/PrepareTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.RequestPrepareTransaction,
    responseType: transaction_pb.ResponsePrepareTransaction,
    requestSerialize: serialize_rippleapi_transaction_RequestPrepareTransaction,
    requestDeserialize: deserialize_rippleapi_transaction_RequestPrepareTransaction,
    responseSerialize: serialize_rippleapi_transaction_ResponsePrepareTransaction,
    responseDeserialize: deserialize_rippleapi_transaction_ResponsePrepareTransaction,
  },
  signTransaction: {
    path: '/rippleapi.transaction.RippleTransactionAPI/SignTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.RequestSignTransaction,
    responseType: transaction_pb.ResponseSignTransaction,
    requestSerialize: serialize_rippleapi_transaction_RequestSignTransaction,
    requestDeserialize: deserialize_rippleapi_transaction_RequestSignTransaction,
    responseSerialize: serialize_rippleapi_transaction_ResponseSignTransaction,
    responseDeserialize: deserialize_rippleapi_transaction_ResponseSignTransaction,
  },
  submitTransaction: {
    path: '/rippleapi.transaction.RippleTransactionAPI/SubmitTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.RequestSubmitTransaction,
    responseType: transaction_pb.ResponseSubmitTransaction,
    requestSerialize: serialize_rippleapi_transaction_RequestSubmitTransaction,
    requestDeserialize: deserialize_rippleapi_transaction_RequestSubmitTransaction,
    responseSerialize: serialize_rippleapi_transaction_ResponseSubmitTransaction,
    responseDeserialize: deserialize_rippleapi_transaction_ResponseSubmitTransaction,
  },
  waitValidation: {
    path: '/rippleapi.transaction.RippleTransactionAPI/WaitValidation',
    requestStream: false,
    responseStream: true,
    requestType: google_protobuf_empty_pb.Empty,
    responseType: transaction_pb.ResponseWaitValidation,
    requestSerialize: serialize_google_protobuf_Empty,
    requestDeserialize: deserialize_google_protobuf_Empty,
    responseSerialize: serialize_rippleapi_transaction_ResponseWaitValidation,
    responseDeserialize: deserialize_rippleapi_transaction_ResponseWaitValidation,
  },
  getTransaction: {
    path: '/rippleapi.transaction.RippleTransactionAPI/GetTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.RequestGetTransaction,
    responseType: transaction_pb.ResponseGetTransaction,
    requestSerialize: serialize_rippleapi_transaction_RequestGetTransaction,
    requestDeserialize: deserialize_rippleapi_transaction_RequestGetTransaction,
    responseSerialize: serialize_rippleapi_transaction_ResponseGetTransaction,
    responseDeserialize: deserialize_rippleapi_transaction_ResponseGetTransaction,
  },
  combineTransaction: {
    path: '/rippleapi.transaction.RippleTransactionAPI/CombineTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.RequestCombineTransaction,
    responseType: transaction_pb.ResponseCombineTransaction,
    requestSerialize: serialize_rippleapi_transaction_RequestCombineTransaction,
    requestDeserialize: deserialize_rippleapi_transaction_RequestCombineTransaction,
    responseSerialize: serialize_rippleapi_transaction_ResponseCombineTransaction,
    responseDeserialize: deserialize_rippleapi_transaction_ResponseCombineTransaction,
  },
};

exports.RippleTransactionAPIClient = grpc.makeGenericClientConstructor(RippleTransactionAPIService);
