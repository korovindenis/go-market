package entity

import "testing"

func TestBalanceUpdate_IsValidNumber(t *testing.T) {
	tests := []struct {
		name    string
		b       *BalanceUpdate
		wantErr bool
	}{
		{
			name: "positive",
			b:    &BalanceUpdate{Order: "9278923470"},
		},
		{
			name:    "negative",
			b:       &BalanceUpdate{Order: "1"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.IsValidNumber(); (err != nil) != tt.wantErr {
				t.Errorf("BalanceUpdate.IsValidNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
