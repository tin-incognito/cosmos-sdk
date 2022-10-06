package models

import (
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
	"reflect"
	"testing"
)

func TestGeneratePrivateKey(t *testing.T) {
	tests := []struct {
		name string
		want key.PrivateKey
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeneratePrivateKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GeneratePrivateKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
