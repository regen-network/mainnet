package main

import (
	"reflect"
	"testing"
)

func TestMergeAccounts(t *testing.T) {
	tests := []struct {
		name     string
		accounts []Account
		want     []Account
		wantErr  bool
	}{
		// TODO
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeAccounts(tt.accounts)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeAccounts() got = %v, want %v", got, tt.want)
			}
		})
	}
}
