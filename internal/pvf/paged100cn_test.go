package pvf

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

const test100CNDirEnv = "PVF_100CN_DIR"

func open100CNFixture(t *testing.T) (*Archive, string) {
	t.Helper()
	dir := os.Getenv(test100CNDirEnv)
	if dir == "" {
		t.Skipf("%s not set; skipping 100CN integration test", test100CNDirEnv)
	}
	path := filepath.Join(dir, "Script.pvf")
	a, err := Open(path)
	if err != nil {
		t.Fatalf("open 100CN archive: %v", err)
	}
	return a, path
}

func TestPaged100CNOpen(t *testing.T) {
	a, _ := open100CNFixture(t)
	if got := a.Format(); got != FormatPaged100CN {
		t.Fatalf("format = %q, want %q", got, FormatPaged100CN)
	}
	if got := a.ClientVersion(); got != "100CN" {
		t.Fatalf("client version = %q, want 100CN", got)
	}
	if !a.IsPaged110() {
		t.Fatal("100CN archive did not retain the paged container")
	}
	if got, want := a.FileCount(), int32(2407803); got != want {
		t.Fatalf("file count = %d, want %d", got, want)
	}
	if got, want := a.Header().GroupCount, int32(28579); got != want {
		t.Fatalf("group count = %d, want %d", got, want)
	}
	if len(a.pageKeys) != 30*paged110PageKeySize {
		t.Fatalf("page key bytes = %d, want %d", len(a.pageKeys), 30*paged110PageKeySize)
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

func TestPaged100CNNoOpSaveIsByteExact(t *testing.T) {
	a, source := open100CNFixture(t)
	var saved bytes.Buffer
	if err := a.SaveTo(&saved); err != nil {
		t.Fatalf("SaveTo: %v", err)
	}
	original, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(saved.Bytes(), original) {
		t.Fatal("no-op 100CN save differs from the original encrypted archive")
	}
}
