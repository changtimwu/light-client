package sign

import (
	"github.com/pkg/errors"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
	lightclient "github.com/tendermint/light-client"
)

func init() {
	registerGoCryptoGoWire()
}

// we must register these types here, to make sure they parse (maybe go-wire issue??)
// TODO: fix go-wire, remove this code
func registerGoCryptoGoWire() {
	wire.RegisterInterface(
		struct{ crypto.PubKey }{},
		wire.ConcreteType{O: crypto.PubKeyEd25519{}, Byte: crypto.PubKeyTypeEd25519},
		wire.ConcreteType{O: crypto.PubKeySecp256k1{}, Byte: crypto.PubKeyTypeSecp256k1},
	)
	wire.RegisterInterface(
		struct{ crypto.Signature }{},
		wire.ConcreteType{O: crypto.SignatureEd25519{}, Byte: crypto.SignatureTypeEd25519},
		wire.ConcreteType{O: crypto.SignatureSecp256k1{}, Byte: crypto.SignatureTypeSecp256k1},
	)
}

func Single(data []byte) lightclient.Signable {
	return &single{data: data}
}

type single struct {
	data   []byte
	sig    []byte
	pubkey []byte
}

func (s *single) Bytes() []byte {
	return s.data
}

func (s *single) Sign(addr, pubkey, sig []byte) error {
	if len(sig) == 0 || len(pubkey) == 0 {
		return errors.New("Signature or Key missing")
	}
	if s.sig != nil {
		return errors.New("Transaction can only be signed once")
	}
	s.sig = sig
	s.pubkey = pubkey
	return nil
}

func (s *single) Signed() ([]byte, error) {
	if s.sig == nil {
		return nil, errors.New("Transaction was never signed")
	}
	return wire.BinaryBytes(s), nil
}

// TODO: how do we verify this??? need some function to deserialize and verify the sigs!

// // Validate will deserialize the contained action, and validate the signature or return an error
// func (tx SignedAction) Validate() (ValidatedAction, error) {
// 	res := ValidatedAction{
// 		SignedAction: tx,
// 	}
// 	valid := tx.Signer.VerifyBytes(tx.ActionData, tx.Signature)
// 	if !valid {
// 		return res, errors.New("Invalid signature")
// 	}

// 	var err error
// 	res.action, err = ActionFromBytes(tx.ActionData)
// 	if err == nil {
// 		res.valid = true
// 	}
// 	return res, err
// }

// // SignAction will serialize the action and sign it with your key
// func SignAction(action Action, privKey crypto.PrivKey) (res SignedAction, err error) {
// 	res.ActionData, err = ActionToBytes(action)
// 	if err != nil {
// 		return res, err
// 	}
// 	res.Signature = privKey.Sign(res.ActionData)
// 	res.Signer = privKey.PubKey()
// 	return res, nil
// }
