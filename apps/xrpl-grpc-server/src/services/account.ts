/**
 * RippleAccountAPI Service
 *
 * Implements account information methods for the XRP Ledger.
 * Uses xrpl.js 4.5.0 to retrieve validated account data.
 *
 * Important: Always uses ledger_index: "validated" to ensure finalized data.
 */

import { type Client, dropsToXrp } from 'xrpl';
import { getClient } from '../xrpl';

/**
 * Request for GetAccountInfo RPC
 */
export interface GetAccountInfoRequest {
  address: string;
}

/**
 * Response for GetAccountInfo RPC
 */
export interface GetAccountInfoResponse {
  sequence: bigint;
  xrpBalance: string;
  ownerCount: bigint;
  previousAffectingTransactionId: string;
  previousAffectingTransactionLedgerVersion: bigint;
}

/**
 * Error thrown when an account is not found on the ledger
 */
export class AccountNotFoundError extends Error {
  constructor(address: string) {
    super(`Account not found: ${address}`);
    this.name = 'AccountNotFoundError';
  }
}

/**
 * Error thrown when account info request fails
 */
export class AccountInfoError extends Error {
  constructor(
    message: string,
    public readonly code?: string,
  ) {
    super(message);
    this.name = 'AccountInfoError';
  }
}

/**
 * Create the account service with a client getter function
 *
 * This factory pattern allows for dependency injection in tests.
 *
 * @param clientGetter - Function to get the XRPL client
 * @returns Account service implementation
 */
export function createAccountService(clientGetter: () => Promise<Client>) {
  return {
    /**
     * Get account information from the XRP Ledger
     *
     * Retrieves validated account data including:
     * - sequence: The current sequence number for transactions
     * - xrpBalance: The account's XRP balance (converted from drops)
     * - ownerCount: Number of objects owned by this account
     * - previousAffectingTransactionId: Hash of last transaction affecting this account
     * - previousAffectingTransactionLedgerVersion: Ledger index of that transaction
     *
     * @param request - Contains the address to query
     * @returns Account information
     * @throws AccountNotFoundError if the account does not exist
     * @throws AccountInfoError for other API errors
     */
    getAccountInfo: async (request: GetAccountInfoRequest): Promise<GetAccountInfoResponse> => {
      const client = await clientGetter();

      try {
        const response = await client.request({
          command: 'account_info',
          account: request.address,
          ledger_index: 'validated', // Important: always use validated ledger for finalized data
        });

        const accountData = response.result.account_data;

        return {
          sequence: BigInt(accountData.Sequence),
          xrpBalance: String(dropsToXrp(accountData.Balance)),
          ownerCount: BigInt(accountData.OwnerCount),
          previousAffectingTransactionId: accountData.PreviousTxnID ?? '',
          previousAffectingTransactionLedgerVersion: BigInt(accountData.PreviousTxnLgrSeq ?? 0),
        };
      } catch (error) {
        // Handle specific XRPL errors
        if (error instanceof Error) {
          // RippledError has a `data` property with the error code.
          // See: https://xrpl.org/docs/references/http-websocket-apis/error-formatting/
          const errorCode = (error as { data?: { error?: string } }).data?.error;

          // actNotFound is the error code for non-existent accounts
          if (errorCode === 'actNotFound') {
            throw new AccountNotFoundError(request.address);
          }

          // Re-throw with more context
          throw new AccountInfoError(
            `Failed to get account info for ${request.address}: ${error.message}`,
            errorCode,
          );
        }

        throw new AccountInfoError(`Unknown error getting account info for ${request.address}`);
      }
    },
  };
}

/**
 * RippleAccountAPI service implementation
 *
 * Provides methods for XRP account information retrieval.
 * Uses the singleton XRPL client for network connectivity.
 */
export const accountService = createAccountService(getClient);
