// ariago project ariago.go
package ariago

type AriaMode int

const (
	ARIA_MODE_ECB AriaMode = iota
	ARIA_MODE_MAX
)

const (
	ARIA_KEY_LENGTH = 32	// 256bit
	ARIA_BLOCK_SIZE = 16	// 128bit
)

const (
	max_round_key_length = (16 + 1) * 16
)


type ARIA struct {
	mode	 			AriaMode
	masterKey 		[ARIA_KEY_LENGTH]byte
	encryptRoundKey	[max_round_key_length]byte
	decryptRoundKey	[max_round_key_length]byte
}

func MakeARIA(mode AriaMode) *ARIA {
	if mode >= ARIA_MODE_MAX {
		return nil
	}
	
	aria := new(ARIA)
	
	if aria != nil {
		aria.mode = mode
		
		for i := 0; i < ARIA_KEY_LENGTH; i++ {
			aria.masterKey[i] = 0
		}
		
		for i := 0; i < max_round_key_length; i++ {
			aria.encryptRoundKey[i] = 0
			aria.decryptRoundKey[i] = 0
		}
	}
	
	return aria
}

func (aria *ARIA) SetMasterKeyWithByte(masterKey []byte) bool {	
	if len(masterKey) == 0 || len(masterKey) > ARIA_KEY_LENGTH {
		return false
	}
					
	copy(aria.masterKey[:], masterKey[:])

	aria.setupRoundKey()
	
	return true
}

func (aria *ARIA) SetMasterKeyWithString(masterKey string) bool {
	if len(masterKey) == 0 || len(masterKey) > ARIA_KEY_LENGTH {
		return false
	}
			
	for i, v := range masterKey {
		aria.masterKey[i] = byte(v)		
	}

	aria.setupRoundKey()
	
	return true
}

func (aria *ARIA) setupRoundKey() {
	encKeySetup(aria.masterKey[:ARIA_KEY_LENGTH], aria.encryptRoundKey[:max_round_key_length], ARIA_KEY_LENGTH * 8)	
	decKeySetup(aria.masterKey[:ARIA_KEY_LENGTH], aria.decryptRoundKey[:max_round_key_length], ARIA_KEY_LENGTH * 8)	
}

func (aria* ARIA) ariacrypt(inputData []byte, isEncrypt bool) []byte {
	blockCount := len(inputData) / ARIA_BLOCK_SIZE
	remainSize := len(inputData) % ARIA_BLOCK_SIZE
		
	outputDataSize := blockCount * ARIA_BLOCK_SIZE
	
	if remainSize > 0 {
		outputDataSize += ARIA_BLOCK_SIZE
	}
		
	var outputData []byte = make([]byte, outputDataSize)
	
	for i := 0; i < blockCount; i++ {
		if isEncrypt == true {
			crypt(inputData[i*ARIA_BLOCK_SIZE:(i+1)*ARIA_BLOCK_SIZE], ARIA_BLOCK_SIZE, aria.encryptRoundKey[:max_round_key_length], outputData[i*ARIA_BLOCK_SIZE:(i+1)*ARIA_BLOCK_SIZE])
		} else {
			crypt(inputData[i*ARIA_BLOCK_SIZE:(i+1)*ARIA_BLOCK_SIZE], ARIA_BLOCK_SIZE, aria.decryptRoundKey[:max_round_key_length], outputData[i*ARIA_BLOCK_SIZE:(i+1)*ARIA_BLOCK_SIZE])	
		}		
	}
	
	if remainSize > 0 {
		var tmpRemainData []byte = make([]byte, ARIA_BLOCK_SIZE)		
		for i := 0; i < ARIA_BLOCK_SIZE; i++ {
			tmpRemainData[i] = 0
		}
		
		copy(tmpRemainData, inputData[blockCount*ARIA_BLOCK_SIZE:])
		
		if isEncrypt == true {
			crypt(tmpRemainData, ARIA_BLOCK_SIZE, aria.encryptRoundKey[:max_round_key_length], outputData[blockCount*ARIA_BLOCK_SIZE:])
		} else {
			crypt(tmpRemainData, ARIA_BLOCK_SIZE, aria.decryptRoundKey[:max_round_key_length], outputData[blockCount*ARIA_BLOCK_SIZE:])
		}
	}
	
	return outputData
}

func (aria *ARIA) Encrypt(inputData []byte) []byte {
	return aria.ariacrypt(inputData, true)
}

func (aria *ARIA) Decrypt(inputData []byte) []byte {
	return aria.ariacrypt(inputData, false)
}

