import Web3 from 'web3';
import abi from 'ethereumjs-abi';
//import { TransactionConfig } from 'web3-core';
import { TransactionReceipt, TransactionConfig } from 'web3-core';
import { AbiItem } from 'web3-utils';
import { Contract } from 'web3-eth-contract';

export class ERC20 {
  private _web3: Web3;
  private _contractAddr: string;
  private _contract: Contract;

  public constructor(
    nodeURL: string,
    contractAddr: string,
    contractAbi: AbiItem[]
  ) {
    this._web3 = new Web3(nodeURL);
    this._web3.eth.handleRevert = true;

    this._contractAddr = contractAddr;
    this._contract = new this._web3.eth.Contract(
      contractAbi,
      this._contractAddr
    );
  }

  public async callBalanceOf(owner: string, account: string): Promise<string> {
    const txObject: TransactionConfig = {
      from: owner,
      to: this._contractAddr,
      data: this._contract.methods.balanceOf(account).encodeABI() as string,
    } as TransactionConfig;
    return await this._web3.eth.call(txObject);
  }

  public async callTransfer(
    from: string,
    to: string,
    amount: number
  ): Promise<TransactionReceipt> {
    const txObject: TransactionConfig = {
      from: from,
      to: this._contractAddr,
      data: this._contract.methods.transfer(to, amount).encodeABI() as string,
    } as TransactionConfig;
    return await this._web3.eth.sendTransaction(txObject);
  }

  public async callTransferFrom(
    from: string,
    to: string,
    amount: number
  ): Promise<TransactionReceipt> {
    const txObject: TransactionConfig = {
      from: from,
      to: this._contractAddr,
      data: this._contract.methods
        .transferFrom(from, to, amount)
        .encodeABI() as string,
    } as TransactionConfig;
    return await this._web3.eth.sendTransaction(txObject);
  }

  // just temporary implementation because golang cli doesn't work...
  public async callEstimateGas(from: string): Promise<number> {
    const encodedData: string =
      '0x' +
      abi
        .simpleEncode('transfer(address,uint256)', this._contractAddr, 100)
        .toString('hex');

    const txObject: TransactionConfig = {
      from: from,
      to: this._contractAddr,
      data: encodedData,
    } as TransactionConfig;
    const gas = await this._web3.eth.estimateGas(txObject);
    return gas;
  }
}
