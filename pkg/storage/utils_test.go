package storage

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Test_getFilesFromFolder(t *testing.T) {
	type args struct {
		folderPath string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFilesFromFolder(tt.args.folderPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFilesFromFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFilesFromFolder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTotalVersionedSize(t *testing.T) {
	type args struct {
		bucket string
		client *s3.Client
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTotalVersionedSize(tt.args.bucket, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTotalVersionedSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getTotalVersionedSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
