package pdfutil

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_pdfReader_scanTxt(t *testing.T) {
	buffer, err := getPdfBytes("/home/oyamo/Documents/Interruptions - 01.04.2021 (Upload).pdf", 0)
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

			got, err := scanTxt(tt.args.buffer)
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

func TestScanPDF(t *testing.T) {
	path := "/home/oyamo/.cache/3339f6b2-8c3b-4dd0-9f07-581405a0154f.pdf"
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *BlackoutResult
		wantErr bool
	}{
		{
			name:    "Parse Test",
			args:    args{path: path},
			wantErr: false,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ScanPDF(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanPDF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanPDF() got = %v, want %v", got, tt.want)
			}
		})
	}
}
