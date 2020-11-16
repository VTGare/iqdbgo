package iqdbgo

import (
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    *Result
		wantErr bool
	}{
		{"Test", args{"https://i.pinimg.com/originals/eb/bb/9c/ebbb9c6067fd1f30c0c1b5261833e051.jpg"}, &Result{}, false},
		{"Test bad", args{"https://pbs.twimg.com/media/Eivr_7PUcAAgwYp.jpg"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Search(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
