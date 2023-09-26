package client

import (
	"testing"
)

func Test_addQueryInURI(t *testing.T) {
	type args struct {
		base   string
		params map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid base case",
			args: args{
				base: "/some_link/",
				params: map[string]string{
					"age":  "30",
					"name": "kirill",
				},
			},
			want:    "/some_link/?age=30&name=kirill",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addQueryInURI(tt.args.base, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("addQueryInURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("addQueryInURI() got = %v, want %v", got, tt.want)
			}
		})
	}
}
