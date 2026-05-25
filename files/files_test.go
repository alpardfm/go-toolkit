package files

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetExtension(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "simple txt extension",
			filename: "document.txt",
			want:     "txt",
		},
		{
			name:     "csv extension",
			filename: "data.csv",
			want:     "csv",
		},
		{
			name:     "multiple dots returns last extension",
			filename: "archive.tar.gz",
			want:     "gz",
		},
		{
			name:     "no extension",
			filename: "README",
			want:     "",
		},
		{
			name:     "hidden file with extension",
			filename: ".gitignore.bak",
			want:     "bak",
		},
		{
			name:     "hidden file without extension",
			filename: ".gitignore",
			want:     "gitignore",
		},
		{
			name:     "empty string",
			filename: "",
			want:     "",
		},
		{
			name:     "path with extension",
			filename: "/path/to/file.json",
			want:     "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetExtension(tt.filename)
			if got != tt.want {
				t.Errorf("GetExtension(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}

func TestIsExist(t *testing.T) {
	t.Run("existing file returns true", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, "testfile.txt")
		if err := os.WriteFile(filePath, []byte("hello"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		if !IsExist(filePath) {
			t.Errorf("IsExist(%q) = false, want true", filePath)
		}
	})

	t.Run("non-existing file returns false", func(t *testing.T) {
		if IsExist("/nonexistent/path/file.txt") {
			t.Error("IsExist() = true for non-existing file, want false")
		}
	})

	t.Run("directory returns false", func(t *testing.T) {
		dir := t.TempDir()
		if IsExist(dir) {
			t.Errorf("IsExist(%q) = true for directory, want false", dir)
		}
	})

	t.Run("empty path returns false", func(t *testing.T) {
		if IsExist("") {
			t.Error("IsExist(\"\") = true, want false")
		}
	})
}

func TestGetCurrentFileLocation(t *testing.T) {
	t.Run("returns non-empty path", func(t *testing.T) {
		loc := GetCurrentFileLocation()
		if loc == "" {
			t.Error("GetCurrentFileLocation() returned empty string")
		}
	})

	t.Run("path contains files_test.go", func(t *testing.T) {
		loc := GetCurrentFileLocation()
		if !strings.HasSuffix(loc, "files_test.go") {
			t.Errorf("GetCurrentFileLocation() = %q, want suffix 'files_test.go'", loc)
		}
	})
}

func TestGetCurrentMethodName(t *testing.T) {
	t.Run("returns non-empty name", func(t *testing.T) {
		name := GetCurrentMethodName()
		if name == "" {
			t.Error("GetCurrentMethodName() returned empty string")
		}
	})

	t.Run("contains function suffix", func(t *testing.T) {
		name := GetCurrentMethodName()
		if !strings.HasSuffix(name, "()") {
			t.Errorf("GetCurrentMethodName() = %q, want suffix '()'", name)
		}
	})
}
