package transferring

import (
	"testing"
)

func Test_service_BetweenAccounts(t *testing.T) {
	type args struct {
		originBalance      float64
		destinationBalance float64
		amount             float64
	}
	tt := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "When transfer occurs successfully",
			args: args{
				originBalance:      112.52,
				destinationBalance: 23.29,
				amount:             112.51,
			},
			wantErr: false,
		},
		{
			name: "When there's not enough balance",
			args: args{
				originBalance:      112.52,
				destinationBalance: 23.29,
				amount:             112.53,
			},
			wantErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := NewService()
			newOBalance, newDBalance, err := s.BalanceBetweenAccounts(tc.args.originBalance, tc.args.destinationBalance, tc.args.amount)
			tDBalance, tOBalance, _ := s.BalanceBetweenAccounts(newDBalance, newOBalance, tc.args.amount)

			if (err != nil) && !tc.wantErr {
				t.Errorf("BalanceBetweenAccounts() error = %v; wantErr = %v", err, tc.wantErr)
			}

			if !tc.wantErr && tOBalance != tc.args.originBalance {
				t.Errorf("Expected reverse operation lead to equality to originBalance %f, but got %f", tc.args.originBalance, tOBalance)
			}

			if !tc.wantErr && tDBalance != tc.args.destinationBalance {
				t.Errorf("Expected reverse operation lead to equality to destinationBalance %f, but got %f", tc.args.destinationBalance, tDBalance)
			}
		})
	}
}
