import grpc, {sendUnaryData, ServerUnaryCall, ServiceError} from 'grpc';
import * as ripple from 'ripple-lib';
import * as grpc_pb from '../pb/account_grpc_pb';
import * as pb from '../pb/account_pb';
import { rippledError } from './errors';


export class RippleAccountAPIService implements grpc_pb.IRippleAccountAPIServer {
  private rippleAPI: ripple.RippleAPI;

  public constructor(rippleAPI: ripple.RippleAPI) {
    this.rippleAPI = rippleAPI;
  }

  // getAccountInfo handler
  getAccountInfo = (
    call: ServerUnaryCall<pb.RequestGetAccountInfo>,
    callback: sendUnaryData<pb.ResponseGetAccountInfo>,
  ) : void => {
    console.log("[getAccountInfo] is called");
    const address = call.request.getAddress();
    //const ledgerversion = call.request.getLedgerversion();

    this.rippleAPI.getAccountInfo(address)
    .then(info => {
      console.log("account xrpBalance:", info.xrpBalance);
      // response
      const res = new pb.ResponseGetAccountInfo();
      res.setSequence(info.sequence);
      res.setXrpbalance(info.xrpBalance);
      res.setOwnercount(info.ownerCount);
      res.setPreviousaffectingtransactionid(info.previousAffectingTransactionID);
      res.setPreviousaffectingtransactionledgerversion(info.previousAffectingTransactionLedgerVersion);
      callback(null, res);
    })
    .catch((error: rippledError) => {
      if (error) console.log(error.data); //sometimes, data is undefined 
      const statusError: ServiceError = {
        name: 'getAccountInfo error',
        message: error? error.data.error_message: 'something error',
        code: grpc.status.INVALID_ARGUMENT,
      };
      callback(statusError, null);
    });
  }

};

export const service = grpc_pb.RippleAccountAPIService;
