/**
 * Services Module Exports
 *
 * Aggregates all gRPC service implementations.
 */

export {
  AccountInfoError,
  AccountNotFoundError,
  accountService,
  createAccountService,
  type GetAccountInfoRequest,
  type GetAccountInfoResponse,
} from './account';

export {
  addressService,
  type GenerateAddressResponse,
  type GenerateXAddressResponse,
  type IsValidAddressRequest,
  type IsValidAddressResponse,
} from './address';
