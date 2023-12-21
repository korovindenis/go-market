package entity

import "testing"

func TestOrder_IsValidNumber(t *testing.T) {
	tests := []struct {
		name    string
		o       *Order
		wantErr bool
	}{
		{
			name: "positive",
			o:    &Order{Number: "9278923470"},
		},
		{
			name:    "negative",
			o:       &Order{Number: "1"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.IsValidNumber(); (err != nil) != tt.wantErr {
				t.Errorf("Order.IsValidNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
