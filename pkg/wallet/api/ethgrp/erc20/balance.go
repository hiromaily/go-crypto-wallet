package erc20

//func (c *Contract) GetBalance(hexAddr string, _ eth.QuantityTag) (*big.Int, error) {
//
//	balance, err := c.token.BalanceOf(common.HexToAddress(hexAddr)) // (*big.Int, error)
//	if err != nil {
//		return nil, errors.Wrapf(err, "fail to call e.contract.BalanceOf(%s)", hexAddr)
//	}
//	h, err := hexutil.DecodeBig(balance)
//	if err != nil {
//		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
//	}
//
//	return h, nil
//}
