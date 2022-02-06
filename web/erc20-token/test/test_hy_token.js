const testHyToken = artifacts.require('testHyToken');

/*
 * uncomment accounts to access the test accounts made available by the
 * Ethereum client
 * See docs: https://www.trufflesuite.com/docs/truffle/testing/writing-tests-in-javascript
 */
contract('testHyToken', function (/* accounts */) {
  it('should assert true', async function () {
    await testHyToken.deployed();
    return assert.isTrue(true);
  });
});
