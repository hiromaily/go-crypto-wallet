import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import * as grpc from 'grpc';
import * as ripple from 'ripple-lib';
import * as grpc_pb from '../pb/address_grpc_pb';
import * as pb from '../pb/address_pb';


export class RippleAddressAPIService implements grpc_pb.IRippleAddressAPIServer {
  private rippleAPI: ripple.RippleAPI;

  public constructor(rippleAPI: ripple.RippleAPI) {
    this.rippleAPI = rippleAPI;
  }

  // generateAddress handler
  generateAddress = (
    call: grpc.ServerUnaryCall<Empty>,
    callback: grpc.sendUnaryData<pb.ResponseGenerateAddress>,
  ) : void => {
    console.log("[generateAddress] is called");

    const generated = this.rippleAPI.generateAddress();
    // response
    const res = new pb.ResponseGenerateAddress();
    res.setXaddress(generated.xAddress);
    if (generated.classicAddress) {
      res.setClassicaddress(generated.classicAddress);
    }
    if (generated.address) {
      res.setAddress(generated.address);
    }
    res.setSecret(generated.secret);

    callback(null, res);
  }

  // generateXAddress handler
  generateXAddress = (
    call: grpc.ServerUnaryCall<Empty>,
    callback: grpc.sendUnaryData<pb.ResponseGenerateXAddress>,
  ) : void => {
    console.log("[generateXAddress] is called");

    const generated = this.rippleAPI.generateXAddress();
    // response
    const res = new pb.ResponseGenerateXAddress();
    res.setXaddress(generated.xAddress);
    res.setSecret(generated.secret);

    callback(null, res);
  }

  // isValidAddress handler
  isValidAddress = (
    call: grpc.ServerUnaryCall<pb.RequestIsValidAddress>,
    callback: grpc.sendUnaryData<pb.ResponseIsValidAddress>,
  ) : void => {
    const address = call.request.getAddress();
    const isValid = this.rippleAPI.isValidAddress(address);
    
    // response
    const res = new pb.ResponseIsValidAddress();
    res.setIsvalid(isValid);

    callback(null, res);
  }

};

export const service = grpc_pb.RippleAddressAPIService;
