package utils

import "testing"

func TestGetDownloadFileName(t *testing.T) {
	type args struct {
		contentDisposition string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{contentDisposition: "attachment; filename=\"blindbox_v1_mix.safetensors\""},
			want: "blindbox_v1_mix.safetensors",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDownloadFileName(tt.args.contentDisposition)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDownloadFileName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDownloadFileName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
