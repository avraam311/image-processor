package images

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wb-go/wbf/dbpg"
)

func TestRepository_SetImageStatus(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		mockSetup   func(sqlmock.Sqlmock)
		expectedID  uint
		expectError bool
	}{
		{
			name:   "success",
			status: "in process",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO image \(status\) VALUES \(\$1\) RETURNING id`).
					WithArgs("in process").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedID:  1,
			expectError: false,
		},
		{
			name:   "db error",
			status: "in process",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO image \(status\) VALUES \(\$1\) RETURNING id`).
					WithArgs("in process").
					WillReturnError(errors.New("db error"))
			},
			expectedID:  0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			repo := &Repository{db: &dbpg.DB{Master: db}}

			id, err := repo.SetImageStatus(context.Background(), tt.status)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedID, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_CheckImage(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		mockSetup   func(sqlmock.Sqlmock)
		expectError error
	}{
		{
			name: "success - processed",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT status FROM image WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow("processed"))
			},
			expectError: nil,
		},
		{
			name: "not found",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT status FROM image WHERE id = \$1`).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			expectError: ErrImageNotFound,
		},
		{
			name: "in process",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT status FROM image WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow("in process"))
			},
			expectError: ErrImageInProcess,
		},
		{
			name: "db error",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT status FROM image WHERE id = \$1`).
					WithArgs(1).
					WillReturnError(errors.New("db error"))
			},
			expectError: errors.New("repository/check_image.go - failed to check image - db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			repo := &Repository{db: &dbpg.DB{Master: db}}

			err = repo.CheckImage(context.Background(), tt.id)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_DeleteImage(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		mockSetup   func(sqlmock.Sqlmock)
		expectError error
	}{
		{
			name: "success",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM image WHERE id = \$1`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: nil,
		},
		{
			name: "not found",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM image WHERE id = \$1`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectError: ErrImageNotFound,
		},
		{
			name: "db error",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM image WHERE id = \$1`).
					WithArgs(1).
					WillReturnError(errors.New("db error"))
			},
			expectError: errors.New("repository/change_image_status.go - failed to delete image - db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			repo := &Repository{db: &dbpg.DB{Master: db}}

			err = repo.DeleteImage(context.Background(), tt.id)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_ChangeImageStatus(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		status      string
		mockSetup   func(sqlmock.Sqlmock)
		expectError error
	}{
		{
			name:   "success",
			id:     1,
			status: "processed",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE image SET status = \$2 WHERE id = \$1`).
					WithArgs(1, "processed").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: nil,
		},
		{
			name:   "not found",
			id:     1,
			status: "processed",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE image SET status = \$2 WHERE id = \$1`).
					WithArgs(1, "processed").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectError: ErrImageNotFound,
		},
		{
			name:   "db error",
			id:     1,
			status: "processed",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE image SET status = \$2 WHERE id = \$1`).
					WithArgs(1, "processed").
					WillReturnError(errors.New("db error"))
			},
			expectError: errors.New("repository/change_image_status.go - failed to change image status - db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			repo := &Repository{db: &dbpg.DB{Master: db}}

			err = repo.ChangeImageStatus(context.Background(), tt.id, tt.status)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
