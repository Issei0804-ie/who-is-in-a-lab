package util

import "testing"

func TestValidateMacAddress(t *testing.T) {
	type args struct {
		macAddress string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "正常系",
			args: args{macAddress: "b8:85:84:33:7a:61"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateMacAddress(tt.args.macAddress); got != tt.want {
				t.Errorf("ValidateMacAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
