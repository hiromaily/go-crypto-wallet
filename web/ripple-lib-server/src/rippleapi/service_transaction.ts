import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import grpc, {ServiceError} from 'grpc';
import * as ripple from 'ripple-lib';
import * as grpc_pb from '../pb/transaction_grpc_pb';
import * as pb from '../pb/transaction_pb';
import { enumTransactionTypeString } from './enum';
//import { rippledError } from './errors';

// this document may be useful
// https://qiita.com/aanrii/items/699b4cda0babb3f47a2f

interface transaction {
  TransactionType: string;
  Account: string;
  Amount: string;
  Destination: string;
}

interface resSubmitTransaction {
  resJSON: any;
  earlistLedgerVersion: number;
}

interface resGetTransaction {
  txJSON: string;
  errMessage: string;
}

interface resCombineTransaction {
  signedTransaction: string;
  txJSON: string;
}

interface rippleInstructions {
  sequence?: number;
  fee?: string;
  maxFee?: string;
  maxLedgerVersion?: number;
  maxLedgerVersionOffset?: number;
  signersCount?: number;
}

export class RippleTransactionAPIService implements grpc_pb.IRippleTransactionAPIServer {
  private rippleAPI: ripple.RippleAPI;

  // public constructor(wsURL: string) {
  //   this.rippleAPI = new ripple.RippleAPI({server: wsURL});
  //   this.rippleAPI.connect();
  // }
  public constructor(rippleAPI: ripple.RippleAPI) {
    this.rippleAPI = rippleAPI;
  }

  private async _prepareTransaction(call: grpc.ServerUnaryCall<pb.RequestPrepareTransaction>) : Promise<string> {
    console.log("_prepareTransaction()");

    const txType = call.request.getTxType();
    const instructions = call.request.getInstructions();
    console.log('maxLedgerVersionOffset: ', instructions?.getMaxledgerversionoffset());
    console.log('sequence: ', instructions?.getSequence());

    // create parameter instructions
    const paramInst: rippleInstructions = {};
    if (instructions?.getMaxledgerversionoffset()) {
      paramInst.maxLedgerVersionOffset = instructions?.getMaxledgerversionoffset();
    }
    if (instructions?.getMaxledgerversion()) {
      paramInst.maxLedgerVersion = instructions?.getMaxledgerversion();
    }
    if (instructions?.getSequence()) {
      paramInst.sequence = instructions?.getSequence();
    }
    if (instructions?.getFee()) {
      paramInst.fee = instructions?.getFee();
    }
    if (instructions?.getMaxfee()) {
      paramInst.maxFee = instructions?.getMaxfee();
    }
    if (instructions?.getSignerscount()) {
      paramInst.signersCount = instructions?.getSignerscount();
    }
    console.log('paramInst:', paramInst);

    // prepareTransaction()
    const preparedTx = await this.rippleAPI.prepareTransaction({
      "TransactionType": enumTransactionTypeString[txType],
      "Account": call.request.getSenderaccount(),
      "Amount": this.rippleAPI.xrpToDrops(call.request.getAmount().toString()),
      "Destination": call.request.getReceiveraccount(),      
    }, paramInst);
    return preparedTx.txJSON;
  }

  private async _submitTransaction(call: grpc.ServerUnaryCall<pb.RequestSubmitTransaction>) : Promise<resSubmitTransaction> {
    console.log("_submitTransaction()");

    const latestLedgerVersion = await this.rippleAPI.getLedgerVersion();
    const txBlob = call.request.getTxblob();
    const resJSON = await this.rippleAPI.submit(txBlob);
    console.log("Tentative result code:", resJSON.resultCode);
    console.log("Tentative result message:", resJSON.resultMessage);  
    
    return { resJSON: resJSON, earlistLedgerVersion: latestLedgerVersion + 1 };
  }

  // prepareTransaction handler
  prepareTransaction = (
    call: grpc.ServerUnaryCall<pb.RequestPrepareTransaction>,
    callback: grpc.sendUnaryData<pb.ResponsePrepareTransaction>,
  ) : void => {
    console.log("[prepareTransaction] is called");

    if (!this.rippleAPI.isConnected()) {
      const statusError: ServiceError = {
        name: 'connection error',
        message: 'connection error',
        code: grpc.status.INVALID_ARGUMENT,
      };
      callback(statusError, null);
    }

    // call API as async
    this._prepareTransaction(call).then(resJSON => {
      const txJSON = JSON.stringify(resJSON);
      //console.log("txJSON", txJSON);

      // response
      const res = new pb.ResponsePrepareTransaction();
      res.setTxjson(txJSON);
      callback(null, res);
    })
  }

  // signTransaction handler
  // https://xrpl.org/rippleapi-reference.html#offline-functionality
  signTransaction = (
    call: grpc.ServerUnaryCall<pb.RequestSignTransaction>,
    callback: grpc.sendUnaryData<pb.ResponseSignTransaction>,
  ) : void => {
    console.log("[signTransaction] is called");
  
    // call API
    const signed = this.rippleAPI.sign(call.request.getTxjson(), call.request.getSecret());
    console.log("txID: Identifying hash:", signed.id);
    console.log("txBlob: Signed blob:", signed.signedTransaction);
  
    // response
    const res = new pb.ResponseSignTransaction();
    res.setTxid(signed.id);
    res.setTxblob(signed.signedTransaction);
    callback(null, res);
  }

  // submitTransaction handler
  submitTransaction = (
    call: grpc.ServerUnaryCall<pb.RequestSubmitTransaction>,
    callback: grpc.sendUnaryData<pb.ResponseSubmitTransaction>,
  ) : void => {
    console.log("[submitTransaction] is called");

    // call API as async
    this._submitTransaction(call).then(resAPI => {
      const resJSON = JSON.stringify(resAPI.resJSON);
      console.log("resJSON", resJSON);
      console.log("earlistLedgerVersion", resAPI.earlistLedgerVersion);

      // response
      const res = new pb.ResponseSubmitTransaction();
      res.setResultjsonstring(resJSON);
      res.setEarliestledgerversion(resAPI.earlistLedgerVersion);
      callback(null, res);
    })
  }

  // waitValidation as server streaming
  waitValidation = (call: grpc.ServerWritableStream<Empty>,
  ) : void => {
    console.log("[waitValidation] is called");

    const ledgerHandler = (ledger: any) => {
      if (call.cancelled) {
        console.log("canceled");
        call.end();
        this.rippleAPI.removeListener('ledger', ledgerHandler);
        return;
      }

      console.log("Ledger version", ledger.ledgerVersion, "was just validated.", call.cancelled);
      // response
      const res = new pb.ResponseWaitValidation();
      res.setLedgerversion(<number>ledger.ledgerVersion);
      call.write(res);
    }
    this.rippleAPI.on('ledger', ledgerHandler);
    // this.rippleAPI.on('ledger', ledger => {
    //   console.log("Ledger version", ledger.ledgerVersion, "was just validated.");
    //   // if (ledger.ledgerVersion > maxLedgerVersion) {
    //   //   console.log("If the transaction hasn't succeeded by now, it's expired")
    //   // }
    //   call.write(ledger.ledgerVersion);
    // });

    // when disconnected, remove listener
    // FIXME: this event is not called
    call.on('close', () => {
      console.log("[close] is called");
      call.end();
      this.rippleAPI.removeListener('ledger', ledgerHandler);
    });
  }

  // getTransaction handler
  // - Ledger History
  // - https://xrpl.org/ledger-history.html
  getTransaction = (
    call: grpc.ServerUnaryCall<pb.RequestGetTransaction>,
    callback: grpc.sendUnaryData<pb.ResponseGetTransaction>,
  ) : void => {
    console.log("[getTransaction] is called");

    const txID = call.request.getTxid();
    const earliestLedgerVersion = call.request.getMinledgerversion();
    console.log(`earliestLedgerVersion: ${earliestLedgerVersion}`);

    this.rippleAPI.getTransaction(txID, {minLedgerVersion: earliestLedgerVersion})
    .then(tx => {
      console.log("Transaction result:", tx.outcome.result);
      console.log("Balance changes:", JSON.stringify(tx.outcome.balanceChanges));

      // response
      const res = new pb.ResponseGetTransaction();
      res.setResultjsonstring(JSON.stringify(tx));
      callback(null, res);

    })
    .catch((error: Error) => {
      if (error) {
        console.log(error.name);
        console.log(error.message);
      }   
      // MissingLedgerHistoryError: Server is missing ledger history in the specified range
      // NotFoundError: Transaction has not been validated yet; try again later
      // NotFoundError: Transaction not found
      const statusError: ServiceError = {
        name: error.name? `getTransaction error ${error.name}`: 'getTransaction error',
        message: error.message? error.message: 'something error',
        code: grpc.status.INVALID_ARGUMENT,
      };
      callback(statusError, null);
    });
  }

  // combineTransaction handler
  combineTransaction = (
    call: grpc.ServerUnaryCall<pb.RequestCombineTransaction>,
    callback: grpc.sendUnaryData<pb.ResponseCombineTransaction>,
  ) : void => {
    console.log("[combineTransaction] is called");

    const signedObj = this.rippleAPI.combine(call.request.getSignedtransactionsList());
    console.log('signedObj:', signedObj);
    //resCombineTransaction
    // response
    const res = new pb.ResponseCombineTransaction();
    //if(typeof signed === 'object' && 'signedTransaction' in signed) {
    if (isResCombineTransaction(signedObj)){
      const signed = signedObj as resCombineTransaction;
      res.setSignedtransaction(<string>signed.signedTransaction);
      res.setTxid(signed.txJSON);  
    }
    callback(null, res);
  }
};

const isResCombineTransaction = (obj: any): obj is resCombineTransaction =>
  obj.signedTransaction && obj.txJSON;

// export default {
//   service: transaction_grpc_pb.RippleTransactionAPIService,  // Service interface
//   impl: new RippleTransactionAPIService(wsURL),              // Service interface definitions
// };
export const service = grpc_pb.RippleTransactionAPIService;
