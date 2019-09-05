package config

import "testing"

func TestIsReservedWord(t *testing.T) {
	type args struct {
		applicationName string
	}
	tests := []struct {
		name            string
		applicationName string
		want            bool
	}{
		{name: "is not reserved", applicationName: "foo", want: false},
		{name: "is reserved", applicationName: "managed", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsReservedWord(tt.applicationName); got != tt.want {
				t.Errorf("IsReservedWord() = %v, want %v", got, tt.want)
			}
		})
	}
}
