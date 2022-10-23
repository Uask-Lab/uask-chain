package types

type StoreInfo struct {
	Content []byte
	// Hash is ipfs hash
	Hash string `json:"hash"`
	// Url is ipfs node URL
	Url string `json:"url"`
}
