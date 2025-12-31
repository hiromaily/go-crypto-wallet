// package: rippleapi.account
// file: account.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class RequestGetAccountInfo extends jspb.Message { 
    getAddress(): string;
    setAddress(value: string): RequestGetAccountInfo;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RequestGetAccountInfo.AsObject;
    static toObject(includeInstance: boolean, msg: RequestGetAccountInfo): RequestGetAccountInfo.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RequestGetAccountInfo, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RequestGetAccountInfo;
    static deserializeBinaryFromReader(message: RequestGetAccountInfo, reader: jspb.BinaryReader): RequestGetAccountInfo;
}

export namespace RequestGetAccountInfo {
    export type AsObject = {
        address: string,
    }
}

export class ResponseGetAccountInfo extends jspb.Message { 
    getSequence(): number;
    setSequence(value: number): ResponseGetAccountInfo;
    getXrpbalance(): string;
    setXrpbalance(value: string): ResponseGetAccountInfo;
    getOwnercount(): number;
    setOwnercount(value: number): ResponseGetAccountInfo;
    getPreviousaffectingtransactionid(): string;
    setPreviousaffectingtransactionid(value: string): ResponseGetAccountInfo;
    getPreviousaffectingtransactionledgerversion(): number;
    setPreviousaffectingtransactionledgerversion(value: number): ResponseGetAccountInfo;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ResponseGetAccountInfo.AsObject;
    static toObject(includeInstance: boolean, msg: ResponseGetAccountInfo): ResponseGetAccountInfo.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ResponseGetAccountInfo, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ResponseGetAccountInfo;
    static deserializeBinaryFromReader(message: ResponseGetAccountInfo, reader: jspb.BinaryReader): ResponseGetAccountInfo;
}

export namespace ResponseGetAccountInfo {
    export type AsObject = {
        sequence: number,
        xrpbalance: string,
        ownercount: number,
        previousaffectingtransactionid: string,
        previousaffectingtransactionledgerversion: number,
    }
}
