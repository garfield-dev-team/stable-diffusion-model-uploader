package api

import (
	"log"
	"testing"
)

func TestFetchModelList(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FetchModelList()
			log.Printf("%+v", got)
		})
	}
}
