package bitcoin

// Version represents Bitcoin Core version
type Version int

// Bitcoin Core versions
const (
	Ver17 Version = 170000
	Ver18 Version = 180000
	Ver19 Version = 190000
	Ver20 Version = 200000
	Ver21 Version = 210000 // for BCH version
)

// Int converts Version to int
func (v Version) Int() int {
	return int(v)
}

// RequiredVersion is the minimum required version for this wallet
const RequiredVersion = Ver19
