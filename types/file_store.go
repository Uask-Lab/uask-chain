package types

type StoreInfo struct {
	OnchainStore bool `json:"onchain_store"`
	// If OnchainStore is true, it is not nil.
	Content []byte `json:"content"`

	// If OnchainStore is false, these are not nil.
	// Hash is ipfs hash
	Hash string `json:"hash"`
	Url  string `json:"url"`
}
