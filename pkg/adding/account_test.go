package adding

import "testing"

func TestErrInvalidAccountField_Error(t *testing.T) {
	tests := []struct {
		name      string
		e         ErrInvalidAccountField
		want      string
		willPanic bool
	}{
		{
			name: "When error with field and message",
			e: ErrInvalidAccountField{
				field:   "cpf",
				message: "foo bar",
			},
			want:      "Field cpf contains an invalid value: foo bar",
			willPanic: false,
		},
		{
			name: "When error with message but no field",
			e: ErrInvalidAccountField{
				message: "foo bar",
			},
			want:      "",
			willPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.willPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Error() did not panic")
					}
				}()
			}
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_name_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		n       name
		input   []byte
		wantErr bool
	}{
		{
			name:    "When runs smoothly",
			n:       "",
			input:   []byte(`"Chewbacca Solo"`),
			wantErr: false,
		},
		{
			name:    "When name is an empty string",
			n:       "",
			input:   []byte(`""`),
			wantErr: true,
		},
		{
			name:    "When name is not a string",
			n:       "",
			input:   []byte(`56416`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.UnmarshalJSON(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cpf_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		c       cpf
		input   []byte
		wantErr bool
	}{
		{
			name:    "When CPF is valid",
			c:       "",
			input:   []byte(`"11111111030"`),
			wantErr: false,
		},
		{
			name:    "When CPF is invalid with expected length",
			c:       "",
			input:   []byte(`"11111111111""`),
			wantErr: true,
		},
		{
			name:    "When CPF is invalid with lower length",
			c:       "",
			input:   []byte(`"11111111111"`),
			wantErr: true,
		},
		{
			name:    "When CPF is invalid with higher length",
			c:       "",
			input:   []byte(`"11111111111"`),
			wantErr: true,
		},
		{
			name:    "When CPF is not a string",
			c:       "",
			input:   []byte("11111111111"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.UnmarshalJSON(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_secret_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		s       secret
		input   []byte
		wantErr bool
	}{
		{
			name:    "When runs smoothly",
			s:       nil,
			input:   []byte(`"254855"`),
			wantErr: false,
		},
		{
			name:    "When input is not of string type",
			s:       nil,
			input:   []byte("123456"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.UnmarshalJSON(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_balance_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		b       balance
		input   []byte
		wantErr bool
	}{
		{
			name:    "When runs smoothly",
			b:       0,
			input:   []byte(`42.42`),
			wantErr: false,
		},
		{
			name:    "When input is a negative number",
			b:       0,
			input:   []byte(`-42.42`),
			wantErr: true,
		},
		{
			name:    "When input is not of numeric type",
			b:       0,
			input:   []byte(`"abc"`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.UnmarshalJSON(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
