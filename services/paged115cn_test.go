package services

import (
	"os"
	"path/filepath"
	"testing"

	"pvfine/internal/pvf"
)

func TestPaged115CNReadOnlyBrowsing(t *testing.T) {
	dir := os.Getenv("PVF_115CN_DIR")
	if dir == "" {
		t.Skip("set PVF_115CN_DIR to the 115CN client directory")
	}

	archive, err := pvf.Open(filepath.Join(dir, "Script.pvf"))
	if err != nil {
		t.Fatalf("open encrypted 115CN archive: %v", err)
	}

	core := NewCore()
	defer core.closeArchive()
	if err := core.setArchive(archive); err != nil {
		t.Fatalf("build browsing index: %v", err)
	}

	archiveService := NewArchiveService(core)
	info := archiveService.Info()
	if info.Format != string(pvf.FormatPaged115CN) {
		t.Fatalf("format = %q, want %q", info.Format, pvf.FormatPaged115CN)
	}

	children, err := archiveService.ListChildren("")
	if err != nil {
		t.Fatalf("list root: %v", err)
	}
	if len(children) == 0 {
		t.Fatal("root listing is empty")
	}

	index, ok := archive.Find("etc/accountcargo.etc")
	if !ok {
		t.Fatal("etc/accountcargo.etc not found")
	}
	file, err := NewEditorService(core).GetFile(index)
	if err != nil {
		t.Fatalf("read sample file: %v", err)
	}
	if file.Path != "etc/accountcargo.etc" || file.Text == "" {
		t.Fatalf("unexpected editor result: path=%q textBytes=%d", file.Path, len(file.Text))
	}
}
