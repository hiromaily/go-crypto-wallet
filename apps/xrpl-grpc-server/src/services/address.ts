/**
 * RippleAddressAPI Service
 *
 * Implements address generation and validation methods for the XRP Ledger.
 * Uses xrpl.js 4.5.0 Wallet and validation utilities.
 *
 * Security Note: Generated secrets must never be logged.
 */

import { Wallet, isValidAddress as xrplIsValidAddress } from 'xrpl';

/**
 * Response for GenerateAddress RPC
 */
export interface GenerateAddressResponse {
  xAddress: string;
  classicAddress: string;
  address: string;
  secret: string;
}

/**
 * Response for GenerateXAddress RPC
 */
export interface GenerateXAddressResponse {
  xAddress: string;
  secret: string;
}

/**
 * Request for IsValidAddress RPC
 */
export interface IsValidAddressRequest {
  address: string;
}

/**
 * Response for IsValidAddress RPC
 */
export interface IsValidAddressResponse {
  isValid: boolean;
}

/**
 * RippleAddressAPI service implementation
 *
 * Provides methods for XRP address generation and validation.
 * These methods do not require network connectivity as they operate
 * on cryptographic primitives only.
 */
export const addressService = {
  /**
   * Generate a new XRP address with all address formats
   *
   * Creates a new random wallet and returns:
   * - xAddress: X-address format (includes destination tag support)
   * - classicAddress: Classic r... address format
   * - address: Same as classicAddress (for compatibility)
   * - secret: The seed for signing (NEVER log this)
   *
   * @returns Generated address information with secret
   */
  generateAddress: (): GenerateAddressResponse => {
    const wallet = Wallet.generate();
    if (!wallet.seed) {
      // This should not happen with Wallet.generate(), but as a safeguard
      throw new Error('Failed to generate a wallet seed.');
    }

    return {
      xAddress: wallet.getXAddress(),
      classicAddress: wallet.classicAddress,
      address: wallet.address,
      secret: wallet.seed,
    };
  },

  /**
   * Generate a new XRP X-address only
   *
   * Creates a new random wallet and returns only the X-address format.
   * X-addresses encode the classic address and optional destination tag
   * into a single address string.
   *
   * @returns Generated X-address with secret
   */
  generateXAddress: (): GenerateXAddressResponse => {
    const wallet = Wallet.generate();
    if (!wallet.seed) {
      // This should not happen with Wallet.generate(), but as a safeguard
      throw new Error('Failed to generate a wallet seed.');
    }

    return {
      xAddress: wallet.getXAddress(),
      secret: wallet.seed,
    };
  },

  /**
   * Validate an XRP address
   *
   * Checks if the provided address is a valid XRP address.
   * Supports both classic (r...) and X-address formats.
   *
   * @param request - Contains the address to validate
   * @returns Whether the address is valid
   */
  isValidAddress: (request: IsValidAddressRequest): IsValidAddressResponse => {
    return {
      isValid: xrplIsValidAddress(request.address),
    };
  },
};
