package u2u

import (
	"fmt"
	"testing"
)

var privateKeyUsedForTest = "1b734ae16eb3b7470d99780dff19bc7e2d8ce5b04785a7390d7363e78d37c6e8"

func TestSignText(t *testing.T) {
	u2uNode := NewU2U(nil)
	if err := u2uNode.SetAccount(privateKeyUsedForTest); err != nil {
		t.Fatal(err)
	}

	signature, err := u2uNode.SignText([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("signature %x\n", signature)
}
