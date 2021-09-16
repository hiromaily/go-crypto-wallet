package coin

// ERC20Token erc20 token
type ERC20Token string

// erc20_token
const (
	HYT ERC20Token = "hyt"
	BAT ERC20Token = "bat"
)

// String converter
func (e ERC20Token) String() string {
	return string(e)
}

// ERC20Map
var ERC20Map = map[ERC20Token]bool{
	HYT: true,
	BAT: true,
}

// IsERC20Token validate
func IsERC20Token(val string) bool {
	if _, ok := ERC20Map[ERC20Token(val)]; ok {
		return true
	}
	return false
}
