package medias

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/clock"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/sync/errgroup"
)

var (
	ErrInvalidPath   = errors.New("invalid path")
	ErrInodeNotAFile = errors.New("inode doesn't point to a file")
	ErrNotExist      = errors.New("file not exists")
)

type mediaStorage interface {
	Save(ctx context.Context, meta *FileMeta) error
	GetByID(ctx context.Context, id uuid.UUID) (*FileMeta, error)
	GetByChecksum(ctx context.Context, checksum string) (*FileMeta, error)
	Delete(ctx context.Context, fileID uuid.UUID) error
}

type fileStorage interface {
	NewFileUploader() (uuid.UUID, io.WriteCloser, error)
	NewFileDownloader(fileID uuid.UUID) (io.ReadSeekCloser, error)
	DeleteFile(fileID uuid.UUID) error
}

type service struct {
	uuid         uuid.Service
	clock        clock.Clock
	fileStorage  fileStorage
	mediaStorage mediaStorage
}

func newService(fileStorage fileStorage, mediaStorage mediaStorage, tools tools.Tools) *service {
	return &service{
		fileStorage:  fileStorage,
		mediaStorage: mediaStorage,
		uuid:         tools.UUID(),
		clock:        tools.Clock(),
	}
}

func (s *service) Upload(ctx context.Context, mediaType MediaType, r io.Reader) (*FileMeta, error) {
	g, ctx := errgroup.WithContext(ctx)

	fileID, file, err := s.fileStorage.NewFileUploader()
	if err != nil {
		return nil, fmt.Errorf("failed to create the FileUploader: %w", err)
	}
	defer file.Close()

	// Start the hasher job
	hashReader, hashWriter := io.Pipe()
	hasher := sha256.New()
	g.Go(func() error {
		_, err := io.Copy(hasher, hashReader)
		if err != nil {
			hashReader.CloseWithError(fmt.Errorf("failed to calculate the file hash: %w", err))
		}

		return nil
	})

	// Start the mime type detection
	var mimeStr string
	mimeReader, mimeWriter := io.Pipe()
	g.Go(func() error {
		mime, err := mimetype.DetectReader(mimeReader)
		if err != nil {
			mimeReader.CloseWithError(fmt.Errorf("failed to detect the mime type: %w", err))
			return nil
		}

		mimeStr = mime.String()

		io.Copy(io.Discard, mimeReader)

		return nil
	})

	multiWrite := io.MultiWriter(mimeWriter, hashWriter, file)

	written, err := io.Copy(multiWrite, r)
	if err != nil {
		_ = s.Delete(context.WithoutCancel(ctx), fileID)
		return nil, fmt.Errorf("upload error: %w", err)
	}

	_ = mimeWriter.Close()
	_ = hashWriter.Close()

	err = file.Close()
	if err != nil {
		_ = s.Delete(context.WithoutCancel(ctx), fileID)
		return nil, fmt.Errorf("failed to end the file encryption: %w", err)
	}

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	ctx = context.WithoutCancel(ctx)

	checksum := base64.RawStdEncoding.Strict().EncodeToString(hasher.Sum(nil))

	existingFile, err := s.mediaStorage.GetByChecksum(ctx, checksum)
	if err != nil && !errors.Is(err, errNotFound) {
		return nil, errs.Internal(fmt.Errorf("failed to GetByChecksum: %w", err))
	}

	if existingFile != nil {
		_ = s.Delete(context.WithoutCancel(ctx), fileID)
		return existingFile, nil
	}

	fileMeta := FileMeta{
		id:         fileID,
		size:       uint64(written),
		mimetype:   mimeStr,
		mediaType:  mediaType,
		checksum:   checksum,
		uploadedAt: s.clock.Now(),
	}

	// XXX:MULTI-WRITE
	err = s.mediaStorage.Save(ctx, &fileMeta)
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to save the file meta: %w", err))
	}

	return &fileMeta, nil
}

func (s *service) GetMetadataByChecksum(ctx context.Context, checksum string) (*FileMeta, error) {
	res, err := s.mediaStorage.GetByChecksum(ctx, checksum)
	if errors.Is(err, errNotFound) {
		return nil, ErrNotExist
	}

	return res, err
}

func (s *service) GetMetadata(ctx context.Context, fileID uuid.UUID) (*FileMeta, error) {
	res, err := s.mediaStorage.GetByID(ctx, fileID)
	if errors.Is(err, errNotFound) {
		return nil, ErrNotExist
	}

	return res, err
}

func (s *service) Download(ctx context.Context, fileID uuid.UUID) (io.ReadSeekCloser, error) {
	file, err := s.fileStorage.NewFileDownloader(fileID)
	if errors.Is(err, errNotExist) {
		return nil, errs.NotFound(fmt.Errorf("%s not found", fileID))
	}

	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to open the file: %w", err))
	}

	return file, nil
}

func (s *service) Delete(ctx context.Context, fileID uuid.UUID) error {
	err := s.mediaStorage.Delete(ctx, fileID)
	if err != nil {
		return errs.Internal(fmt.Errorf("failed to delete the file metadatas: %w", err))
	}

	return s.fileStorage.DeleteFile(fileID)
}
