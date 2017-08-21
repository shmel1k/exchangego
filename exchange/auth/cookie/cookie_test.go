package cookie

import (
	"encoding/hex"
	"testing"
)

func TestDecryptCookie(t *testing.T) {
	var testData = []struct {
		testUser     string
		expectedUser string
		cookie       string
		err          bool
	}{
		{
			testUser:     "petrusha",
			expectedUser: "petrusha",
			cookie:       "",
			err:          false,
		},
	}

	for _, v := range testData {
		if v.cookie == "" {
			cook, err := GenerateCookie(v.testUser)
			if err != nil {
				t.Errorf("got err %q while generating cookie", err)
			}
			v.cookie = cook
		}
		unpacked, err := hex.DecodeString(v.cookie)
		if err != nil {
			t.Errorf("got err while decoding cookie: %s", err)
		}
		user, err := decryptCookie([]byte(unpacked))
		if err != nil && !v.err {
			t.Errorf("got error %q, want nil", err)
		}
		if user != v.expectedUser {
			t.Errorf("got user %q, expected %q", user, v.expectedUser)
		}
	}

}
