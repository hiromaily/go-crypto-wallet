package bitcoin

// Version represents Bitcoin Core version
type Version int

// Bitcoin Core versions
const (
	Ver28 Version = 280000
	Ver29 Version = 290000
	Ver30 Version = 300000
)

// Int converts Version to int
func (v Version) Int() int {
	return int(v)
}

// RequiredVersion is the minimum required version for this wallet
const RequiredVersion = Ver28
