package lightclient

// Broadcaster provides a way to send a signed transaction to a tendermint node
type Broadcaster interface {
	// Broadcast sends into to the chain
	// We only implement BroadcastCommit for now, add others???
	// The return result cannot be fully trusted without downloading signed headers
	Broadcast(tx []byte) (TmBroadcastResult, error)
}

// Checker provides access to calls to get data from the tendermint core
// and all cryptographic proof of its validity
type Checker interface {
	// Prove returns a merkle proof for the given key
	Prove(key []byte) (TmQueryResult, error)

	// SignedHeader gives us Header data along with the backing signatures,
	// so we can validate it externally (matching with the list of
	// known validators)
	SignedHeader(height uint64) (TmSignedHeader, error)
}

type Searcher interface {
	// Query gets data from the Blockchain state, possibly with a
	// complex path.  It doesn't worry about proofs
	Query(path string, data []byte) (TmQueryResult, error)
}

// Node represents someway to reliably read and write to the
// tendermint core (or a mock implenetation).
// TODO: is this interface even needed????
// Designed to be minimal, use the RPC directly for unverified info
type Node interface {
	Broadcaster
	Checker

	// Status and Validators give some info, nothing to be trusted though...
	// Unless we find that eg. the ValidatorResult matches the ValidatorHash
	// in a properly signed block header
	// Status() (TmStatusResult, error)
	// Validators() (TmValidatorResult, error)

	// TODO: let's make this reactive if possible
	// TODO: listen for a transaction being committed?
	// TODO: listen for a new block?
	// TODO: listen for change to a given key in the merkle store?

	// 	NetInfo() (*ctypes.ResultNetInfo, error)
	// 	DialSeeds(seeds []string) (*ctypes.ResultDialSeeds, error)
	// 	Genesis() (*ctypes.ResultGenesis, error)
	//  Block(height int) (*ctypes.ResultBlock, error)

	// 	ABCIInfo() (*ctypes.ResultABCIInfo, error)
}
