package coin

// ERC20Token erc20 token
type ERC20Token string

// erc20_token
const (
	BAT ERC20Token = "bat"
)

// String converter
func (e ERC20Token) String() string {
	return string(e)
}

// ERC20Address
var ERC20Address = map[ERC20Token]string{
	BAT: "0x0D8775F648430679A709E98d2b0Cb6250d2887EF",
}

// IsERC20Token validate
func IsERC20Token(val string) bool {
	if _, ok := ERC20Address[ERC20Token(val)]; ok {
		return true
	}
	return false
}
