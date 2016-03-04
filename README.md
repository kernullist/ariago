# ariago
ARIA implementation with Go


# [개요]
- ARIA 구현 (based on KISA's impl)
- 현재 키 길이는 256bit만 지원
- 모드는 ECB 모드만 지원
- 키를 []byte, string 형태로 지정 가능


# [사용 예제]

```go
package main

import (
	"github.com/kernullist/ariago"	
	"fmt"
)

func main() {
	
	// Make a ariago instance with ECB mode
	aria := ariago.MakeARIA(ariago.ARIA_MODE_ECB)		
	if aria == nil {
		// Error
		return
	}
	
	// Master Key
	masterKey := "This is my Secret Key."
	if aria.SetMasterKeyWithString(masterKey) == false {
		// Error
		return
	}
	
	// Original Text
	originalText := "hello world!! this is a plain text"
	fmt.Printf("Original Text : %s\n", originalText)
	
	// Encrypt
	encryptedData := aria.Encrypt([]byte(originalText))
		
	// Decrypt	
	decryptedData := aria.Decrypt(encryptedData)
	
	fmt.Printf("Decrypted Data : %s\n", string(decryptedData))
}
```
