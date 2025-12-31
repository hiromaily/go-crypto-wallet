const Web3 = require('web3');
const web3 = new Web3('http://127.0.0.1:7545');
web3.eth.handleRevert = true;

const contractAbi = require('../build/contracts/HyToken.json').abi;
const conAddress = '0x014C2061ba81a6Da4b8dD32b1322598D99B711D0';

const contract = new web3.eth.Contract(contractAbi, conAddress);

// const callUpdateCounter = async (owner, param) => {
//   const txObject = {
//     from: owner,
//     to: conAddress,
//     data: contract.methods.updateCounter(param).encodeABI(),
//   };
//   return await web3.eth.sendTransaction(txObject);
// };

const callBalanceOf = async (owner, param) => {
  const txObject = {
    from: owner,
    to: conAddress,
    data: contract.methods.balanceOf(param).encodeABI(),
  };
  return await web3.eth.call(txObject);
};

const main = async () => {
  const owner = '0xd0446b3eD62f23815bE724C17e34C0617A186e34';
  try {
    const txHash = await callBalanceOf(owner, owner);
    //const txHash = await callUpdateCounter(owner, 1);
    console.log(`tx hash: ${txHash}`);
  } catch (e) {
    console.log(`error in main: ${e}`);
    console.dir(e);
  }
  //when calling callReturnParam()
  //Error: Error: Returned error: execution reverted
};
main();
