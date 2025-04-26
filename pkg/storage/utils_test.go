package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getFilesFromFolder(t *testing.T) {
	filesFolder := "../../files"

	t.Run("folder with files", func(t *testing.T) {
		// This assumes the folder already has some files inside
		got, err := getFilesFromFolder(filesFolder)
		require.NoError(t, err)

		if len(got) < 5 {
			t.Errorf("expected at least 5 files, got %d", len(got))
		}
	})

	t.Run("folder does not exist", func(t *testing.T) {
		_, err := getFilesFromFolder("/not/exist")
		require.Error(t, err)
	})

	t.Run("empty folder", func(t *testing.T) {
		// Create an empty temp folder
		dir := t.TempDir()

		got, err := getFilesFromFolder(dir)
		require.NoError(t, err)
		require.Len(t, got, 0)
	})
}
