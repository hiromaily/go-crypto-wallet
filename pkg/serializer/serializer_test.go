package serializer

import (
	"testing"
)

type testStruct struct {
	Name  string
	Value int
	Data  []byte
}

func TestGobSerializer_EncodeDecode(t *testing.T) {
	serializer := NewGobSerializer()

	original := &testStruct{
		Name:  "test",
		Value: 42,
		Data:  []byte{1, 2, 3, 4, 5},
	}

	// Encode
	encoded, err := serializer.EncodeToString(original)
	if err != nil {
		t.Fatalf("EncodeToString failed: %v", err)
	}
	if encoded == "" {
		t.Fatal("Encoded string is empty")
	}

	// Decode
	var decoded testStruct
	err = serializer.DecodeFromString(encoded, &decoded)
	if err != nil {
		t.Fatalf("DecodeFromString failed: %v", err)
	}

	// Verify
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, original.Name)
	}
	if decoded.Value != original.Value {
		t.Errorf("Value mismatch: got %d, want %d", decoded.Value, original.Value)
	}
	if len(decoded.Data) != len(original.Data) {
		t.Errorf("Data length mismatch: got %d, want %d", len(decoded.Data), len(original.Data))
	}
	for i := range decoded.Data {
		if decoded.Data[i] != original.Data[i] {
			t.Errorf("Data[%d] mismatch: got %d, want %d", i, decoded.Data[i], original.Data[i])
		}
	}
}

func TestJSONSerializer_EncodeDecode(t *testing.T) {
	serializer := NewJSONSerializer()

	original := &testStruct{
		Name:  "test",
		Value: 42,
		Data:  []byte{1, 2, 3, 4, 5},
	}

	// Encode
	encoded, err := serializer.EncodeToString(original)
	if err != nil {
		t.Fatalf("EncodeToString failed: %v", err)
	}
	if encoded == "" {
		t.Fatal("Encoded string is empty")
	}

	// Decode
	var decoded testStruct
	err = serializer.DecodeFromString(encoded, &decoded)
	if err != nil {
		t.Fatalf("DecodeFromString failed: %v", err)
	}

	// Verify
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, original.Name)
	}
	if decoded.Value != original.Value {
		t.Errorf("Value mismatch: got %d, want %d", decoded.Value, original.Value)
	}
	if len(decoded.Data) != len(original.Data) {
		t.Errorf("Data length mismatch: got %d, want %d", len(decoded.Data), len(original.Data))
	}
	for i := range decoded.Data {
		if decoded.Data[i] != original.Data[i] {
			t.Errorf("Data[%d] mismatch: got %d, want %d", i, decoded.Data[i], original.Data[i])
		}
	}
}

func TestDefaultSerializer_BackwardCompatibility(t *testing.T) {
	original := &testStruct{
		Name:  "test",
		Value: 42,
		Data:  []byte{1, 2, 3, 4, 5},
	}

	// Test default serializer (should be gob)
	encoded, err := EncodeToString(original)
	if err != nil {
		t.Fatalf("EncodeToString failed: %v", err)
	}

	var decoded testStruct
	err = DecodeFromString(encoded, &decoded)
	if err != nil {
		t.Fatalf("DecodeFromString failed: %v", err)
	}

	// Verify
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, original.Name)
	}
	if decoded.Value != original.Value {
		t.Errorf("Value mismatch: got %d, want %d", decoded.Value, original.Value)
	}
}

func TestSetDefaultSerializer(t *testing.T) {
	// Save original
	original := GetDefaultSerializer()

	// Set to JSON
	jsonSerializer := NewJSONSerializer()
	SetDefaultSerializer(jsonSerializer)

	// Verify it's changed
	current := GetDefaultSerializer()
	if current != jsonSerializer {
		t.Error("SetDefaultSerializer did not change the default serializer")
	}

	// Restore original
	SetDefaultSerializer(original)
}

func TestGobSerializer_CrossCompatibility(t *testing.T) {
	// Test that gob-encoded data can be decoded by gob decoder
	gobSerializer := NewGobSerializer()

	original := &testStruct{
		Name:  "test",
		Value: 42,
		Data:  []byte{1, 2, 3, 4, 5},
	}

	encoded, err := gobSerializer.EncodeToString(original)
	if err != nil {
		t.Fatalf("EncodeToString failed: %v", err)
	}

	// Decode using default serializer (should be gob)
	var decoded testStruct
	err = DecodeFromString(encoded, &decoded)
	if err != nil {
		t.Fatalf("DecodeFromString failed: %v", err)
	}

	// Verify
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, original.Name)
	}
}
