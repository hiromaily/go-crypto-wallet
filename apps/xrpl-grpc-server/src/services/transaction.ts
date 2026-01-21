/**
 * RippleTransactionAPI Service
 *
 * Implements transaction methods for the XRP Ledger.
 * Uses xrpl.js 4.5.0 for transaction preparation, signing, submission, and monitoring.
 *
 * Security Note: Secrets must never be logged.
 */

import { type Client, multisign, type SubmittableTransaction, Wallet, xrpToDrops } from 'xrpl';
import { getClient } from '../xrpl';

// ============================================================================
// Transaction Type Enum Mapping
// ============================================================================

/**
 * Enum values matching proto/rippleapi/transaction.proto EnumTransactionType
 */
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

/**
 * Maps enum values to XRPL TransactionType strings
 */
const enumToTransactionType: Record<EnumTransactionType, string> = {
  [EnumTransactionType.TX_ACCOUNT_SET]: 'AccountSet',
  [EnumTransactionType.TX_ACCOUNT_DELETE]: 'AccountDelete',
  [EnumTransactionType.TX_CHECK_CANCEL]: 'CheckCancel',
  [EnumTransactionType.TX_CHECK_CASH]: 'CheckCash',
  [EnumTransactionType.TX_CHECK_CREATE]: 'CheckCreate',
  [EnumTransactionType.TX_DEPOSIT_PREAUTH]: 'DepositPreauth',
  [EnumTransactionType.TX_ESCROW_CANCEL]: 'EscrowCancel',
  [EnumTransactionType.TX_ESCROW_CREATE]: 'EscrowCreate',
  [EnumTransactionType.TX_ESCROW_FINISH]: 'EscrowFinish',
  [EnumTransactionType.TX_OFFER_CANCEL]: 'OfferCancel',
  [EnumTransactionType.TX_OFFER_CREATE]: 'OfferCreate',
  [EnumTransactionType.TX_PAYMENT]: 'Payment',
  [EnumTransactionType.TX_PAYMENT_CHANNEL_CLAIM]: 'PaymentChannelClaim',
  [EnumTransactionType.TX_PAYMENT_CHANNEL_CREATE]: 'PaymentChannelCreate',
  [EnumTransactionType.TX_PAYMENT_CHANNEL_FUND]: 'PaymentChannelFund',
  [EnumTransactionType.TX_SET_REGULAR_KEY]: 'SetRegularKey',
  [EnumTransactionType.TX_SINGER_LIST_SET]: 'SignerListSet',
  [EnumTransactionType.TX_TRUST_SET]: 'TrustSet',
};

// ============================================================================
// Request/Response Types
// ============================================================================

/**
 * Transaction instructions (optional parameters)
 */
export interface Instructions {
  fee?: string;
  maxFee?: string;
  maxLedgerVersion?: bigint;
  maxLedgerVersionOffset?: bigint;
  sequence?: bigint;
  signersCount?: bigint;
}

/**
 * SignerEntry for SignerListSet transaction
 */
export interface SignerEntry {
  account: string;
  weight: number;
}

/**
 * Request for PrepareTransaction RPC
 */
export interface PrepareTransactionRequest {
  txType: EnumTransactionType;
  senderAccount: string;
  amount: number;
  receiverAccount: string;
  instructions?: Instructions;
  /** For SetRegularKey: the address to set as regular key (empty to remove) */
  regularKey?: string;
  /** For AccountSet: flag to set (e.g., 4 = asfDisableMaster) */
  setFlag?: number;
  /** For AccountSet: flag to clear */
  clearFlag?: number;
  /** For SignerListSet: minimum total weight of signatures required (0 to remove signer list) */
  signerQuorum?: number;
  /** For SignerListSet: list of signers with their weights */
  signerEntries?: SignerEntry[];
}

/**
 * Response for PrepareTransaction RPC
 */
export interface PrepareTransactionResponse {
  txJson: string;
  instructions?: Instructions;
}

/**
 * Request for SignTransaction RPC
 */
export interface SignTransactionRequest {
  txJson: string;
  secret: string;
}

/**
 * Response for SignTransaction RPC
 */
export interface SignTransactionResponse {
  txId: string;
  txBlob: string;
}

/**
 * Request for SubmitTransaction RPC
 */
export interface SubmitTransactionRequest {
  txBlob: string;
}

/**
 * Response for SubmitTransaction RPC
 */
export interface SubmitTransactionResponse {
  resultJsonString: string;
  earliestLedgerVersion: bigint;
}

/**
 * Request for GetTransaction RPC
 */
export interface GetTransactionRequest {
  txId: string;
  minLedgerVersion: bigint;
}

/**
 * Response for GetTransaction RPC
 */
export interface GetTransactionResponse {
  resultJsonString: string;
}

/**
 * Request for CombineTransaction RPC
 */
export interface CombineTransactionRequest {
  signedTransactions: string[];
}

/**
 * Response for CombineTransaction RPC
 */
export interface CombineTransactionResponse {
  signedTransaction: string;
  txId: string;
}

/**
 * Response for WaitValidation RPC (streaming)
 */
export interface WaitValidationResponse {
  ledgerVersion: bigint;
}

/**
 * Context for streaming operations (includes abort signal)
 */
export interface StreamingContext {
  signal: AbortSignal;
}

// ============================================================================
// Error Classes
// ============================================================================

/**
 * Error thrown when transaction preparation fails
 */
export class PrepareTransactionError extends Error {
  constructor(
    message: string,
    public readonly code?: string,
  ) {
    super(message);
    this.name = 'PrepareTransactionError';
  }
}

/**
 * Error thrown when transaction signing fails
 */
export class SignTransactionError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'SignTransactionError';
  }
}

/**
 * Error thrown when transaction submission fails
 */
export class SubmitTransactionError extends Error {
  constructor(
    message: string,
    public readonly code?: string,
  ) {
    super(message);
    this.name = 'SubmitTransactionError';
  }
}

/**
 * Error thrown when transaction retrieval fails
 */
export class GetTransactionError extends Error {
  constructor(
    message: string,
    public readonly code?: string,
  ) {
    super(message);
    this.name = 'GetTransactionError';
  }
}

/**
 * Error thrown when transaction is not yet validated
 */
export class TransactionNotValidatedError extends Error {
  constructor(txId: string) {
    super(`Transaction not yet validated: ${txId}`);
    this.name = 'TransactionNotValidatedError';
  }
}

/**
 * Error thrown when transaction combine fails
 */
export class CombineTransactionError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'CombineTransactionError';
  }
}

// ============================================================================
// Service Factory
// ============================================================================

/**
 * Create the transaction service with a client getter function
 *
 * This factory pattern allows for dependency injection in tests.
 *
 * @param clientGetter - Function to get the XRPL client
 * @returns Transaction service implementation
 */
export function createTransactionService(clientGetter: () => Promise<Client>) {
  return {
    /**
     * Prepare a transaction for signing
     *
     * Uses client.autofill() to automatically fill in transaction fields
     * such as Sequence, Fee, and LastLedgerSequence.
     *
     * @param request - Transaction parameters
     * @returns Prepared transaction JSON and instructions
     * @throws PrepareTransactionError if preparation fails
     */
    prepareTransaction: async (
      request: PrepareTransactionRequest,
    ): Promise<PrepareTransactionResponse> => {
      const client = await clientGetter();

      try {
        // Build the base transaction object dynamically based on transaction type.
        // Using Record<string, unknown> for type safety while allowing dynamic properties.
        // biome-ignore lint/suspicious/noExplicitAny: Dynamic transaction building requires flexible typing
        const tx: Record<string, any> = {
          TransactionType: enumToTransactionType[request.txType],
          Account: request.senderAccount,
        };

        // Add transaction-type-specific fields
        switch (request.txType) {
          case EnumTransactionType.TX_PAYMENT:
            tx.Amount = xrpToDrops(request.amount.toString());
            tx.Destination = request.receiverAccount;
            break;

          case EnumTransactionType.TX_SET_REGULAR_KEY:
            // SetRegularKey transaction
            // - RegularKey field: set to assign, omit to remove regular key
            // Reference: https://xrpl.org/docs/references/protocol/transactions/types/setregularkey
            if (request.regularKey && request.regularKey.length > 0) {
              tx.RegularKey = request.regularKey;
            }
            // If regularKey is empty/undefined, the transaction will remove the regular key
            break;

          case EnumTransactionType.TX_ACCOUNT_SET:
            // AccountSet transaction for account configuration
            // Reference: https://xrpl.org/docs/references/protocol/transactions/types/accountset
            // Common flags:
            //   4 (asfDisableMaster): Disable master key signing
            //   8 (asfNoFreeze): Permanently give up ability to freeze
            if (request.setFlag !== undefined && request.setFlag > 0) {
              tx.SetFlag = request.setFlag;
            }
            if (request.clearFlag !== undefined && request.clearFlag > 0) {
              tx.ClearFlag = request.clearFlag;
            }
            break;

          case EnumTransactionType.TX_SINGER_LIST_SET:
            // SignerListSet transaction for multi-signature configuration
            // Reference: https://xrpl.org/docs/references/protocol/transactions/types/signerlistset
            // - SignerQuorum: Minimum weight of signatures required (0 to delete signer list)
            // - SignerEntries: Array of signers with their weights
            if (request.signerQuorum !== undefined) {
              tx.SignerQuorum = request.signerQuorum;
            }
            if (request.signerEntries && request.signerEntries.length > 0) {
              // Convert our SignerEntry format to XRPL SignerEntries format
              tx.SignerEntries = request.signerEntries.map((entry) => ({
                SignerEntry: {
                  Account: entry.account,
                  SignerWeight: entry.weight,
                },
              }));
            }
            // If signerQuorum is 0 and no signerEntries, this will remove the signer list
            break;

          default:
            // Other transaction types can be added here as needed
            break;
        }

        // Apply instructions if provided
        if (request.instructions) {
          if (request.instructions.fee) {
            tx.Fee = request.instructions.fee;
          }
          if (request.instructions.sequence !== undefined) {
            tx.Sequence = Number(request.instructions.sequence);
          }
          if (request.instructions.maxLedgerVersionOffset !== undefined) {
            const currentLedgerIndex = await client.getLedgerIndex();
            tx.LastLedgerSequence =
              currentLedgerIndex + Number(request.instructions.maxLedgerVersionOffset);
          } else if (request.instructions.maxLedgerVersion !== undefined) {
            tx.LastLedgerSequence = Number(request.instructions.maxLedgerVersion);
          }
        }

        // Use autofill to populate remaining fields
        // Cast to SubmittableTransaction as autofill accepts partial transactions
        const prepared = await client.autofill(tx as SubmittableTransaction);

        return {
          txJson: JSON.stringify(prepared),
          ...(request.instructions && { instructions: request.instructions }),
        };
      } catch (error) {
        if (error instanceof Error) {
          throw new PrepareTransactionError(
            `Failed to prepare transaction: ${error.message}`,
            (error as { data?: { error?: string } }).data?.error,
          );
        }
        throw new PrepareTransactionError('Unknown error preparing transaction');
      }
    },

    /**
     * Sign a transaction (offline operation)
     *
     * This operation does not require network connectivity.
     * Uses Wallet.fromSeed() to derive the signing keys.
     *
     * Security Note: The secret is used only for signing and is not logged.
     *
     * @param request - Transaction JSON and secret
     * @returns Signed transaction ID and blob
     * @throws SignTransactionError if signing fails
     */
    signTransaction: async (request: SignTransactionRequest): Promise<SignTransactionResponse> => {
      try {
        const wallet = Wallet.fromSeed(request.secret);
        const tx = JSON.parse(request.txJson);
        const signed = wallet.sign(tx);

        return {
          txId: signed.hash,
          txBlob: signed.tx_blob,
        };
      } catch (error) {
        if (error instanceof Error) {
          throw new SignTransactionError(`Failed to sign transaction: ${error.message}`);
        }
        throw new SignTransactionError('Unknown error signing transaction');
      }
    },

    /**
     * Submit a signed transaction to the network
     *
     * @param request - Signed transaction blob
     * @returns Submission result and earliest ledger version for validation check
     * @throws SubmitTransactionError if submission fails
     */
    submitTransaction: async (
      request: SubmitTransactionRequest,
    ): Promise<SubmitTransactionResponse> => {
      const client = await clientGetter();

      try {
        const latestLedgerVersion = await client.getLedgerIndex();
        const response = await client.submit(request.txBlob);

        return {
          resultJsonString: JSON.stringify(response.result),
          earliestLedgerVersion: BigInt(latestLedgerVersion + 1),
        };
      } catch (error) {
        if (error instanceof Error) {
          throw new SubmitTransactionError(
            `Failed to submit transaction: ${error.message}`,
            (error as { data?: { error?: string } }).data?.error,
          );
        }
        throw new SubmitTransactionError('Unknown error submitting transaction');
      }
    },

    /**
     * Get a transaction by ID
     *
     * Important: This method throws TransactionNotValidatedError if the
     * transaction has not yet been validated. Callers should use
     * waitValidation to wait for ledger advancement before retrying.
     *
     * @param request - Transaction ID and minimum ledger version
     * @returns Transaction result JSON
     * @throws TransactionNotValidatedError if transaction is not validated
     * @throws GetTransactionError for other failures
     */
    getTransaction: async (request: GetTransactionRequest): Promise<GetTransactionResponse> => {
      const client = await clientGetter();

      try {
        const response = await client.request({
          command: 'tx',
          transaction: request.txId,
          min_ledger: Number(request.minLedgerVersion),
        });

        // Check if the transaction is validated
        if (!response.result.validated) {
          throw new TransactionNotValidatedError(request.txId);
        }

        return {
          resultJsonString: JSON.stringify(response.result),
        };
      } catch (error) {
        // Re-throw TransactionNotValidatedError as-is
        if (error instanceof TransactionNotValidatedError) {
          throw error;
        }

        if (error instanceof Error) {
          throw new GetTransactionError(
            `Failed to get transaction ${request.txId}: ${error.message}`,
            (error as { data?: { error?: string } }).data?.error,
          );
        }
        throw new GetTransactionError(`Unknown error getting transaction ${request.txId}`);
      }
    },

    /**
     * Combine multiple signed transactions into a multi-signature transaction
     *
     * This operation does not require network connectivity.
     * Uses xrpl.multisign() to combine the signatures.
     *
     * @param request - Array of signed transaction blobs
     * @returns Combined signed transaction
     * @throws CombineTransactionError if combination fails
     */
    combineTransaction: async (
      request: CombineTransactionRequest,
    ): Promise<CombineTransactionResponse> => {
      // Validate input before calling external library
      if (request.signedTransactions.length < 2) {
        throw new CombineTransactionError(
          'At least 2 signed transactions are required for multi-signature',
        );
      }

      try {
        const combined = multisign(request.signedTransactions);

        return {
          signedTransaction: combined,
          txId: '', // The txId changes after combining, caller should compute if needed
        };
      } catch (error) {
        if (error instanceof Error) {
          throw new CombineTransactionError(`Failed to combine transactions: ${error.message}`);
        }
        throw new CombineTransactionError('Unknown error combining transactions');
      }
    },

    /**
     * Wait for ledger validation (server streaming)
     *
     * This is an async generator that yields ledger versions as they are validated.
     * The caller can iterate over the results until their transaction is validated.
     *
     * Important: This method explicitly subscribes to the ledger stream and
     * unsubscribes on cleanup. Always use try/finally or for-await-of to ensure
     * proper cleanup.
     *
     * @param context - Streaming context with abort signal
     * @yields Ledger versions as they are validated
     */
    waitValidation: async function* (
      context: StreamingContext,
    ): AsyncGenerator<WaitValidationResponse> {
      // Early exit if already aborted
      if (context.signal.aborted) {
        return;
      }

      const client = await clientGetter();

      // Subscribe to ledger events
      await client.request({ command: 'subscribe', streams: ['ledger'] });

      // Create a promise factory for ledger events with abort support
      const ledgerPromise = (): Promise<number | null> =>
        new Promise((resolve) => {
          // Handle abort signal
          const onAbort = () => {
            client.off('ledgerClosed', onLedger);
            resolve(null);
          };

          // Handle ledger event
          const onLedger = (ledger: { ledger_index: number }) => {
            context.signal.removeEventListener('abort', onAbort);
            resolve(ledger.ledger_index);
          };

          // Set up listeners
          context.signal.addEventListener('abort', onAbort, { once: true });
          client.once('ledgerClosed', onLedger);
        });

      try {
        while (!context.signal.aborted) {
          const ledgerIndex = await ledgerPromise();
          // Exit if aborted (ledgerIndex will be null)
          if (ledgerIndex === null) {
            break;
          }
          yield { ledgerVersion: BigInt(ledgerIndex) };
        }
      } finally {
        // Always unsubscribe when done
        await client.request({ command: 'unsubscribe', streams: ['ledger'] });
      }
    },
  };
}

// ============================================================================
// Default Service Instance
// ============================================================================

/**
 * RippleTransactionAPI service implementation
 *
 * Provides methods for XRP transaction operations.
 * Uses the singleton XRPL client for network connectivity.
 */
export const transactionService = createTransactionService(getClient);
