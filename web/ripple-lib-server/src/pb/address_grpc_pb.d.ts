// package: rippleapi.address
// file: address.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "grpc";
import * as address_pb from "./address_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";

interface IRippleAddressAPIService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    generateAddress: IRippleAddressAPIService_IGenerateAddress;
    generateXAddress: IRippleAddressAPIService_IGenerateXAddress;
    isValidAddress: IRippleAddressAPIService_IIsValidAddress;
}

interface IRippleAddressAPIService_IGenerateAddress extends grpc.MethodDefinition<google_protobuf_empty_pb.Empty, address_pb.ResponseGenerateAddress> {
    path: "/rippleapi.address.RippleAddressAPI/GenerateAddress";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    requestDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
    responseSerialize: grpc.serialize<address_pb.ResponseGenerateAddress>;
    responseDeserialize: grpc.deserialize<address_pb.ResponseGenerateAddress>;
}
interface IRippleAddressAPIService_IGenerateXAddress extends grpc.MethodDefinition<google_protobuf_empty_pb.Empty, address_pb.ResponseGenerateXAddress> {
    path: "/rippleapi.address.RippleAddressAPI/GenerateXAddress";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    requestDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
    responseSerialize: grpc.serialize<address_pb.ResponseGenerateXAddress>;
    responseDeserialize: grpc.deserialize<address_pb.ResponseGenerateXAddress>;
}
interface IRippleAddressAPIService_IIsValidAddress extends grpc.MethodDefinition<address_pb.RequestIsValidAddress, address_pb.ResponseIsValidAddress> {
    path: "/rippleapi.address.RippleAddressAPI/IsValidAddress";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<address_pb.RequestIsValidAddress>;
    requestDeserialize: grpc.deserialize<address_pb.RequestIsValidAddress>;
    responseSerialize: grpc.serialize<address_pb.ResponseIsValidAddress>;
    responseDeserialize: grpc.deserialize<address_pb.ResponseIsValidAddress>;
}

export const RippleAddressAPIService: IRippleAddressAPIService;

export interface IRippleAddressAPIServer {
    generateAddress: grpc.handleUnaryCall<google_protobuf_empty_pb.Empty, address_pb.ResponseGenerateAddress>;
    generateXAddress: grpc.handleUnaryCall<google_protobuf_empty_pb.Empty, address_pb.ResponseGenerateXAddress>;
    isValidAddress: grpc.handleUnaryCall<address_pb.RequestIsValidAddress, address_pb.ResponseIsValidAddress>;
}

export interface IRippleAddressAPIClient {
    generateAddress(request: google_protobuf_empty_pb.Empty, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseGenerateAddress) => void): grpc.ClientUnaryCall;
    generateAddress(request: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseGenerateAddress) => void): grpc.ClientUnaryCall;
    generateAddress(request: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseGenerateAddress) => void): grpc.ClientUnaryCall;
    generateXAddress(request: google_protobuf_empty_pb.Empty, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseGenerateXAddress) => void): grpc.ClientUnaryCall;
    generateXAddress(request: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseGenerateXAddress) => void): grpc.ClientUnaryCall;
    generateXAddress(request: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseGenerateXAddress) => void): grpc.ClientUnaryCall;
    isValidAddress(request: address_pb.RequestIsValidAddress, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseIsValidAddress) => void): grpc.ClientUnaryCall;
    isValidAddress(request: address_pb.RequestIsValidAddress, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseIsValidAddress) => void): grpc.ClientUnaryCall;
    isValidAddress(request: address_pb.RequestIsValidAddress, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseIsValidAddress) => void): grpc.ClientUnaryCall;
}

export class RippleAddressAPIClient extends grpc.Client implements IRippleAddressAPIClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public generateAddress(request: google_protobuf_empty_pb.Empty, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseGenerateAddress) => void): grpc.ClientUnaryCall;
    public generateAddress(request: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseGenerateAddress) => void): grpc.ClientUnaryCall;
    public generateAddress(request: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseGenerateAddress) => void): grpc.ClientUnaryCall;
    public generateXAddress(request: google_protobuf_empty_pb.Empty, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseGenerateXAddress) => void): grpc.ClientUnaryCall;
    public generateXAddress(request: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseGenerateXAddress) => void): grpc.ClientUnaryCall;
    public generateXAddress(request: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseGenerateXAddress) => void): grpc.ClientUnaryCall;
    public isValidAddress(request: address_pb.RequestIsValidAddress, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseIsValidAddress) => void): grpc.ClientUnaryCall;
    public isValidAddress(request: address_pb.RequestIsValidAddress, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseIsValidAddress) => void): grpc.ClientUnaryCall;
    public isValidAddress(request: address_pb.RequestIsValidAddress, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: address_pb.ResponseIsValidAddress) => void): grpc.ClientUnaryCall;
}
