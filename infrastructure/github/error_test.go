package github

import (
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	type fields struct {
		StatusCode int
		Message    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Get 401 error",
			fields: fields{
				StatusCode: http.StatusUnauthorized,
				Message:    "unauthorized",
			},
			want: "GitHub API error: status code 401, message: unauthorized",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &APIError{
				StatusCode: tt.fields.StatusCode,
				Message:    tt.fields.Message,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
