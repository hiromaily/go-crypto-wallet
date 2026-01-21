/**
 * ConnectRPC Server Setup
 *
 * Sets up the ConnectRPC router with all XRP Ledger services.
 * Provides a fetch handler for use with Bun.serve().
 */

import { create } from '@bufbuild/protobuf';
import { type ConnectRouter, createConnectRouter } from '@connectrpc/connect';
import {
  universalServerRequestFromFetch,
  universalServerResponseToFetch,
} from '@connectrpc/connect/protocol';

import {
  type RequestGetAccountInfo,
  ResponseGetAccountInfoSchema,
  RippleAccountAPI,
} from './gen/account_pb';
import {
  type RequestIsValidAddress,
  ResponseGenerateAddressSchema,
  ResponseGenerateXAddressSchema,
  ResponseIsValidAddressSchema,
  RippleAddressAPI,
} from './gen/address_pb';
import {
  type EnumTransactionType,
  type RequestCombineTransaction,
  type RequestGetTransaction,
  type RequestPrepareTransaction,
  type RequestSignTransaction,
  type RequestSubmitTransaction,
  ResponseCombineTransactionSchema,
  ResponseGetTransactionSchema,
  ResponsePrepareTransactionSchema,
  ResponseSignTransactionSchema,
  ResponseSubmitTransactionSchema,
  ResponseWaitValidationSchema,
  RippleTransactionAPI,
} from './gen/transaction_pb';
import { accountService } from './services/account';
import { addressService } from './services/address';
import {
  type EnumTransactionType as ServiceEnumTxType,
  transactionService,
} from './services/transaction';

/**
 * Create the ConnectRPC router with all services registered
 *
 * @returns Configured ConnectRouter with all XRP services
 */
export function createRouter(): ConnectRouter {
  return createConnectRouter()
    .service(RippleAccountAPI, {
      getAccountInfo: async (request: RequestGetAccountInfo) => {
        const result = await accountService.getAccountInfo({
          address: request.address,
        });

        return create(ResponseGetAccountInfoSchema, {
          sequence: result.sequence,
          xrpBalance: result.xrpBalance,
          ownerCount: result.ownerCount,
          previousAffectingTransactionID: result.previousAffectingTransactionId,
          previousAffectingTransactionLedgerVersion:
            result.previousAffectingTransactionLedgerVersion,
        });
      },
    })
    .service(RippleAddressAPI, {
      generateAddress: () => {
        const result = addressService.generateAddress();

        return create(ResponseGenerateAddressSchema, {
          xAddress: result.xAddress,
          classicAddress: result.classicAddress,
          address: result.address,
          secret: result.secret,
        });
      },
      generateXAddress: () => {
        const result = addressService.generateXAddress();

        return create(ResponseGenerateXAddressSchema, {
          xAddress: result.xAddress,
          secret: result.secret,
        });
      },
      isValidAddress: (request: RequestIsValidAddress) => {
        const result = addressService.isValidAddress({
          address: request.address,
        });

        return create(ResponseIsValidAddressSchema, {
          isValid: result.isValid,
        });
      },
    })
    .service(RippleTransactionAPI, {
      prepareTransaction: async (request: RequestPrepareTransaction) => {
        // Build the request object
        type ServiceRequest = {
          txType: ServiceEnumTxType;
          senderAccount: string;
          amount: number;
          receiverAccount: string;
          instructions?: {
            fee?: string;
            maxFee?: string;
            maxLedgerVersion?: bigint;
            maxLedgerVersionOffset?: bigint;
            sequence?: bigint;
            signersCount?: bigint;
          };
          regularKey?: string;
          setFlag?: number;
          clearFlag?: number;
        };

        const serviceRequest: ServiceRequest = {
          txType: mapEnumTransactionType(request.txType),
          senderAccount: request.senderAccount,
          amount: request.amount,
          receiverAccount: request.receiverAccount,
        };

        // Map regular key fields for SetRegularKey transaction
        if (request.regularKey) {
          serviceRequest.regularKey = request.regularKey;
        }

        // Map AccountSet flags for enabling/disabling master key
        if (request.setFlag > 0) {
          serviceRequest.setFlag = request.setFlag;
        }
        if (request.clearFlag > 0) {
          serviceRequest.clearFlag = request.clearFlag;
        }

        // Map instructions if provided
        if (request.instructions) {
          const serviceInstructions: NonNullable<ServiceRequest['instructions']> = {};
          if (request.instructions.fee) {
            serviceInstructions.fee = request.instructions.fee;
          }
          if (request.instructions.maxFee) {
            serviceInstructions.maxFee = request.instructions.maxFee;
          }
          if (request.instructions.maxLedgerVersion) {
            serviceInstructions.maxLedgerVersion = request.instructions.maxLedgerVersion;
          }
          if (request.instructions.maxLedgerVersionOffset) {
            serviceInstructions.maxLedgerVersionOffset =
              request.instructions.maxLedgerVersionOffset;
          }
          if (request.instructions.sequence) {
            serviceInstructions.sequence = request.instructions.sequence;
          }
          if (request.instructions.signersCount) {
            serviceInstructions.signersCount = request.instructions.signersCount;
          }
          serviceRequest.instructions = serviceInstructions;
        }

        const result = await transactionService.prepareTransaction(serviceRequest);

        // Build response with optional instructions
        const responseInit: {
          txJSON: string;
          instructions?: {
            fee: string;
            maxFee: string;
            maxLedgerVersion: bigint;
            maxLedgerVersionOffset: bigint;
            sequence: bigint;
            signersCount: bigint;
          };
        } = {
          txJSON: result.txJson,
        };

        if (result.instructions) {
          responseInit.instructions = {
            fee: result.instructions.fee ?? '',
            maxFee: result.instructions.maxFee ?? '',
            maxLedgerVersion: result.instructions.maxLedgerVersion ?? 0n,
            maxLedgerVersionOffset: result.instructions.maxLedgerVersionOffset ?? 0n,
            sequence: result.instructions.sequence ?? 0n,
            signersCount: result.instructions.signersCount ?? 0n,
          };
        }

        return create(ResponsePrepareTransactionSchema, responseInit);
      },
      signTransaction: async (request: RequestSignTransaction) => {
        const result = await transactionService.signTransaction({
          txJson: request.txJSON,
          secret: request.secret,
        });

        return create(ResponseSignTransactionSchema, {
          txID: result.txId,
          txBlob: result.txBlob,
        });
      },
      submitTransaction: async (request: RequestSubmitTransaction) => {
        const result = await transactionService.submitTransaction({
          txBlob: request.txBlob,
        });

        return create(ResponseSubmitTransactionSchema, {
          resultJSONString: result.resultJsonString,
          earliestLedgerVersion: result.earliestLedgerVersion,
        });
      },
      getTransaction: async (request: RequestGetTransaction) => {
        const result = await transactionService.getTransaction({
          txId: request.txID,
          minLedgerVersion: request.minLedgerVersion,
        });

        return create(ResponseGetTransactionSchema, {
          resultJSONString: result.resultJsonString,
        });
      },
      combineTransaction: async (request: RequestCombineTransaction) => {
        const result = await transactionService.combineTransaction({
          signedTransactions: request.signedTransactions,
        });

        return create(ResponseCombineTransactionSchema, {
          signedTransaction: result.signedTransaction,
          txID: result.txId,
        });
      },
      waitValidation: async function* (_, context) {
        for await (const response of transactionService.waitValidation({
          signal: context.signal,
        })) {
          yield create(ResponseWaitValidationSchema, {
            ledgerVersion: response.ledgerVersion,
          });
        }
      },
    });
}

/**
 * Map proto enum to service enum.
 * Both enums have identical numeric values, so direct return is safe.
 */
function mapEnumTransactionType(protoEnum: EnumTransactionType): ServiceEnumTxType {
  return protoEnum;
}

/**
 * Create a fetch handler from the ConnectRPC router.
 *
 * Routes incoming requests to the appropriate service handler.
 * This is required because ConnectRouter is not directly a fetch handler.
 *
 * @returns Fetch handler function for use with Bun.serve()
 */
export function createFetchHandler(): (request: Request) => Promise<Response> {
  const router = createRouter();

  return async (request: Request): Promise<Response> => {
    const url = new URL(request.url);

    // Find matching handler based on path
    for (const handler of router.handlers) {
      if (url.pathname.startsWith(handler.requestPath)) {
        const uRequest = universalServerRequestFromFetch(request, {});
        const uResponse = await handler(uRequest);
        return universalServerResponseToFetch(uResponse);
      }
    }

    // No matching handler found
    return new Response('Not Found', { status: 404 });
  };
}
