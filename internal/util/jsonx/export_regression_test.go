package jsonx

import (
	"testing"

	"example.com/inkkeep/internal/model"
)

// TestExportAgeRoundTrip is a regression test for the bug where ExportFile.Age
// carried the wrong JSON tag ("years" instead of "age"), causing the "age" key
// to be dropped on decode and Age to read as 0.
func TestExportAgeRoundTrip(t *testing.T) {
	orig := model.ExportFile{Name: "alpha", Age: 9}
	raw, err := EncodeExport(orig)
	if err != nil {
		t.Fatal(err)
	}
	// The on-wire key must be "age" so external producers/consumers agree.
	if want := `{"name":"alpha","age":9}`; string(raw) != want {
		t.Fatalf("encoded=%s want=%s", raw, want)
	}
	got, err := DecodeExport(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got != orig {
		t.Fatalf("round-trip mismatch: got=%+v want=%+v", got, orig)
	}
}

// TestExportAgeFromExternalJSON ensures decoding JSON authored by an external
// producer (using the "age" key) maps into the Age field rather than leaving 0.
func TestExportAgeFromExternalJSON(t *testing.T) {
	got, err := DecodeExport([]byte(`{"name":"beta","age":42}`))
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "beta" || got.Age != 42 {
		t.Fatalf("got=%+v", got)
	}
}
