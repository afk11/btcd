package btcec_test

import (
	"encoding/hex"
	"gopkg.in/yaml.v2"
	"testing"
	"time"
	"log"
	"io/ioutil"
	"strings"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
)

func readFile(filename string) []byte {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return source
}

func removeSigHash(sig string) string {
	return strings.TrimSuffix(sig, "01")
}

type EcdsaTestCase struct {
	PrivateKey string `yaml:"privkey"`
	Message    string `yaml:"msg"`
	Sig        string `yaml:"sig"`
}

func (t *EcdsaTestCase) GetPrivateKey() []byte {
	private, err := hex.DecodeString(t.PrivateKey)
	if err != nil {
		panic("Invalid private key")
	}
	return private
}

func (t *EcdsaTestCase) GetMessage() []byte {
	msg, err := hex.DecodeString(t.Message)
	if err != nil {
		panic("Invalid msg32")
	}
	return msg
}

func (t *EcdsaTestCase) GetSigBytes() []byte {
	sig, err := hex.DecodeString(removeSigHash(t.Sig))
	if err != nil {
		panic("Invalid msg32")
	}
	return sig
}

type EcdsaFixtures []EcdsaTestCase

func GetEcdsaFixtures() []EcdsaTestCase {
	source := readFile("sign_vectors.yaml")
	testCase := EcdsaFixtures{}
	err := yaml.Unmarshal(source, &testCase)
	if err != nil {
		panic(err)
	}
	return testCase
}

func Benchmark_Ecdsa_Sign(t *testing.B) {

	fixtures := GetEcdsaFixtures()

	sum := 0
	for j := 0; j < t.N; j++ {
		for i := 0; i < len(fixtures); i++ {
			testCase := fixtures[i]
			msg32 := testCase.GetMessage()

			privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), testCase.GetPrivateKey())

			// Sign a message using the private key.
			S := time.Now()
			_, err := privKey.Sign(msg32)
			dur := time.Since(S)
			log.Println(dur.Nanoseconds())
			sum = sum + int(dur)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	avg := sum / len(fixtures)
	log.Printf("Sign time: %d when run %d times ", avg, len(fixtures))
}
