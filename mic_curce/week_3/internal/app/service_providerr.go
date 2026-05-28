package app

import (
	"context"
	"log"

	//	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/api/note"
	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/client/db"
	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/client/db/pg"
	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/client/db/transaction"
	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/closer"
	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/config"
	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/repository"
	noteRepositore "github.com/vivaldi7/golang_code/mic_curce/week_3/internal/repository/note"
	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/service"

	//	note "github.com/vivaldi7/golang_code/mic_curce/week_3/internal/api/note"
	//	noteRepositore "github.com/vivaldi7/golang_code/mic_curce/week_3/internal/repository/note"
	//	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/repository/note"
	//	noteRepositore "github.com/vivaldi7/golang_code/mic_curce/week_3/internal/repository/note"
	//	noteService "github.com/vivaldi7/golang_code/mic_curce/week_3/internal/service/note"
	//	"github.com/vivaldi7/golang_code/mic_curce/week_3/internal/service/note"
	noteService "github.com/vivaldi7/golang_code/mic_curce/week_3/internal/service/note"
	//	noteService "github.com/vivaldi7/golang_code/mic_curce/week_3/internal/service/note"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient  db.Client
	txManager db.TxManager
	//	pgPool     *pgxpool.Pool

	noteRepository repository.NoteRepository
	noteService    service.NoteService
	noteImpl       *note.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}
		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = cfg
	}
	return s.grpcConfig
}

func (s *serviceProvider) DBCClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}
	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBCClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) NoteRepositore(ctx context.Context) repository.NoteRepository {
	if s.noteRepository == nil {
		s.noteRepository = noteRepositore.NewRepository(s.DBCClient(ctx))
	}
	return s.noteRepository
}

func (s *serviceProvider) NoteService(ctx context.Context) service.NoteService {
	if s.noteService == nil {
		s.noteService = noteService.NewService(
			s.NoteRepositore(ctx),
			s.TxManager(ctx),
		)
	}
	return s.noteService
}

func (s *serviceProvider) GetnoteImpl(ctx context.Context) *note.Implementation {
	if s.noteImpl == nil {
		s.noteImpl = note.NewImplementation(s.NoteService(ctx))
	}
	return s.noteImpl
}
