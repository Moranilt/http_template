package repository

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/Moranilt/http_template/custom_errors"
	"github.com/Moranilt/http_template/models"
	"github.com/stretchr/testify/assert"
)

func TestFiles(t *testing.T) {
	// Mock repository dependencies
	mockedRepo := mockRepository(t)
	mockedFiles := []*multipart.FileHeader{
		{
			Filename: "file.txt",
			Header: textproto.MIMEHeader{
				"Content-Type": {"text/plain"},
			},
			Size: 123,
		},
	}
	mockedFile := &multipart.FileHeader{
		Filename: "test.txt",
		Header: textproto.MIMEHeader{
			"Content-Type": {"text/plain"},
		},
		Size: 123,
	}
	mockedRequest := &models.FileRequest{
		Name:        "Test",
		Files:       mockedFiles,
		OneMoreFile: mockedFile,
	}
	t.Run("Success", func(t *testing.T) {
		b, _ := json.Marshal(mockedRequest)
		mockedRepo.rabbitmqMock.ExpectPush(b, nil)

		response, err := mockedRepo.repo.Files(context.Background(), &models.FileRequest{
			Name:        "Test",
			Files:       mockedFiles,
			OneMoreFile: mockedFile,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		assert.Equal(t, "Test", response.Name)
		assert.Equal(t, mockedFiles, response.Files)
		assert.Equal(t, mockedFile, response.OneMoreFile)
	})

	t.Run("empty request", func(t *testing.T) {
		response, err := mockedRepo.repo.Files(context.Background(), nil)
		assert.Equal(t, custom_errors.ERR_CODE_BodyRequired, err.GetCode())
		assert.Nil(t, response)
	})

	t.Run("rabbitmq error", func(t *testing.T) {
		expectedError := errors.New("rabbitmq error")
		b, _ := json.Marshal(mockedRequest)
		mockedRepo.rabbitmqMock.ExpectPush(b, expectedError)

		response, err := mockedRepo.repo.Files(context.Background(), mockedRequest)
		assert.Equal(t, expectedError.Error(), err.Error())
		assert.Nil(t, response)
	})
}
