package micro

import (
	"reflect"
	"testing"
)

type testClient struct{}

func (c *testClient) Name() string {
	return "client"
}

func TestClientSet_NewClientSet(t *testing.T) {
	tests := []struct {
		name string
		args []Client
		want int
	}{
		{"With no clients", []Client{}, 0},
		{"With a client", []Client{&testClient{}}, 1},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			cs := NewClientSet(test.args...)
			if got := cs.Length(); got != test.want {
				t.Errorf("ClientSet.Length() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestClientSet_Get(t *testing.T) {
	c := &testClient{}
	cs := NewClientSet(c)

	tests := []struct {
		name string
		c    *ClientSet
		args string
		want Client
	}{
		{"A valid key", cs, c.Name(), c},
		{"An invalid key", cs, "invalid", nil},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			if got := test.c.Get(test.args); !reflect.DeepEqual(got, test.want) {
				t.Errorf("ClientSet.Get() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestClientSet_Keys(t *testing.T) {
	c := &testClient{}

	tests := []struct {
		name string
		c    *ClientSet
		want []string
	}{
		{"An empty ClientSet", NewClientSet(), []string{}},
		{"A non-empty ClientSet", NewClientSet(c), []string{c.Name()}},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			if got := test.c.Keys(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("ClientSet.Keys() = %v, want %v", got, test.want)
			}
		})
	}
}
