package files

import (
	"reflect"
	"testing"
)

func Test_mergeMap(t *testing.T) {
	type args struct {
		a map[string]string
		b map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "no overlaps",
			args: args{
				a: map[string]string{"name": "flanksource"},
				b: map[string]string{"foo": "bar"},
			},
			want: map[string]string{
				"name": "flanksource",
				"foo":  "bar",
			},
		},
		{
			name: "overlaps",
			args: args{
				a: map[string]string{"name": "flanksource", "foo": "baz"},
				b: map[string]string{"foo": "bar"},
			},
			want: map[string]string{
				"name": "flanksource",
				"foo":  "bar",
			},
		},
		{
			name: "overlaps II",
			args: args{
				a: map[string]string{"name": "github", "foo": "baz"},
				b: map[string]string{"name": "flanksource", "foo": "bar"},
			},
			want: map[string]string{
				"name": "flanksource",
				"foo":  "bar",
			},
		},
		{
			name: "ditto",
			args: args{
				a: map[string]string{"name": "flanksource", "foo": "bar"},
				b: map[string]string{"name": "flanksource", "foo": "bar"},
			},
			want: map[string]string{
				"name": "flanksource",
				"foo":  "bar",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeMap(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
