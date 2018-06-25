package http

import (
	"testing"

	"github.com/v-zhidu/orb/logging"
)

func init() {
	logging.SetLevel("debug")
}

func TestGet(t *testing.T) {
	type args struct {
		url    string
		params map[string]string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse bool
		wantErr      bool
	}{
		{
			name: "Valid endpoint",
			args: args{
				url: "https://postman-echo.com/get",
				params: map[string]string{
					"test": "3",
				},
			},
			wantResponse: true,
			wantErr:      false,
		},
		{
			name: "Invalid endpoint",
			args: args{
				url: "https://invalid.com/get",
				params: map[string]string{
					"test": "3",
				},
			},
			wantResponse: false,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.url, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResponse == true && got == nil {
				t.Errorf("Get() = %v", got)
				return
			}
		})
	}
}

func TestPostJSON(t *testing.T) {
	type args struct {
		url     string
		body    []byte
		headers map[string]string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse bool
		wantErr      bool
	}{
		{
			name: "Valid endpoint",
			args: args{
				url:     "https://postman-echo.com/post",
				body:    []byte("{\"foo1\":\"bar1\",\"foo2\":\"bar2\"}"),
				headers: nil,
			},
			wantResponse: true,
			wantErr:      false,
		},
		{
			name: "Invalid endpoint",
			args: args{
				url:     "https://invalid.com/post",
				body:    []byte("{\"foo1\":\"bar1\",\"foo2\":\"bar2\"}"),
				headers: nil,
			},
			wantResponse: false,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostJSON(tt.args.url, tt.args.body, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResponse == true && got == nil {
				t.Errorf("Get() = %v", got)
				return
			}
		})
	}
}

func TestPostMap(t *testing.T) {
	type args struct {
		url     string
		body    map[string]interface{}
		headers map[string]string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse bool
		wantErr      bool
	}{
		{
			name: "Valid endpoint",
			args: args{
				url: "https://postman-echo.com/post",
				body: map[string]interface{}{
					"foo1": "bar1",
					"foo2": "bar2",
				},
				headers: nil,
			},
			wantResponse: true,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostMap(tt.args.url, tt.args.body, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResponse == true && got == nil {
				t.Errorf("Get() = %v", got)
				return
			}
		})
	}
}
