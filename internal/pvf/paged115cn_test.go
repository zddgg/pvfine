package pvf

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

const test115CNDirEnv = "PVF_115CN_DIR"

func open115CNFixture(t *testing.T) (*Archive, string) {
	t.Helper()
	dir := os.Getenv(test115CNDirEnv)
	if dir == "" {
		t.Skipf("%s not set; skipping 115CN integration test", test115CNDirEnv)
	}
	path := filepath.Join(dir, "Script.pvf")
	a, err := Open(path)
	if err != nil {
		t.Fatalf("open 115CN archive: %v", err)
	}
	return a, path
}

func TestPaged115CNOpen(t *testing.T) {
	a, _ := open115CNFixture(t)
	if got := a.Format(); got != FormatPaged115CN {
		t.Fatalf("format = %q, want %q", got, FormatPaged115CN)
	}
	if got := a.ClientVersion(); got != "115CN" {
		t.Fatalf("client version = %q, want 115CN", got)
	}
	if !a.IsPaged110() {
		t.Fatal("115CN archive did not retain the paged container")
	}
	if got, want := a.FileCount(), int32(6259057); got != want {
		t.Fatalf("file count = %d, want %d", got, want)
	}
	if got, want := a.Header().GroupCount, int32(96237); got != want {
		t.Fatalf("group count = %d, want %d", got, want)
	}
	if len(a.pageKeys) != 81*paged110PageKeySize {
		t.Fatalf("page key bytes = %d, want %d", len(a.pageKeys), 81*paged110PageKeySize)
	}
	index, ok := a.Find("etc/accountcargo.etc")
	if !ok {
		t.Fatal("etc/accountcargo.etc not found")
	}
	text, err := a.Text(index)
	if err != nil {
		t.Fatalf("read sample script: %v", err)
	}
	if text == "" {
		t.Fatal("sample script decoded to empty text")
	}
}

func TestPaged115CNNoOpSaveIsByteExact(t *testing.T) {
	a, source := open115CNFixture(t)
	var saved bytes.Buffer
	if err := a.SaveTo(&saved); err != nil {
		t.Fatalf("SaveTo: %v", err)
	}
	original, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(saved.Bytes(), original) {
		t.Fatal("no-op 115CN save differs from the original encrypted archive")
	}
}
