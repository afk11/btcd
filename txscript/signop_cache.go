package txscript

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"fmt"
)

// SigHashCache persists a map of SigHashType's to SigHash's
// during the execution of OP_CHECKSIG, OP_CHECKSIGVERIFY,
// OP_CHECKMULTISIG, OP_CHECKMULTISIGVERIFY. It's state is
// applicable only for a specific txin, and a specific
// signature opcode.
type SigHashCache struct {
	cache map[SigHashType]*chainhash.Hash
}

// Add associates the hashType with a ECDSA signature hash.
func (c *SigHashCache) Add(hashType SigHashType, hash *chainhash.Hash) {
	c.cache[hashType] = hash
}

// Contains will check whether the hash for this hashType is already cached.
func (c *SigHashCache) Contains(hashType SigHashType) bool {
	_, ok := c.cache[hashType]
	return ok
}

// Find will return the cached signature hash for hashType, or nil
// if none existed.
func (c *SigHashCache) Find(hashType SigHashType) *chainhash.Hash {
	hash, ok := c.cache[hashType]
	if !ok {
		return nil
	}
	return hash
}

// NewSigHashCache initializes a new SigHashCache to be used during a
// TxIns script evaluation.
func NewSigHashCache() *SigHashCache {
	return &SigHashCache{
		cache: make(map[SigHashType]*chainhash.Hash),
	}
}

// SignOp captures state about the type of checksig operation
type SignOp struct {
	op           int
	requiredSigs int
	pubkeys      []*btcec.PublicKey
	sigHashes    *SigHashCache
}

// IsCheckSig returns whether the operation was a OP_CHECKSIG-like
// operation, ie, OP_CHECKSIG{,VERIFY}
func (op *SignOp) IsCheckSig() bool {
	return op.op == OP_CHECKSIG || op.op == OP_CHECKSIGVERIFY
}

// IsCheckMultiSig returns whether the operation was a OP_CHECKMULTISIG-like
// operation, ie, OP_CHECKMULTISIG{,VERIFY}
func (op *SignOp) IsCheckMultiSig() bool {
	return op.op == OP_CHECKMULTISIG || op.op == OP_CHECKMULTISIGVERIFY
}

func (op *SignOp) GetKeys() []*btcec.PublicKey {
	return op.pubkeys
}

func (op *SignOp) GetSigHash(hashType SigHashType) []*btcec.PublicKey {
	if op.sigHashes.Contains(hashType) {
		return op.sigHashes.Find(hashType)
	}
	return nil
}

//func (op *SignOp) GetSignatures() map[*btcec.PublicKey] {
//
//}

// SignOpCache represents a sequence of signing operations executed
// during script evaluation
type SignOpCache struct {
	ops []*SignOp
}

// Helper function to create a SignOp from an opcode, whether
// it called the abstractVerify opcode, the number of required sigs,
// the cache of requested sigHashTypes+sigHashes, and the set
// of public keys
func newSignOp(op int, verify bool, requiredSigs int, cache *SigHashCache, keys []*btcec.PublicKey) *SignOp {
	if verify {
		op++
	}
	return &SignOp{op: op, requiredSigs: requiredSigs, pubkeys: keys, sigHashes: cache}
}

// CheckSig records a CheckSig operation into the cache(specifying whether it was *VERIFY)
func (opCache *SignOpCache) CheckSig(isVerify bool, cache *SigHashCache, keys []*btcec.PublicKey) {
	opCache.ops = append(opCache.ops, newSignOp(OP_CHECKSIG, isVerify, 1, cache, keys))
}

// CheckMultiSig records a CheckSig operation into the cache(specifying whether it was *VERIFY)
func (opCache *SignOpCache) CheckMultiSig(isVerify bool, cache *SigHashCache, requiredSigs int, keys []*btcec.PublicKey) {
	opCache.ops = append(opCache.ops, newSignOp(OP_CHECKMULTISIG, isVerify, requiredSigs, cache, keys))
}

// GetOp returns the SignOp at the specified idx, or an error
// if it did not exist.
func (opCache *SignOpCache) GetOp(idx int) (*SignOp, error) {
	if idx < 0 || idx > len(opCache.ops) {
		return nil, fmt.Errorf("No operation at %d", idx)
	}

	return opCache.ops[idx]
}

// NewSignOpCache returns an initialized SignOpCache.
func NewSignOpCache() *SignOpCache {
	return &SignOpCache{
		ops: make([]*SignOp, 0),
	}
}
