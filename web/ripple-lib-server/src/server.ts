import 'dotenv/config';
import * as grpc from 'grpc';
import * as ripple from 'ripple-lib';
import * as rippleTxAPI from './rippleapi/service_transaction';
import * as rippleAccountAPI from './rippleapi/service_account';
import * as rippleAddressAPI from './rippleapi/service_address';


const port: string | number = process.env.PORT || 50051;
const wsURL: string = process.env.RippleAPIURL || 'wss://s.altnet.rippletest.net:51233';

type StartServerType = () => void;
export const startServer: StartServerType = (): void => {
  // Note: if it run as offline mode, run without parameter. `new RippleAPI();`
  // https://xrpl.org/rippleapi-reference.html#offline-functionality
  const rippleAPI = new ripple.RippleAPI({server: wsURL});

  // event handler
  rippleAPI.on('error', (errorCode, errorMessage) => {
    console.log(errorCode + ': ' + errorMessage);
  });
  rippleAPI.on('disconnected', (code) => {
    // code - [close code](https://developer.mozilla.org/en-US/docs/Web/API/CloseEvent) sent by the server
    // will be 1000 if this was normal closure
    console.log('disconnected, code:', code);
  });

  // connect to ripple server
  // https://xrpl.org/rippleapi-reference.html#boilerplate
  rippleAPI.connect().then(() => {
    console.log('connected');
  }).catch((e) => {
    console.log(e);
  });

  // grpc setting
  const server = new grpc.Server();
  server.addService(
    rippleTxAPI.service,
    new rippleTxAPI.RippleTransactionAPIService(rippleAPI),
  );
  server.addService(
    rippleAccountAPI.service,
    new rippleAccountAPI.RippleAccountAPIService(rippleAPI),
  );
  server.addService(
    rippleAddressAPI.service,
    new rippleAddressAPI.RippleAddressAPIService(rippleAPI),
  );

  // run server
  server.bindAsync(
    `0.0.0.0:${ port }`,
    grpc.ServerCredentials.createInsecure(),
    (error: Error | null, port: number) => {
      if (error != null) {
          return console.error(error);
      }
      console.log(`gRPC listening on ${ port }`);
    },
  );

  server.start();
};

startServer();
