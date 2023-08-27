package repository

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/textproto"
	"reflect"
	"testing"

	"github.com/Moranilt/http_template/models"
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

		if response.Name != "Test" {
			t.Errorf("Expected name to be 'Test', got %q", response.Name)
		}

		if !reflect.DeepEqual(response.Files, mockedFiles) {
			t.Errorf("Expected files to be %v, got %v", mockedFiles, response.Files)
		}

		if !reflect.DeepEqual(response.OneMoreFile, mockedFile) {
			t.Errorf("Expected oneMoreFile to be %v, got %v", mockedFile, response.OneMoreFile)
		}
	})

	t.Run("empty request", func(t *testing.T) {
		response, err := mockedRepo.repo.Files(context.Background(), nil)
		if err.Error() != ERR_BodyRequired {
			t.Errorf("Expected error %q but got %q", ERR_BodyRequired, err)
		}

		if response != nil {
			t.Errorf("Expected nil response on error, got %v", response)
		}
	})

	t.Run("rabbitmq error", func(t *testing.T) {
		expectedError := errors.New("rabbitmq error")
		b, _ := json.Marshal(mockedRequest)
		mockedRepo.rabbitmqMock.ExpectPush(b, expectedError)

		response, err := mockedRepo.repo.Files(context.Background(), mockedRequest)
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected error %q but got %q", ERR_BodyRequired, err)
		}

		if response != nil {
			t.Errorf("Expected nil response on error, got %v", response)
		}
	})
}
