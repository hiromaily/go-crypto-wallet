import { Command } from 'commander';
import { AbiItem } from 'web3-utils';
import HyToken from '../build/contracts/HyToken.json';
import { ERC20 } from './erc20/erc20';

const program = new Command();

interface Web3Params {
  nodeURL: string;
  contractAddr: string;
  ownerAddr: string;
  mode: string;
  targetAddr?: string;
  amount?: number;
}

const checkArgs = (): Web3Params | undefined => {
  program
    .option('-m, --mode <string>', 'balance, transfer')
    .option('-a, --address <string>', 'target address')
    .option('-u, --amount <number>', 'amount')
    .parse(process.argv);
  const opts = program.opts();

  const params = ['mode'];
  for (const param of params) {
    if (!opts[param]) {
      console.error(`${param} option is required`);
      return undefined;
    }
  }

  const envParams = ['NODE_URL', 'CONTRACT_ADDRESS', 'OWNER_ADDRESS'];
  for (const param of envParams) {
    if (!process.env[param]) {
      console.error(`${param} environment variable is required`);
      return undefined;
    }
  }

  const web3Params: Web3Params = {
    nodeURL: process.env.NODE_URL || '',
    contractAddr: process.env.CONTRACT_ADDRESS || '',
    ownerAddr: process.env.OWNER_ADDRESS || '',
    mode: opts.mode,
  };
  if (opts.address) web3Params.targetAddr = opts.address;
  if (opts.amount) web3Params.amount = opts.amount;
  return web3Params;
};

const main = async (): Promise<void> => {
  const args = checkArgs();
  if (args === undefined) throw 'args is invalid';

  try {
    const contractAbi: AbiItem[] = HyToken.abi as AbiItem[];
    const erc20: ERC20 = new ERC20(
      args.nodeURL,
      args.contractAddr,
      contractAbi
    );

    console.log(`command: ${args.mode}`);
    switch (args.mode) {
      case 'balance': {
        // validate args
        if (!args.targetAddr) throw new Error('--address option is required');
        const hexBalance: string = await erc20.callBalanceOf(
          args.ownerAddr,
          args.targetAddr
        );
        console.log(`balance: ${parseInt(hexBalance, 16)}`);
        break;
      }
      case 'transfer': {
        // validate args
        if (!args.targetAddr) throw new Error('--address option is required');
        if (!args.amount) throw new Error('--amount option is required');
        const resultJSON = await erc20.callTransfer(
          args.ownerAddr,
          args.targetAddr,
          args.amount
        );
        console.log('result:', resultJSON);
        break;
      }
      case 'estimateGas': {
        const gas = await erc20.callEstimateGas(args.ownerAddr);
        console.log(gas);
      }
    }
  } catch (e) {
    console.error(e);
  }
};

void main();
