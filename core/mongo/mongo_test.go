package mongo

import (
	"context"
	"testing"
)

func TestMustNewClient(t *testing.T) {
	type args struct {
		uri string
		db  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "one",
			args: args{
				uri: "mongodb://127.0.0.1:27017/test?maxPoolSize=100&retryWrites=true&w=majority&connect=direct",
				db:  "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := MustNewClient(tt.args.uri, tt.args.db)
			t.Log(client.Close(context.TODO()))
		})
	}
}
