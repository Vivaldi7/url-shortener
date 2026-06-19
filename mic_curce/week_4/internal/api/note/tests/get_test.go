package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"github.com/vivaldi7/golang_code/mic_curce/week_4/internal/api/note"
	"github.com/vivaldi7/golang_code/mic_curce/week_4/internal/model"
	"github.com/vivaldi7/golang_code/mic_curce/week_4/internal/service"
	serviceMocks "github.com/vivaldi7/golang_code/mic_curce/week_4/internal/service/mocks"
	desc "github.com/vivaldi7/golang_code/mic_curce/week_4/pkg/note_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGet(t *testing.T) {
	type noteServiseMockFunc func(mc *minimock.Controller) service.NoteService

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		title     = gofakeit.Animal()
		content   = gofakeit.Animal()
		createdAt = gofakeit.Date()
		updateAt  = gofakeit.Date()

		serviceErr = fmt.Errorf("service error")

		req = &desc.GetRequest{
			Id: id,
		}

		serviceRes = &model.Note{
			ID: id,
			Info: model.NoteInfo{
				Title:   title,
				Content: content,
			},
			CreatedAt: createdAt,
			UpdateAt: sql.NullTime{
				Time:  updateAt,
				Valid: true,
			},
		}

		res = &desc.GetResponse{
			Note: &desc.Note{
				Id: id,
				Info: &desc.NoteInfo{
					Title:   title,
					Content: content,
				},
				CreatedAt: timestamppb.New(createdAt),
				UpdateAt:  timestamppb.New(updateAt),
			},
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetResponse
		err             error
		noteServiceMock noteServiseMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			noteServiceMock: func(mc *minimock.Controller) service.NoteService {
				mock := serviceMocks.NewNoteServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(serviceRes, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			noteServiceMock: func(mc *minimock.Controller) service.NoteService {
				mock := serviceMocks.NewNoteServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			noteServiceMock := tt.noteServiceMock(mc)
			api := note.NewImplementation(noteServiceMock)

			resHandler, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resHandler)
		})
	}
}
