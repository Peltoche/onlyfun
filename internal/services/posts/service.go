package posts

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/clock"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
)

const (
	maxImgSizeBytes  = 5 * 1024 * 1024 // 5MiB
	maxPostBatchSize = 100
)

var ErrToMuchPostsAsked = errors.New("too much posts asks")

type storage interface {
	Save(ctx context.Context, post *Post) error
	GetLatestPostWithStatus(ctx context.Context, status Status) (*Post, error)
	GetOldestPostWithStatus(ctx context.Context, status Status) (*Post, error)
	GetListedPosts(ctx context.Context, start uint, limit uint) ([]Post, error)
	GetByID(ctx context.Context, postID uint) (*Post, error)
	CountPostsWithStatus(ctx context.Context, status Status) (int, error)
	CountUserPostsByStatus(ctx context.Context, userID uuid.UUID, status Status) (int, error)
	Update(ctx context.Context, post *Post) error
}

type service struct {
	storage      storage
	mediasSvc    medias.Service
	permsSvc     perms.Service
	clock        clock.Clock
	uuid         uuid.Service
	newPostChans []chan Post
	l            *sync.Mutex
}

func newService(tools tools.Tools, posts storage, mediasSvc medias.Service, permsSvc perms.Service) *service {
	svc := &service{
		storage:      posts,
		mediasSvc:    mediasSvc,
		permsSvc:     permsSvc,
		clock:        tools.Clock(),
		uuid:         tools.UUID(),
		newPostChans: make([]chan Post, 0),
		l:            new(sync.Mutex),
	}

	return svc
}

func (s *service) SuscribeToNewPost() <-chan Post {
	s.l.Lock()
	defer s.l.Unlock()

	ch := make(chan Post)
	s.newPostChans = append(s.newPostChans, ch)

	return ch
}

func (s *service) GetByID(ctx context.Context, postID uint) (*Post, error) {
	res, err := s.storage.GetByID(ctx, postID)
	if errors.Is(err, errNotFound) {
		return nil, errs.NotFound(fmt.Errorf("no post available"))
	}

	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to GetLatestPost: %w", err))
	}

	return res, nil
}

func (s *service) ValidatePost(ctx context.Context, cmd *ValidatePostcmd) error {
	if !s.permsSvc.IsAuthorized(cmd.User, perms.Moderation) {
		return errs.Unauthorized(fmt.Errorf("user %q doesn't have the authorization %q", cmd.User.ID(), perms.Moderation))
	}

	cmd.Post.status = Listed

	err := s.storage.Update(ctx, cmd.Post)
	if err != nil {
		return fmt.Errorf("failed to Update post %d: %w", cmd.Post.id, err)
	}

	return nil
}

func (s *service) Create(ctx context.Context, cmd *CreateCmd) (*Post, error) {
	err := cmd.Validate()
	if err != nil {
		return nil, errs.Validation(err)
	}

	meta, err := s.mediasSvc.Upload(ctx, medias.Post, cmd.Media)
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to upload the media: %w", err))
	}

	// TODO: Need a lot of validation:
	// - Size
	// - Mimetype
	// - Checksum already exists
	// - Ratio

	post := Post{
		status:    Uploaded,
		title:     cmd.Title,
		fileID:    meta.ID(),
		createdBy: cmd.CreatedBy.ID(),
		createdAt: s.clock.Now(),
	}

	err = s.storage.Save(ctx, &post)
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to save the post: %w", err))
	}

	go func() {
		for _, ch := range s.newPostChans {
			ch <- post
		}
	}()

	return &post, nil
}

func (s *service) GetNextPostToModerate(ctx context.Context) (*Post, error) {
	res, err := s.storage.GetOldestPostWithStatus(ctx, Uploaded)
	if errors.Is(err, errNotFound) {
		return nil, errs.NotFound(fmt.Errorf("no post available"))
	}

	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to GetLatestPost: %w", err))
	}

	return res, nil
}

func (s *service) GetUserStats(ctx context.Context, user *users.User) (map[Status]int, error) {
	var err error

	status := []Status{Uploaded, Listed, Moderated}

	stats := make(map[Status]int, len(status))

	for _, status := range status {
		stats[status], err = s.storage.CountUserPostsByStatus(ctx, user.ID(), status)
		if err != nil {
			return nil, errs.Internal(fmt.Errorf("failed to CountUserPostsBystatus with status %q for user %q: %w", status, user.ID(), err))
		}
	}

	return stats, nil
}

func (s *service) GetLatestPost(ctx context.Context) (*Post, error) {
	res, err := s.storage.GetLatestPostWithStatus(ctx, Listed)
	if errors.Is(err, errNotFound) {
		return nil, errs.NotFound(fmt.Errorf("no post available"))
	}

	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to GetLatestPost: %w", err))
	}

	return res, nil
}

func (s *service) CountPostsWaitingModeration(ctx context.Context) (int, error) {
	res, err := s.storage.CountPostsWithStatus(ctx, Uploaded)
	if err != nil {
		return 0, errs.Internal(fmt.Errorf("failed to CountPostsWithStatus: %w", err))
	}

	return res, nil
}

func (s *service) GetPosts(ctx context.Context, start uint, nbPosts uint) ([]Post, error) {
	if nbPosts > maxPostBatchSize {
		return nil, errs.Validation(ErrToMuchPostsAsked)
	}

	res, err := s.storage.GetListedPosts(ctx, start, nbPosts)
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to GetAll: %w", err))
	}

	return res, nil
}
