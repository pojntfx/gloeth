package encryptors

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"reflect"
	"testing"

	"github.com/pojntfx/gloeth/v3/pkg/wrappers"
)

func getKey() string {
	return "testtesttesttest"
}

func getFrame() [PlaintextFrameSize]byte {
	return [PlaintextFrameSize]byte{1}
}

func encryptFrame(key string, frame [PlaintextFrameSize]byte) ([wrappers.EncryptedFrameSize]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return [wrappers.EncryptedFrameSize]byte{}, err
	}

	encryptedFrame := [wrappers.EncryptedFrameSize]byte{}
	iv := encryptedFrame[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return [wrappers.EncryptedFrameSize]byte{}, nil
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	stream.XORKeyStream(encryptedFrame[aes.BlockSize:], frame[:])

	return encryptedFrame, nil
}

func decryptFrame(key string, frame [wrappers.EncryptedFrameSize]byte) ([PlaintextFrameSize]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return [PlaintextFrameSize]byte{}, err
	}

	iv := frame[:aes.BlockSize]
	encryptedFrame := frame[aes.BlockSize:]
	decryptedFrame := [PlaintextFrameSize]byte{}

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(decryptedFrame[:], encryptedFrame)

	return decryptedFrame, nil
}

func TestNewEthernet(t *testing.T) {
	expectedKey := getKey()

	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want *Ethernet
	}{
		{
			"New",
			args{
				expectedKey,
			},
			&Ethernet{
				key: expectedKey,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEthernet(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEthernet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEthernet_Encrypt(t *testing.T) {
	key := getKey()
	expectedFrame := getFrame()

	type fields struct {
		key string
	}
	type args struct {
		frame [PlaintextFrameSize]byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    [PlaintextFrameSize]byte
		wantErr bool
	}{
		{
			"Encrypt",
			fields{
				key,
			},
			args{
				expectedFrame,
			},
			expectedFrame,
			false,
		},
		{
			"Encrypt (faulty key)",
			fields{
				"",
			},
			args{
				expectedFrame,
			},
			[PlaintextFrameSize]byte{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ethernet{
				key: tt.fields.key,
			}

			encryptedFrame, err := e.Encrypt(tt.args.frame)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ethernet.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			actualFrame, err := decryptFrame(tt.fields.key, encryptedFrame)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}

			if !reflect.DeepEqual(actualFrame, tt.want) {
				t.Errorf("decrypt(Ethernet.Encrypt()) = %v, want %v", actualFrame, tt.want)
			}
		})
	}
}

func TestEthernet_Decrypt(t *testing.T) {
	key := getKey()
	expectedFrame := getFrame()
	frameIn, err := encryptFrame(key, expectedFrame)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		key string
	}
	type args struct {
		frame [wrappers.EncryptedFrameSize]byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    [PlaintextFrameSize]byte
		wantErr bool
	}{
		{
			"Decrypt",
			fields{
				key,
			},
			args{
				frameIn,
			},
			expectedFrame,
			false,
		},
		{
			"Decrypt (faulty key)",
			fields{
				"",
			},
			args{
				frameIn,
			},
			[PlaintextFrameSize]byte{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ethernet{
				key: tt.fields.key,
			}

			actualFrame, err := e.Decrypt(tt.args.frame)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ethernet.Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(actualFrame, tt.want) {
				t.Errorf("Ethernet.Decrypt() = %v, want %v", actualFrame, tt.want)
			}
		})
	}
}
