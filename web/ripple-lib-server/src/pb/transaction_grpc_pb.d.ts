// package: rippleapi.transaction
// file: transaction.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "grpc";
import * as transaction_pb from "./transaction_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as gogo_protobuf_gogoproto_gogo_pb from "./gogo/protobuf/gogoproto/gogo_pb";

interface IRippleTransactionAPIService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    prepareTransaction: IRippleTransactionAPIService_IPrepareTransaction;
    signTransaction: IRippleTransactionAPIService_ISignTransaction;
    submitTransaction: IRippleTransactionAPIService_ISubmitTransaction;
    waitValidation: IRippleTransactionAPIService_IWaitValidation;
    getTransaction: IRippleTransactionAPIService_IGetTransaction;
    combineTransaction: IRippleTransactionAPIService_ICombineTransaction;
}

interface IRippleTransactionAPIService_IPrepareTransaction extends grpc.MethodDefinition<transaction_pb.RequestPrepareTransaction, transaction_pb.ResponsePrepareTransaction> {
    path: "/rippleapi.transaction.RippleTransactionAPI/PrepareTransaction";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<transaction_pb.RequestPrepareTransaction>;
    requestDeserialize: grpc.deserialize<transaction_pb.RequestPrepareTransaction>;
    responseSerialize: grpc.serialize<transaction_pb.ResponsePrepareTransaction>;
    responseDeserialize: grpc.deserialize<transaction_pb.ResponsePrepareTransaction>;
}
interface IRippleTransactionAPIService_ISignTransaction extends grpc.MethodDefinition<transaction_pb.RequestSignTransaction, transaction_pb.ResponseSignTransaction> {
    path: "/rippleapi.transaction.RippleTransactionAPI/SignTransaction";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<transaction_pb.RequestSignTransaction>;
    requestDeserialize: grpc.deserialize<transaction_pb.RequestSignTransaction>;
    responseSerialize: grpc.serialize<transaction_pb.ResponseSignTransaction>;
    responseDeserialize: grpc.deserialize<transaction_pb.ResponseSignTransaction>;
}
interface IRippleTransactionAPIService_ISubmitTransaction extends grpc.MethodDefinition<transaction_pb.RequestSubmitTransaction, transaction_pb.ResponseSubmitTransaction> {
    path: "/rippleapi.transaction.RippleTransactionAPI/SubmitTransaction";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<transaction_pb.RequestSubmitTransaction>;
    requestDeserialize: grpc.deserialize<transaction_pb.RequestSubmitTransaction>;
    responseSerialize: grpc.serialize<transaction_pb.ResponseSubmitTransaction>;
    responseDeserialize: grpc.deserialize<transaction_pb.ResponseSubmitTransaction>;
}
interface IRippleTransactionAPIService_IWaitValidation extends grpc.MethodDefinition<google_protobuf_empty_pb.Empty, transaction_pb.ResponseWaitValidation> {
    path: "/rippleapi.transaction.RippleTransactionAPI/WaitValidation";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    requestDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
    responseSerialize: grpc.serialize<transaction_pb.ResponseWaitValidation>;
    responseDeserialize: grpc.deserialize<transaction_pb.ResponseWaitValidation>;
}
interface IRippleTransactionAPIService_IGetTransaction extends grpc.MethodDefinition<transaction_pb.RequestGetTransaction, transaction_pb.ResponseGetTransaction> {
    path: "/rippleapi.transaction.RippleTransactionAPI/GetTransaction";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<transaction_pb.RequestGetTransaction>;
    requestDeserialize: grpc.deserialize<transaction_pb.RequestGetTransaction>;
    responseSerialize: grpc.serialize<transaction_pb.ResponseGetTransaction>;
    responseDeserialize: grpc.deserialize<transaction_pb.ResponseGetTransaction>;
}
interface IRippleTransactionAPIService_ICombineTransaction extends grpc.MethodDefinition<transaction_pb.RequestCombineTransaction, transaction_pb.ResponseCombineTransaction> {
    path: "/rippleapi.transaction.RippleTransactionAPI/CombineTransaction";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<transaction_pb.RequestCombineTransaction>;
    requestDeserialize: grpc.deserialize<transaction_pb.RequestCombineTransaction>;
    responseSerialize: grpc.serialize<transaction_pb.ResponseCombineTransaction>;
    responseDeserialize: grpc.deserialize<transaction_pb.ResponseCombineTransaction>;
}

export const RippleTransactionAPIService: IRippleTransactionAPIService;

export interface IRippleTransactionAPIServer {
    prepareTransaction: grpc.handleUnaryCall<transaction_pb.RequestPrepareTransaction, transaction_pb.ResponsePrepareTransaction>;
    signTransaction: grpc.handleUnaryCall<transaction_pb.RequestSignTransaction, transaction_pb.ResponseSignTransaction>;
    submitTransaction: grpc.handleUnaryCall<transaction_pb.RequestSubmitTransaction, transaction_pb.ResponseSubmitTransaction>;
    waitValidation: grpc.handleServerStreamingCall<google_protobuf_empty_pb.Empty, transaction_pb.ResponseWaitValidation>;
    getTransaction: grpc.handleUnaryCall<transaction_pb.RequestGetTransaction, transaction_pb.ResponseGetTransaction>;
    combineTransaction: grpc.handleUnaryCall<transaction_pb.RequestCombineTransaction, transaction_pb.ResponseCombineTransaction>;
}

export interface IRippleTransactionAPIClient {
    prepareTransaction(request: transaction_pb.RequestPrepareTransaction, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponsePrepareTransaction) => void): grpc.ClientUnaryCall;
    prepareTransaction(request: transaction_pb.RequestPrepareTransaction, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponsePrepareTransaction) => void): grpc.ClientUnaryCall;
    prepareTransaction(request: transaction_pb.RequestPrepareTransaction, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponsePrepareTransaction) => void): grpc.ClientUnaryCall;
    signTransaction(request: transaction_pb.RequestSignTransaction, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseSignTransaction) => void): grpc.ClientUnaryCall;
    signTransaction(request: transaction_pb.RequestSignTransaction, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseSignTransaction) => void): grpc.ClientUnaryCall;
    signTransaction(request: transaction_pb.RequestSignTransaction, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseSignTransaction) => void): grpc.ClientUnaryCall;
    submitTransaction(request: transaction_pb.RequestSubmitTransaction, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseSubmitTransaction) => void): grpc.ClientUnaryCall;
    submitTransaction(request: transaction_pb.RequestSubmitTransaction, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseSubmitTransaction) => void): grpc.ClientUnaryCall;
    submitTransaction(request: transaction_pb.RequestSubmitTransaction, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseSubmitTransaction) => void): grpc.ClientUnaryCall;
    waitValidation(request: google_protobuf_empty_pb.Empty, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<transaction_pb.ResponseWaitValidation>;
    waitValidation(request: google_protobuf_empty_pb.Empty, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<transaction_pb.ResponseWaitValidation>;
    getTransaction(request: transaction_pb.RequestGetTransaction, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseGetTransaction) => void): grpc.ClientUnaryCall;
    getTransaction(request: transaction_pb.RequestGetTransaction, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseGetTransaction) => void): grpc.ClientUnaryCall;
    getTransaction(request: transaction_pb.RequestGetTransaction, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseGetTransaction) => void): grpc.ClientUnaryCall;
    combineTransaction(request: transaction_pb.RequestCombineTransaction, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseCombineTransaction) => void): grpc.ClientUnaryCall;
    combineTransaction(request: transaction_pb.RequestCombineTransaction, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseCombineTransaction) => void): grpc.ClientUnaryCall;
    combineTransaction(request: transaction_pb.RequestCombineTransaction, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseCombineTransaction) => void): grpc.ClientUnaryCall;
}

export class RippleTransactionAPIClient extends grpc.Client implements IRippleTransactionAPIClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public prepareTransaction(request: transaction_pb.RequestPrepareTransaction, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponsePrepareTransaction) => void): grpc.ClientUnaryCall;
    public prepareTransaction(request: transaction_pb.RequestPrepareTransaction, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponsePrepareTransaction) => void): grpc.ClientUnaryCall;
    public prepareTransaction(request: transaction_pb.RequestPrepareTransaction, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponsePrepareTransaction) => void): grpc.ClientUnaryCall;
    public signTransaction(request: transaction_pb.RequestSignTransaction, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseSignTransaction) => void): grpc.ClientUnaryCall;
    public signTransaction(request: transaction_pb.RequestSignTransaction, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseSignTransaction) => void): grpc.ClientUnaryCall;
    public signTransaction(request: transaction_pb.RequestSignTransaction, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseSignTransaction) => void): grpc.ClientUnaryCall;
    public submitTransaction(request: transaction_pb.RequestSubmitTransaction, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseSubmitTransaction) => void): grpc.ClientUnaryCall;
    public submitTransaction(request: transaction_pb.RequestSubmitTransaction, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseSubmitTransaction) => void): grpc.ClientUnaryCall;
    public submitTransaction(request: transaction_pb.RequestSubmitTransaction, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseSubmitTransaction) => void): grpc.ClientUnaryCall;
    public waitValidation(request: google_protobuf_empty_pb.Empty, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<transaction_pb.ResponseWaitValidation>;
    public waitValidation(request: google_protobuf_empty_pb.Empty, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<transaction_pb.ResponseWaitValidation>;
    public getTransaction(request: transaction_pb.RequestGetTransaction, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseGetTransaction) => void): grpc.ClientUnaryCall;
    public getTransaction(request: transaction_pb.RequestGetTransaction, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseGetTransaction) => void): grpc.ClientUnaryCall;
    public getTransaction(request: transaction_pb.RequestGetTransaction, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseGetTransaction) => void): grpc.ClientUnaryCall;
    public combineTransaction(request: transaction_pb.RequestCombineTransaction, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseCombineTransaction) => void): grpc.ClientUnaryCall;
    public combineTransaction(request: transaction_pb.RequestCombineTransaction, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseCombineTransaction) => void): grpc.ClientUnaryCall;
    public combineTransaction(request: transaction_pb.RequestCombineTransaction, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: transaction_pb.ResponseCombineTransaction) => void): grpc.ClientUnaryCall;
}
