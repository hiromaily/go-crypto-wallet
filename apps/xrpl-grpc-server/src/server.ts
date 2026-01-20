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
  EnumTransactionType,
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
  EnumTransactionType as ServiceEnumTxType,
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
        };

        const serviceRequest: ServiceRequest = {
          txType: mapEnumTransactionType(request.txType),
          senderAccount: request.senderAccount,
          amount: request.amount,
          receiverAccount: request.receiverAccount,
        };

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
            fee: result.instructions.fee || '',
            maxFee: result.instructions.maxFee || '',
            maxLedgerVersion: result.instructions.maxLedgerVersion || BigInt(0),
            maxLedgerVersionOffset: result.instructions.maxLedgerVersionOffset || BigInt(0),
            sequence: result.instructions.sequence || BigInt(0),
            signersCount: result.instructions.signersCount || BigInt(0),
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
        const abortController = new AbortController();

        // Connect the context signal to our abort controller
        context.signal.addEventListener('abort', () => {
          abortController.abort();
        });

        for await (const response of transactionService.waitValidation({
          signal: abortController.signal,
        })) {
          yield create(ResponseWaitValidationSchema, {
            ledgerVersion: response.ledgerVersion,
          });
        }
      },
    });
}

/**
 * Map proto enum to service enum
 */
function mapEnumTransactionType(protoEnum: EnumTransactionType): ServiceEnumTxType {
  const mapping: Record<EnumTransactionType, ServiceEnumTxType> = {
    [EnumTransactionType.TX_ACCOUNT_SET]: ServiceEnumTxType.TX_ACCOUNT_SET,
    [EnumTransactionType.TX_ACCOUNT_DELETE]: ServiceEnumTxType.TX_ACCOUNT_DELETE,
    [EnumTransactionType.TX_CHECK_CANCEL]: ServiceEnumTxType.TX_CHECK_CANCEL,
    [EnumTransactionType.TX_CHECK_CASH]: ServiceEnumTxType.TX_CHECK_CASH,
    [EnumTransactionType.TX_CHECK_CREATE]: ServiceEnumTxType.TX_CHECK_CREATE,
    [EnumTransactionType.TX_DEPOSIT_PREAUTH]: ServiceEnumTxType.TX_DEPOSIT_PREAUTH,
    [EnumTransactionType.TX_ESCROW_CANCEL]: ServiceEnumTxType.TX_ESCROW_CANCEL,
    [EnumTransactionType.TX_ESCROW_CREATE]: ServiceEnumTxType.TX_ESCROW_CREATE,
    [EnumTransactionType.TX_ESCROW_FINISH]: ServiceEnumTxType.TX_ESCROW_FINISH,
    [EnumTransactionType.TX_OFFER_CANCEL]: ServiceEnumTxType.TX_OFFER_CANCEL,
    [EnumTransactionType.TX_OFFER_CREATE]: ServiceEnumTxType.TX_OFFER_CREATE,
    [EnumTransactionType.TX_PAYMENT]: ServiceEnumTxType.TX_PAYMENT,
    [EnumTransactionType.TX_PAYMENT_CHANNEL_CLAIM]: ServiceEnumTxType.TX_PAYMENT_CHANNEL_CLAIM,
    [EnumTransactionType.TX_PAYMENT_CHANNEL_CREATE]: ServiceEnumTxType.TX_PAYMENT_CHANNEL_CREATE,
    [EnumTransactionType.TX_PAYMENT_CHANNEL_FUND]: ServiceEnumTxType.TX_PAYMENT_CHANNEL_FUND,
    [EnumTransactionType.TX_SET_REGULAR_KEY]: ServiceEnumTxType.TX_SET_REGULAR_KEY,
    [EnumTransactionType.TX_SINGER_LIST_SET]: ServiceEnumTxType.TX_SINGER_LIST_SET,
    [EnumTransactionType.TX_TRUST_SET]: ServiceEnumTxType.TX_TRUST_SET,
  };
  return mapping[protoEnum];
}

/**
 * Create a fetch handler from the ConnectRPC router
 *
 * This converts the ConnectRPC handlers into a single fetch function
 * suitable for use with Bun.serve().
 *
 * @returns Fetch handler function
 */
export function createServerFetchHandler(): (request: Request) => Promise<Response> {
  const router = createRouter();

  return async (request: Request): Promise<Response> => {
    const url = new URL(request.url);

    // Find matching handler based on path
    for (const handler of router.handlers) {
      if (url.pathname.startsWith(handler.requestPath)) {
        const uRequest = universalServerRequestFromFetch(request, {
          httpVersion: '2',
        });

        const uResponse = await handler(uRequest);
        return universalServerResponseToFetch(uResponse);
      }
    }

    // No matching handler found
    return new Response('Not Found', { status: 404 });
  };
}
