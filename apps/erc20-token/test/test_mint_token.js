const HyToken = artifacts.require('HyToken');

contract('HyToken', (accounts) => {
  let hyToken;

  before(async () => {
    hyToken = await HyToken.new();
  });

  describe('call mint()', async () => {
    it('call mint() by admin address', async () => {
      const setAmount = 10000;

      for (let idx = 0; idx < 5; idx++) {
        await hyToken.mint(accounts[idx], setAmount);
        const accountAmmount = await hyToken.balanceOf(accounts[idx]);
        assert.equal(
          accountAmmount,
          setAmount,
          'target account should own set amount'
        );
      }
    });

    // it('call mint() by NOT admin address', async () => {
    //   const setAmount = 10000;

    //   for (let idx = 0; idx < 5; idx++) {
    //     let err;
    //     await hyToken
    //       .mint(accounts[idx], setAmount, { from: accounts[1] })
    //       .catch((e) => {
    //         err = e;
    //         return;
    //       });
    //     assert.equal(
    //       err.reason,
    //       'Caller is not a admin',
    //       'err should not be undefined'
    //     );
    //   }
    // });
  });
});
