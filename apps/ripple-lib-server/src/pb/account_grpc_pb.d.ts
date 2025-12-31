// package: rippleapi.account
// file: account.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "grpc";
import * as account_pb from "./account_pb";

interface IRippleAccountAPIService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    getAccountInfo: IRippleAccountAPIService_IGetAccountInfo;
}

interface IRippleAccountAPIService_IGetAccountInfo extends grpc.MethodDefinition<account_pb.RequestGetAccountInfo, account_pb.ResponseGetAccountInfo> {
    path: "/rippleapi.account.RippleAccountAPI/GetAccountInfo";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<account_pb.RequestGetAccountInfo>;
    requestDeserialize: grpc.deserialize<account_pb.RequestGetAccountInfo>;
    responseSerialize: grpc.serialize<account_pb.ResponseGetAccountInfo>;
    responseDeserialize: grpc.deserialize<account_pb.ResponseGetAccountInfo>;
}

export const RippleAccountAPIService: IRippleAccountAPIService;

export interface IRippleAccountAPIServer {
    getAccountInfo: grpc.handleUnaryCall<account_pb.RequestGetAccountInfo, account_pb.ResponseGetAccountInfo>;
}

export interface IRippleAccountAPIClient {
    getAccountInfo(request: account_pb.RequestGetAccountInfo, callback: (error: grpc.ServiceError | null, response: account_pb.ResponseGetAccountInfo) => void): grpc.ClientUnaryCall;
    getAccountInfo(request: account_pb.RequestGetAccountInfo, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: account_pb.ResponseGetAccountInfo) => void): grpc.ClientUnaryCall;
    getAccountInfo(request: account_pb.RequestGetAccountInfo, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: account_pb.ResponseGetAccountInfo) => void): grpc.ClientUnaryCall;
}

export class RippleAccountAPIClient extends grpc.Client implements IRippleAccountAPIClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public getAccountInfo(request: account_pb.RequestGetAccountInfo, callback: (error: grpc.ServiceError | null, response: account_pb.ResponseGetAccountInfo) => void): grpc.ClientUnaryCall;
    public getAccountInfo(request: account_pb.RequestGetAccountInfo, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: account_pb.ResponseGetAccountInfo) => void): grpc.ClientUnaryCall;
    public getAccountInfo(request: account_pb.RequestGetAccountInfo, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: account_pb.ResponseGetAccountInfo) => void): grpc.ClientUnaryCall;
}
