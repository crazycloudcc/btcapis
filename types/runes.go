package types

// Runestone holds raw data extracted from a runestone encoded in an OP_RETURN
// output. The hex encoded payload is returned so that callers can perform
// higher level decoding if desired.
type Runestone struct {
	BodyHex string `json:"body_hex"`
}
