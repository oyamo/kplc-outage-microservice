package repositories

import (
	"github.com/oyamo/kplc-outage-microservice/scrapping-service/internal/local/scrapper"
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
