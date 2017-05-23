package txscript

import (
	"encoding/hex"
	"testing"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func TestSignOpCache(t *testing.T) {
	params := &chaincfg.MainNetParams

	txHex := "0100000001d129bfacfb07e603b6d5dd160e2795686254684ff2991b9f7775e662c731cb0601000000b600493046022100eb9bb2e33ba279e8caeacb5101ae1e3bed77a6c88c771eaad4f2a1964d0c6d86022100ad4c6ff742db9949ff533dcd559b23694dd75a1a3c54b1b8962315065ac53775014c69522102951f8b7189e7096194cd2461b2d66b561d894ef2d36bd1bb0af86a2fa21fd3fb2102a3e85daaf647c8985727662b3a037c48db3cbe5236c32526fe4d506d9b55a34721032eda18a391eb3db1812810836668980469c02b858f5df4bcf15114b06c5b619453aeffffffff02a08601000000000017a91487148c0201c58fb7223b14c8b7d81443d37c418f8750eb0b040000000017a914a18ee4fb6a3e673b1a41f0710b0dd05ca6483d198700000000"
	serializedTx, err := hex.DecodeString(txHex)
	if err != nil {
		panic(err)
	}

	tx, err := btcutil.NewTxFromBytes(serializedTx)
	if err != nil {
		t.Errorf("bad test %s", err)
	}

	rs, err := hex.DecodeString("522102951f8b7189e7096194cd2461b2d66b561d894ef2d36bd1bb0af86a2fa21fd3fb2102a3e85daaf647c8985727662b3a037c48db3cbe5236c32526fe4d506d9b55a34721032eda18a391eb3db1812810836668980469c02b858f5df4bcf15114b06c5b619453ae")
	if err != nil {
		panic(err)
	}

	addr, err := btcutil.NewAddressScriptHash(rs, params)
	if err != nil {
		panic(err)
	}

	scriptPubKey, err := PayToAddrScript(addr)
	if err != nil {
		panic(err)
	}

	sigCache := NewSigCache(10)
	flags := ScriptFlags(0)
	nIn := 0
	e, err := NewEngine(scriptPubKey, tx.MsgTx(), nIn, flags, sigCache, nil, 0)
	if err != nil {
		panic(err)
	}

	signOpCache := NewSignOpCache()
	e.SetSignOpCache(signOpCache)

	err = e.Execute()
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(signOpCache.ops); i++ {
		op, _ := signOpCache.GetOp(i)

	}

}