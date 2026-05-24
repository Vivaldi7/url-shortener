package note

import (
	"context"
	//	"log"

	sq "github.com/Masterminds/squirrel"
	//	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/client/db"
	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/model"
	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/repository"
	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/repository/note/converter"
	modelRepo "github.com/vivaldi7/golang_code/mic_curce/week_3/internal/repository/note/model"
)

const (
	tableName       = "note"
	idColumn        = "id"
	titleColumn     = "title"
	contentColumn   = "content"
	craetedAtColumn = "created_at"
	updateAtColumn  = "update_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.NoteRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(titleColumn, contentColumn).
		Values(info.Title, info.Content).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		//		log.Fatalf("failed to build query: %v", err)
		return 0, err
	}

	q := db.Query{
		Name:     "note_repositore.Create",
		QueryRaw: query,
	}

	var id int64

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		//		log.Fatalf("failed to insert notes: %v", err)
		return 0, err
	}
	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.Note, error) {
	builder := sq.Select(idColumn, titleColumn, contentColumn, craetedAtColumn, updateAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		//		log.Fatalf("Failed to build query: %v", err)
		return nil, err
	}

	q := db.Query{
		Name:     "note_repositore.Get",
		QueryRaw: query,
	}

	var note modelRepo.Note

	err = r.db.DB().ScanOneContext(ctx, &note, q, args...)
	if err != nil {
		//		log.Fatalf("Failed to select note: %v", err)
		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil
}
