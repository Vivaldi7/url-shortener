package note

import (
	"github.com/vivaldi7/golang_code/mic_curce/week_4/internal/client/db"
	"github.com/vivaldi7/golang_code/mic_curce/week_4/internal/repository"
	"github.com/vivaldi7/golang_code/mic_curce/week_4/internal/service"
	def "github.com/vivaldi7/golang_code/mic_curce/week_4/internal/service"
)

var _ def.NoteService = (*serv)(nil)

type serv struct {
	noteRepository repository.NoteRepository
	txManager      db.TxManager
}

func NewService(
	noteRepository repository.NoteRepository,
	txManager db.TxManager,
) service.NoteService {
	return &serv{
		noteRepository: noteRepository,
		txManager:      txManager,
	}
}
