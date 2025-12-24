package xrp

// XRPError is error object and used as error
type XRPError string

// XRPErrorDisabledAdminAPI is Admin API error
var XRPErrorDisabledAdminAPI XRPError = "Admin Method can not be used"

// Error returns error message
func (e XRPError) Error() string {
	return string(e)
}
