package platio_test

import (
	"encoding/json"
	. "go-platio-example/platio"
	"reflect"
	"testing"
)

func TestRecord(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		record := Record{
			Id: "r11111111111111111111111111",
			Values: Values{
				Name: &StringValue{"abc"},
				Age:  &NumberValue{30},
			},
		}
		expected := `{"id":"r11111111111111111111111111","values":{"cd33ed98":{"type":"String","value":"abc"},"ce0f2361":{"type":"Number","value":30}}}`
		b, err := json.Marshal(record)
		if err != nil || string(b) != expected {
			t.Error("Marshaling Record failed", err, string(b))
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		s := []byte(`{"id":"r11111111111111111111111111","values":{"cd33ed98":{"type":"String","value":"abc"},"ce0f2361":{"type":"Number","value":30}}}`)
		expected := Record{
			Id: "r11111111111111111111111111",
			Values: Values{
				Name: &StringValue{"abc"},
				Age:  &NumberValue{30},
			},
		}
		var record Record
		if err := json.Unmarshal(s, &record); err != nil || !reflect.DeepEqual(record, expected) {
			t.Error("Unmarshaling Record failed", err, record)
		}
	})
}

func TestValues(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		t.Run("all", func(t *testing.T) {
			values := Values{
				Name: &StringValue{"abc"},
				Age:  &NumberValue{30},
			}
			expected := `{"cd33ed98":{"type":"String","value":"abc"},"ce0f2361":{"type":"Number","value":30}}`
			b, err := json.Marshal(values)
			if err != nil || string(b) != expected {
				t.Error("Marshaling Values failed", err, string(b))
			}
		})

		t.Run("omitempty", func(t *testing.T) {
			values := Values{
				Age: &NumberValue{30},
			}
			expected := `{"ce0f2361":{"type":"Number","value":30}}`
			b, err := json.Marshal(values)
			if err != nil || string(b) != expected {
				t.Error("Marshaling Values failed", err, string(b))
			}
		})
	})

	t.Run("unmarshal", func(t *testing.T) {
		t.Run("all", func(t *testing.T) {
			s := []byte(`{"cd33ed98":{"type":"String","value":"abc"},"ce0f2361":{"type":"Number","value":30}}`)
			expected := Values{
				Name: &StringValue{"abc"},
				Age:  &NumberValue{30},
			}
			var values Values
			if err := json.Unmarshal(s, &values); err != nil || !reflect.DeepEqual(values, expected) {
				t.Error("Marshaling Values failed", err, values)
			}
		})

		t.Run("missing", func(t *testing.T) {
			s := []byte(`{"ce0f2361":{"type":"Number","value":30}}`)
			expected := Values{
				Age: &NumberValue{30},
			}
			var values Values
			if err := json.Unmarshal(s, &values); err != nil || !reflect.DeepEqual(values, expected) {
				t.Error("Marshaling Values failed", err, values)
			}
		})
	})
}

func TestStringValue(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		value := &StringValue{"abc"}
		expected := `{"type":"String","value":"abc"}`
		b, err := json.Marshal(value)
		if err != nil || string(b) != expected {
			t.Error("Marshaling StringValue failed", err, string(b))
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		s := []byte(`{"type":"String","value":"abc"}`)
		expected := StringValue{"abc"}
		var value StringValue
		if err := json.Unmarshal(s, &value); err != nil || value != expected {
			t.Error("Unmarshaling StringValue failed", value)
		}
	})
}

func TestNumberValue(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		value := &NumberValue{10}
		expected := `{"type":"Number","value":10}`
		b, err := json.Marshal(value)
		if err != nil || string(b) != expected {
			t.Error("Marshaling NumberValue failed", err, string(b))
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		s := []byte(`{"type":"Number","value":10}`)
		expected := NumberValue{10}
		var value NumberValue
		if err := json.Unmarshal(s, &value); err != nil || value != expected {
			t.Error("Unmarshaling NumberValue failed", err, value)
		}
	})
}
