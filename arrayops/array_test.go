package arrayops

import (
	"reflect"
	"testing"
)

func TestContainsInt(t *testing.T) {
	type args struct {
		s []int
		e int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "return true if element matches",
			args: args{
				s: []int{1, 2, 3, 4},
				e: 1,
			},
			want: true,
		},
		{
			name: "return false if element not matches",
			args: args{
				s: []int{1, 2, 3, 4},
				e: 5,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsInt(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("ContainsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsString(t *testing.T) {
	type args struct {
		s []string
		e string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "return true if element matches",
			args: args{
				s: []string{"a", "b", "c", "d"},
				e: "a",
			},
			want: true,
		},
		{
			name: "return false if element not matches",
			args: args{
				s: []string{"a", "b", "c", "d"},
				e: "e",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsString(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("ContainsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueInt(t *testing.T) {
	type args struct {
		a []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "return unique array",
			args: args{
				a: []int{1, 2, 3, 4, 4, 3},
			},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "return same array if array is already unique",
			args: args{
				a: []int{1, 2, 3, 4},
			},
			want: []int{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UniqueIntArray(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubtractString(t *testing.T) {
	type args struct {
		x []string
		y []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "remove elements from array case 1",
			args: args{
				x: []string{"a", "b", "c", "d"},
				y: []string{"c", "d"},
			},
			want: []string{"a", "b"},
		},
		{
			name: "remove elements from array case 2",
			args: args{
				x: []string{"a", "b", "c", "d"},
				y: []string{"c", "d", "e", "f"},
			},
			want: []string{"a", "b"},
		},
		{
			name: "remove multiple same elements from array",
			args: args{
				x: []string{"a", "b", "c", "c", "d"},
				y: []string{"c", "d"},
			},
			want: []string{"a", "b"},
		},
		{
			name: "do not remove element if removes elements is blank",
			args: args{
				x: []string{"a", "b", "c", "d"},
				y: []string{},
			},
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "do not remove element if  elements to removes not found",
			args: args{
				x: []string{"a", "b", "c", "d"},
				y: []string{"f"},
			},
			want: []string{"a", "b", "c", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubtractString(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubtractString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueStringArray(t *testing.T) {
	type args struct {
		a []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "return unique array",
			args: args{
				a: []string{"a", "b", "c", "c", "d"},
			},
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "return same array if array is already unique",
			args: args{
				a: []string{"a", "b", "c", "d"},
			},
			want: []string{"a", "b", "c", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UniqueStringArray(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueStringArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want []string
	}{
		{
			name: "successfully found the intersection",
			a:    []string{"abc", "xyz"},
			b:    []string{"xyz", "pqr"},
			want: []string{"xyz"},
		},
		{
			name: "no inytersection found",
			a:    []string{"abc", "xyz"},
			b:    []string{"mno", "pqr"},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Intersection(tt.a, tt.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrSliceToLowerCase(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "convert all string to lower case",
			args: []string{"ABC", "XYZ"},
			want: []string{"abc", "xyz"},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrSliceToLowerCase(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrSliceToLowerCase() = %v, want %v", got, tt.want)
			}
		})

	}
}
