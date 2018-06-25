package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		name       string
		paths      []string
		configType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "load a file that not exist in current directory",
			args: args{
				name:       "wrong",
				paths:      []string{"."},
				configType: "yaml",
			},
			wantErr: true,
		},
		{
			name: "load test.yaml in current directory",
			args: args{
				name:       "test",
				paths:      []string{"."},
				configType: "yaml",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadConfig(tt.args.name, tt.args.paths, tt.args.configType); (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get exist string config",
			args: args{
				key: "string",
			},
			want: "echo",
		},
		{
			name: "Get none exist string config",
			args: args{
				key: "no-string",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadConfig("test", []string{"."}, "yaml")
			if got := GetString(tt.args.key); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSlice(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Get exist string slice config",
			args: args{
				key: "slice",
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "Get none exist string slice config",
			args: args{
				key: "no-slice",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadConfig("test", []string{"."}, "yaml")
			if got := GetSlice(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Get exist integer config",
			args: args{
				key: "int",
			},
			want: 1,
		},
		{
			name: "Get none exist string slice config",
			args: args{
				key: "no-slice",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadConfig("test", []string{"."}, "yaml")
			if got := GetInt(tt.args.key); got != tt.want {
				t.Errorf("GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetDefaultConfig(t *testing.T) {
	type args struct {
		values map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Set default config",
			args: args{
				values: map[string]interface{}{
					"default-string": "a",
					"default-int":    1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadConfig("test", []string{"."}, "yaml")
			SetDefaultConfig(tt.args.values)
			if tt.args.values["default-string"] != GetString("default-string") {
				t.Errorf("Get default string error")
			}
			if tt.args.values["default-int"] != GetInt("default-int") {
				t.Errorf("Get default int error")
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		key    string
		rawVal interface{}
	}
	var wc []Config
	tests := []struct {
		name    string
		args    args
		want    *[]Config
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				key:    "configs",
				rawVal: &wc,
			},
			want: &[]Config{
				{
					ID:   1000000,
					Name: "name",
					List: []string{"a", "b"},
				},
				{
					ID:   1000001,
					Name: "name1",
					List: []string{"c", "d"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		LoadConfig("test", []string{"."}, "yaml")
		t.Run(tt.name, func(t *testing.T) {
			if err := Unmarshal(tt.args.key, tt.args.rawVal); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(&wc, tt.want) {
				t.Errorf("Unmarshal() error, got %v, want %v", wc, tt.want)
			}
		})
	}
}

type Config struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	List []string `json:"list"`
}
