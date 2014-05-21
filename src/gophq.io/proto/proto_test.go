package proto

import (
	"reflect"
	"testing"
)

func TestProduceResponse(t *testing.T) {
	expected := ProduceResponse{"topic", 0, 1234}
	b, err := Encode(&expected)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	var actual ProduceResponse
	err = Decode(b, &actual)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %+v, got %v", expected, actual)
	}
}
