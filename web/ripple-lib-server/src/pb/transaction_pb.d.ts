// package: rippleapi.transaction
// file: transaction.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as gogo_protobuf_gogoproto_gogo_pb from "./gogo/protobuf/gogoproto/gogo_pb";

export class Instructions extends jspb.Message { 
    getFee(): string;
    setFee(value: string): Instructions;
    getMaxfee(): string;
    setMaxfee(value: string): Instructions;
    getMaxledgerversion(): number;
    setMaxledgerversion(value: number): Instructions;
    getMaxledgerversionoffset(): number;
    setMaxledgerversionoffset(value: number): Instructions;
    getSequence(): number;
    setSequence(value: number): Instructions;
    getSignerscount(): number;
    setSignerscount(value: number): Instructions;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Instructions.AsObject;
    static toObject(includeInstance: boolean, msg: Instructions): Instructions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Instructions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Instructions;
    static deserializeBinaryFromReader(message: Instructions, reader: jspb.BinaryReader): Instructions;
}

export namespace Instructions {
    export type AsObject = {
        fee: string,
        maxfee: string,
        maxledgerversion: number,
        maxledgerversionoffset: number,
        sequence: number,
        signerscount: number,
    }
}

export class RequestPrepareTransaction extends jspb.Message { 
    getTxType(): EnumTransactionType;
    setTxType(value: EnumTransactionType): RequestPrepareTransaction;
    getSenderaccount(): string;
    setSenderaccount(value: string): RequestPrepareTransaction;
    getAmount(): number;
    setAmount(value: number): RequestPrepareTransaction;
    getReceiveraccount(): string;
    setReceiveraccount(value: string): RequestPrepareTransaction;

    hasInstructions(): boolean;
    clearInstructions(): void;
    getInstructions(): Instructions | undefined;
    setInstructions(value?: Instructions): RequestPrepareTransaction;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RequestPrepareTransaction.AsObject;
    static toObject(includeInstance: boolean, msg: RequestPrepareTransaction): RequestPrepareTransaction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RequestPrepareTransaction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RequestPrepareTransaction;
    static deserializeBinaryFromReader(message: RequestPrepareTransaction, reader: jspb.BinaryReader): RequestPrepareTransaction;
}

export namespace RequestPrepareTransaction {
    export type AsObject = {
        txType: EnumTransactionType,
        senderaccount: string,
        amount: number,
        receiveraccount: string,
        instructions?: Instructions.AsObject,
    }
}

export class ResponsePrepareTransaction extends jspb.Message { 
    getTxjson(): string;
    setTxjson(value: string): ResponsePrepareTransaction;

    hasInstructions(): boolean;
    clearInstructions(): void;
    getInstructions(): Instructions | undefined;
    setInstructions(value?: Instructions): ResponsePrepareTransaction;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ResponsePrepareTransaction.AsObject;
    static toObject(includeInstance: boolean, msg: ResponsePrepareTransaction): ResponsePrepareTransaction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ResponsePrepareTransaction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ResponsePrepareTransaction;
    static deserializeBinaryFromReader(message: ResponsePrepareTransaction, reader: jspb.BinaryReader): ResponsePrepareTransaction;
}

export namespace ResponsePrepareTransaction {
    export type AsObject = {
        txjson: string,
        instructions?: Instructions.AsObject,
    }
}

export class RequestSignTransaction extends jspb.Message { 
    getTxjson(): string;
    setTxjson(value: string): RequestSignTransaction;
    getSecret(): string;
    setSecret(value: string): RequestSignTransaction;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RequestSignTransaction.AsObject;
    static toObject(includeInstance: boolean, msg: RequestSignTransaction): RequestSignTransaction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RequestSignTransaction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RequestSignTransaction;
    static deserializeBinaryFromReader(message: RequestSignTransaction, reader: jspb.BinaryReader): RequestSignTransaction;
}

export namespace RequestSignTransaction {
    export type AsObject = {
        txjson: string,
        secret: string,
    }
}

export class ResponseSignTransaction extends jspb.Message { 
    getTxid(): string;
    setTxid(value: string): ResponseSignTransaction;
    getTxblob(): string;
    setTxblob(value: string): ResponseSignTransaction;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ResponseSignTransaction.AsObject;
    static toObject(includeInstance: boolean, msg: ResponseSignTransaction): ResponseSignTransaction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ResponseSignTransaction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ResponseSignTransaction;
    static deserializeBinaryFromReader(message: ResponseSignTransaction, reader: jspb.BinaryReader): ResponseSignTransaction;
}

export namespace ResponseSignTransaction {
    export type AsObject = {
        txid: string,
        txblob: string,
    }
}

export class RequestSubmitTransaction extends jspb.Message { 
    getTxblob(): string;
    setTxblob(value: string): RequestSubmitTransaction;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RequestSubmitTransaction.AsObject;
    static toObject(includeInstance: boolean, msg: RequestSubmitTransaction): RequestSubmitTransaction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RequestSubmitTransaction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RequestSubmitTransaction;
    static deserializeBinaryFromReader(message: RequestSubmitTransaction, reader: jspb.BinaryReader): RequestSubmitTransaction;
}

export namespace RequestSubmitTransaction {
    export type AsObject = {
        txblob: string,
    }
}

export class ResponseSubmitTransaction extends jspb.Message { 
    getResultjsonstring(): string;
    setResultjsonstring(value: string): ResponseSubmitTransaction;
    getEarliestledgerversion(): number;
    setEarliestledgerversion(value: number): ResponseSubmitTransaction;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ResponseSubmitTransaction.AsObject;
    static toObject(includeInstance: boolean, msg: ResponseSubmitTransaction): ResponseSubmitTransaction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ResponseSubmitTransaction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ResponseSubmitTransaction;
    static deserializeBinaryFromReader(message: ResponseSubmitTransaction, reader: jspb.BinaryReader): ResponseSubmitTransaction;
}

export namespace ResponseSubmitTransaction {
    export type AsObject = {
        resultjsonstring: string,
        earliestledgerversion: number,
    }
}

export class ResponseWaitValidation extends jspb.Message { 
    getLedgerversion(): number;
    setLedgerversion(value: number): ResponseWaitValidation;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ResponseWaitValidation.AsObject;
    static toObject(includeInstance: boolean, msg: ResponseWaitValidation): ResponseWaitValidation.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ResponseWaitValidation, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ResponseWaitValidation;
    static deserializeBinaryFromReader(message: ResponseWaitValidation, reader: jspb.BinaryReader): ResponseWaitValidation;
}

export namespace ResponseWaitValidation {
    export type AsObject = {
        ledgerversion: number,
    }
}

export class RequestGetTransaction extends jspb.Message { 
    getTxid(): string;
    setTxid(value: string): RequestGetTransaction;
    getMinledgerversion(): number;
    setMinledgerversion(value: number): RequestGetTransaction;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RequestGetTransaction.AsObject;
    static toObject(includeInstance: boolean, msg: RequestGetTransaction): RequestGetTransaction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RequestGetTransaction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RequestGetTransaction;
    static deserializeBinaryFromReader(message: RequestGetTransaction, reader: jspb.BinaryReader): RequestGetTransaction;
}

export namespace RequestGetTransaction {
    export type AsObject = {
        txid: string,
        minledgerversion: number,
    }
}

export class ResponseGetTransaction extends jspb.Message { 
    getResultjsonstring(): string;
    setResultjsonstring(value: string): ResponseGetTransaction;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ResponseGetTransaction.AsObject;
    static toObject(includeInstance: boolean, msg: ResponseGetTransaction): ResponseGetTransaction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ResponseGetTransaction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ResponseGetTransaction;
    static deserializeBinaryFromReader(message: ResponseGetTransaction, reader: jspb.BinaryReader): ResponseGetTransaction;
}

export namespace ResponseGetTransaction {
    export type AsObject = {
        resultjsonstring: string,
    }
}

export class RequestCombineTransaction extends jspb.Message { 
    clearSignedtransactionsList(): void;
    getSignedtransactionsList(): Array<string>;
    setSignedtransactionsList(value: Array<string>): RequestCombineTransaction;
    addSignedtransactions(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RequestCombineTransaction.AsObject;
    static toObject(includeInstance: boolean, msg: RequestCombineTransaction): RequestCombineTransaction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RequestCombineTransaction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RequestCombineTransaction;
    static deserializeBinaryFromReader(message: RequestCombineTransaction, reader: jspb.BinaryReader): RequestCombineTransaction;
}

export namespace RequestCombineTransaction {
    export type AsObject = {
        signedtransactionsList: Array<string>,
    }
}

export class ResponseCombineTransaction extends jspb.Message { 
    getSignedtransaction(): string;
    setSignedtransaction(value: string): ResponseCombineTransaction;
    getTxid(): string;
    setTxid(value: string): ResponseCombineTransaction;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ResponseCombineTransaction.AsObject;
    static toObject(includeInstance: boolean, msg: ResponseCombineTransaction): ResponseCombineTransaction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ResponseCombineTransaction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ResponseCombineTransaction;
    static deserializeBinaryFromReader(message: ResponseCombineTransaction, reader: jspb.BinaryReader): ResponseCombineTransaction;
}

export namespace ResponseCombineTransaction {
    export type AsObject = {
        signedtransaction: string,
        txid: string,
    }
}

export enum EnumTransactionType {
    TX_ACCOUNT_SET = 0,
    TX_ACCOUNT_DELETE = 1,
    TX_CHECK_CANCEL = 2,
    TX_CHECK_CASH = 3,
    TX_CHECK_CREATE = 4,
    TX_DEPOSIT_PREAUTH = 5,
    TX_ESCROW_CANCEL = 6,
    TX_ESCROW_CREATE = 7,
    TX_ESCROW_FINISH = 8,
    TX_OFFER_CANCEL = 9,
    TX_OFFER_CREATE = 10,
    TX_PAYMENT = 11,
    TX_PAYMENT_CHANNEL_CLAIM = 12,
    TX_PAYMENT_CHANNEL_CREATE = 13,
    TX_PAYMENT_CHANNEL_FUND = 14,
    TX_SET_REGULAR_KEY = 15,
    TX_SINGER_LIST_SET = 16,
    TX_TRUST_SET = 17,
}
