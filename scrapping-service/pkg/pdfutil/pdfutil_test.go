package pdfutil

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_pdfReader_scanTxt(t *testing.T) {
	r := &pdfReader{}
	buffer, err := r.readPdf("/home/oyamo/Documents/Interruptions - 01.04.2021 (Upload).pdf")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	type args struct {
		buffer bytes.Buffer
	}
	tests := []struct {
		name    string
		args    args
		want    *BlackoutResult
		wantErr bool
	}{
		{
			args: args{
				buffer: buffer,
			},
			want:    nil,
			wantErr: false,
			name:    "Parser test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := r.scanTxt(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanTxt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scanTxt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
