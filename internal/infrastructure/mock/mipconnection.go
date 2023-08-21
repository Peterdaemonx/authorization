package mock

import (
	"encoding/hex"
)

type MipConnection struct {
}

// Mock response for a successful authorization from mastercard to be mapper into a AuthorizationResult entity
const b = "f0f1f1f0766300018601a002f1f6f2f2f2f3f0f0f1f7f6f0f0f0f2f7f0f4f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f1f0f1f0f0f0f0f0f0f0f0f0f1f2f1f0f3f2f9f1f4f0f4f2f7f6f1f2f0f1f4f9f0f1f0f0f0f0f1f0f3f2f9f0f3f2f9f0f6f0f0f0f0f0f0f0f6f0f0f0f0f0f0f0f0f5f9f6f6f0f0f0f9f3e3f4f2f0f7f0f1f0f3f2f1f2f4f3f2f891c9f3d1c2929281d8f197f8c3c2c1c1c1c2a8f0c3c8e4c1c1c1c17ef6f6f4f5f0f1f0f1f2f0f2f3f6f38284f2f1f3f78460f0f886f160f4868582608281f5f060f383f284f4f4f0f183f9f181f9f7f8f8f4f0f0f0f9d4c1c2f0f1f1f0f6c8"

func (m MipConnection) Send(_ []byte) ([]byte, error) {
	return hex.DecodeString(b)
}
