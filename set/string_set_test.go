package set

import (
	"reflect"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		items []string
	}
	tests := []struct {
		name string
		args args
		want *StringSet
	}{
		{
			name: "not nil slice",
			args: args{
				items: []string{"a", "b", "c"},
			},
			want: &StringSet{
				m: map[string]bool{
					"a": true,
					"b": true,
					"c": true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.items...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSet_Add(t *testing.T) {
	type fields struct {
		s *StringSet
	}
	type args struct {
		items []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StringSet
	}{
		{
			name: "succeed",
			fields: fields{
				s: New("a"),
			},
			args: args{
				items: []string{"b", "c"},
			},
			want: &StringSet{
				m: map[string]bool{
					"a": true,
					"b": true,
					"c": true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.s.Add(tt.args.items...)
		})
		if !reflect.DeepEqual(tt.fields.s, tt.want) {
			t.Errorf("Add() = %v, want %v", tt.fields.s, tt.want)
		}
	}
}

func TestStringSet_Union(t *testing.T) {
	type fields struct {
		s *StringSet
	}
	type args struct {
		sets []*StringSet
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StringSet
	}{
		{
			name: "succeed",
			fields: fields{
				s: New("a"),
			},
			args: args{
				sets: []*StringSet{
					New("b", "c"),
					New("d"),
				},
			},
			want: &StringSet{
				m: map[string]bool{
					"a": true,
					"b": true,
					"c": true,
					"d": true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.s
			if got := s.Union(tt.args.sets...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSet.Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSet_Minus(t *testing.T) {
	type fields struct {
		s *StringSet
	}
	type args struct {
		sets []*StringSet
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StringSet
	}{
		{
			name: "case one",
			fields: fields{
				s: New("a"),
			},
			args: args{
				sets: []*StringSet{
					New("b", "c", "d"),
				},
			},
			want: &StringSet{
				m: map[string]bool{
					"a": true,
				},
			},
		},
		{
			name: "case two",
			fields: fields{
				s: New("a", "b", "c", "d"),
			},
			args: args{
				sets: []*StringSet{
					New("b", "c", "d"),
				},
			},
			want: &StringSet{
				m: map[string]bool{
					"a": true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.s
			if got := s.Minus(tt.args.sets...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSet.Minus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSet_Intersect(t *testing.T) {
	type fields struct {
		s *StringSet
	}
	type args struct {
		sets []*StringSet
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StringSet
	}{
		{
			name: "case one",
			fields: fields{
				s: New("a", "b"),
			},
			args: args{
				sets: []*StringSet{
					New("b", "c", "d"),
				},
			},
			want: &StringSet{
				m: map[string]bool{
					"b": true,
				},
			},
		},
		{
			name: "case two",
			fields: fields{
				s: New("a", "b", "c", "d"),
			},
			args: args{
				sets: []*StringSet{
					New("b", "c", "d"),
				},
			},
			want: &StringSet{
				m: map[string]bool{
					"b": true,
					"c": true,
					"d": true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.s
			if got := s.Intersect(tt.args.sets...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSet.Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSet_Remove(t *testing.T) {
	type fields struct {
		RWMutex sync.RWMutex
		m       map[string]bool
	}
	type args struct {
		items []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringSet{
				RWMutex: tt.fields.RWMutex,
				m:       tt.fields.m,
			}
			s.Remove(tt.args.items...)
		})
	}
}
