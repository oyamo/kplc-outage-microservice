package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/core/scrapper"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/pkg/pdfutil"
	"reflect"
	"testing"
)

func Test_webrepo_GetLinks(t *testing.T) {
	type args struct {
		page string
	}
	tests := []struct {
		name    string
		args    args
		want    []scrapper.Link
		wantErr bool
	}{
		{
			name:    "First Page",
			args:    args{page: KplcBase},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := webrepo{}
			got, err := w.GetLinks(tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLinks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webrepo_GetLinksFromLead(t *testing.T) {
	type args struct {
		lead string
	}
	tests := []struct {
		name    string
		args    args
		want    []scrapper.Link
		wantErr bool
	}{
		{
			args: args{
				lead: "https://kplc.co.ke/content/item/3765/interruptions---01.04.2021---special-power-interruption---friday-02.04.2021",
			},
			name:    "Kplc leads",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := webrepo{}
			got, err := w.GetLinksFromLead(tt.args.lead)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLinksFromLead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLinksFromLead() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webrepo_GetBlackoutResultFromPdf(t *testing.T) {
	w := webrepo{}
	pdfUrls, err := w.GetLinks(KplcBase)
	if err != nil {
		t.Fatal("fetch pdfs", err)
	}
	pdfUrl := pdfUrls[0].Url
	pdfPath, err := w.GenerateTmpPDF(pdfUrl)
	t.Log("pdfurl", pdfPath)

	if err != nil {
		t.Fatal("failed:", err)
	}

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *pdfutil.BlackoutResult
		wantErr bool
	}{
		{
			name: "All fns",
			args: args{
				path: pdfPath,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := w.GetBlackoutResultFromPdf(tt.args.path)
			bytes, err := json.Marshal(got)
			fmt.Println(string(bytes))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlackoutResultFromPdf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBlackoutResultFromPdf() got = %v, want %v", got, tt.want)
			}
		})
	}
}
