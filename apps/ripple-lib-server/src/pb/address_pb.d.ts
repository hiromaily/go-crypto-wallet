// package: rippleapi.address
// file: address.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";

export class ResponseGenerateAddress extends jspb.Message { 
    getXaddress(): string;
    setXaddress(value: string): ResponseGenerateAddress;
    getClassicaddress(): string;
    setClassicaddress(value: string): ResponseGenerateAddress;
    getAddress(): string;
    setAddress(value: string): ResponseGenerateAddress;
    getSecret(): string;
    setSecret(value: string): ResponseGenerateAddress;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ResponseGenerateAddress.AsObject;
    static toObject(includeInstance: boolean, msg: ResponseGenerateAddress): ResponseGenerateAddress.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ResponseGenerateAddress, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ResponseGenerateAddress;
    static deserializeBinaryFromReader(message: ResponseGenerateAddress, reader: jspb.BinaryReader): ResponseGenerateAddress;
}

export namespace ResponseGenerateAddress {
    export type AsObject = {
        xaddress: string,
        classicaddress: string,
        address: string,
        secret: string,
    }
}

export class ResponseGenerateXAddress extends jspb.Message { 
    getXaddress(): string;
    setXaddress(value: string): ResponseGenerateXAddress;
    getSecret(): string;
    setSecret(value: string): ResponseGenerateXAddress;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ResponseGenerateXAddress.AsObject;
    static toObject(includeInstance: boolean, msg: ResponseGenerateXAddress): ResponseGenerateXAddress.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ResponseGenerateXAddress, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ResponseGenerateXAddress;
    static deserializeBinaryFromReader(message: ResponseGenerateXAddress, reader: jspb.BinaryReader): ResponseGenerateXAddress;
}

export namespace ResponseGenerateXAddress {
    export type AsObject = {
        xaddress: string,
        secret: string,
    }
}

export class RequestIsValidAddress extends jspb.Message { 
    getAddress(): string;
    setAddress(value: string): RequestIsValidAddress;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RequestIsValidAddress.AsObject;
    static toObject(includeInstance: boolean, msg: RequestIsValidAddress): RequestIsValidAddress.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RequestIsValidAddress, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RequestIsValidAddress;
    static deserializeBinaryFromReader(message: RequestIsValidAddress, reader: jspb.BinaryReader): RequestIsValidAddress;
}

export namespace RequestIsValidAddress {
    export type AsObject = {
        address: string,
    }
}

export class ResponseIsValidAddress extends jspb.Message { 
    getIsvalid(): boolean;
    setIsvalid(value: boolean): ResponseIsValidAddress;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ResponseIsValidAddress.AsObject;
    static toObject(includeInstance: boolean, msg: ResponseIsValidAddress): ResponseIsValidAddress.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ResponseIsValidAddress, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ResponseIsValidAddress;
    static deserializeBinaryFromReader(message: ResponseIsValidAddress, reader: jspb.BinaryReader): ResponseIsValidAddress;
}

export namespace ResponseIsValidAddress {
    export type AsObject = {
        isvalid: boolean,
    }
}
