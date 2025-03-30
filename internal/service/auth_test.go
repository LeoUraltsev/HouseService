package service

import "testing"

func Test_generatePasswordHash(t *testing.T) {
	type args struct {
		pass string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "check pswd",
			args: args{
				pass: "qwerty",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, err := generatePasswordHash(tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("generatePasswordHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			r2, err := generatePasswordHash(tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("generatePasswordHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := r1 == r2

			if got != tt.want {
				t.Errorf("generatePasswordHash() = %v, want %v", got, tt.want)
			}

			if r1 == tt.args.pass {
				t.Errorf("generatePasswordHash() dont work")
			}
		})
	}
}
