package manager

import (
	"testing"
)

func Test_getKeyHash(t *testing.T) {
	type args struct {
		values map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid base case",
			args: args{
				values: map[string]interface{}{
					"name":    "some_name",
					"surname": "surname",
					"gender":  "male",
					"sdf":     "gfd",
					"ertu":    "erti",
				},
			},
			want: "ae1ade8c230e9acc",
		},
		{
			name: "valid with changed order, hash should be the same",
			args: args{
				values: map[string]interface{}{
					"surname": "surname",
					"sdf":     "gfd",
					"ertu":    "erti",
					"gender":  "male",
					"name":    "some_name",
				},
			},
			want: "ae1ade8c230e9acc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getKeyHash(tt.args.values); got != tt.want {
				t.Errorf("getKeyHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
