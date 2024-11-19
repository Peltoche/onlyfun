package medias

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/spf13/afero"
)

var (
	errNotExist = errors.New("file doesn't exists")
)

type storageAfero struct {
	fs   afero.Fs
	uuid uuid.Service
}

func newStorageAfero(fs afero.Fs, dirPath string, tools tools.Tools) (*storageAfero, error) {
	root := path.Clean(path.Join(dirPath, "files"))
	rootFS := afero.NewBasePathFs(fs, root)

	err := fs.MkdirAll(root, 0o700)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return nil, fmt.Errorf("failed to create the files directory: %w", err)
	}

	err = setupFileDirectory(rootFS)
	if err != nil {
		return nil, fmt.Errorf("failed to setup the file storage directory: %w", err)
	}

	return &storageAfero{
		fs:   rootFS,
		uuid: tools.UUID(),
	}, nil
}

func (s *storageAfero) NewFileUploader() (uuid.UUID, io.WriteCloser, error) {
	fileID := s.uuid.New()
	filePath := pathFromFileID(fileID)

	file, err := s.fs.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return uuid.UUID(""), nil, fmt.Errorf("failed to create the file: %w", err)
	}

	return fileID, file, err
}

func (s *storageAfero) NewFileDownloader(fileID uuid.UUID) (io.ReadSeekCloser, error) {
	filePath := pathFromFileID(fileID)

	file, err := s.fs.OpenFile(filePath, os.O_RDONLY, 0o600)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%s: %w", filePath, ErrNotExist)
	}

	return file, nil
}

func (s *storageAfero) DeleteFile(fileID uuid.UUID) error {
	filePath := pathFromFileID(fileID)

	return s.fs.Remove(filePath)
}

func pathFromFileID(fileID uuid.UUID) string {
	idStr := string(fileID)
	return path.Join(idStr[:2], idStr)
}

func setupFileDirectory(rootFS afero.Fs) error {
	for i := 0; i < 256; i++ {
		dir := fmt.Sprintf("%02x", i)
		// XXX:MULTI-WRITE
		//
		// This function is idempotent so no worries. If it fails the server doesn't start
		// so we are sur that it will be run again until it's completely successful.
		err := rootFS.Mkdir(dir, 0o755)
		if errors.Is(err, os.ErrExist) {
			continue
		}

		if err != nil {
			return fmt.Errorf("failed to Mkdir %q: %w", dir, err)
		}
	}

	return nil
}
