package rabbitmq_mock

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/Moranilt/http_template/clients/rabbitmq"
	"github.com/Moranilt/http_template/utils/mock"
)

var rabbitPushTests = []struct {
	name         string
	expectedData []byte
	actualData   []byte
	expectedErr  func(name string) error
	returnErr    error
	runExpect    bool
}{
	{
		name:         "valid push",
		expectedData: []byte("expected data"),
		actualData:   []byte("expected data"),
		expectedErr:  nil,
		returnErr:    nil,
		runExpect:    true,
	},
	{
		name:         "expected push to return error",
		expectedData: []byte("expected data"),
		actualData:   []byte("expected data"),
		expectedErr:  func(name string) error { return errors.New("expected error") },
		returnErr:    errors.New("expected error"),
		runExpect:    true,
	},
	{
		name:         "unexpected call of Push",
		expectedData: nil,
		actualData:   nil,
		expectedErr:  func(name string) error { return fmt.Errorf(mock.ERR_Events_Is_Empty, name) },
		returnErr:    nil,
		runExpect:    false,
	},
	{
		name:         "not valid body data",
		expectedData: []byte("expected body"),
		actualData:   []byte("actual body"),
		expectedErr: func(name string) error {
			return fmt.Errorf(ERR_Not_Valid_Data, name, []byte("expected body"), []byte("actual body"))
		},
		returnErr: nil,
		runExpect: true,
	},
}

var rabbitCloseTests = []struct {
	name        string
	expectedErr error
	mockedErr   error
	runExpect   bool
}{
	{
		name:        "valid close",
		expectedErr: nil,
		mockedErr:   nil,
		runExpect:   true,
	},
	{
		name:        "mocked error",
		expectedErr: errors.New("expected error"),
		mockedErr:   errors.New("expected error"),
		runExpect:   true,
	},
	{
		name:        "unexpected call",
		expectedErr: fmt.Errorf(mock.ERR_Events_Is_Empty, "Close"),
		mockedErr:   nil,
		runExpect:   false,
	},
}

func TestRabbitMQ(t *testing.T) {
	ctx := context.Background()

	t.Run("Push", func(t *testing.T) {
		for _, test := range rabbitPushTests {
			t.Run(test.name, func(t *testing.T) {
				mockRabbit := NewRabbitMQ(t)

				if test.runExpect {
					mockRabbit.ExpectPush(test.expectedData, test.returnErr)
				}

				var expectedErr error
				if test.expectedErr != nil {
					expectedErr = test.expectedErr("Push")
				}

				err := mockRabbit.Push(ctx, test.actualData)
				if err != nil && err.Error() != expectedErr.Error() {
					t.Errorf("unexpected error %q, expect %q", err, expectedErr)
				}

				if err := mockRabbit.AllExpectationsDone(); err != nil {
					t.Error(err)
				}
			})
		}
	})

	t.Run("UnsafePush", func(t *testing.T) {
		for _, test := range rabbitPushTests {
			t.Run(test.name, func(t *testing.T) {
				mockRabbit := NewRabbitMQ(t)

				if test.runExpect {
					mockRabbit.ExpectUnsafePush(test.expectedData, test.returnErr)
				}

				var expectedErr error
				if test.expectedErr != nil {
					expectedErr = test.expectedErr("UnsafePush")
				}

				err := mockRabbit.UnsafePush(ctx, test.actualData)
				if err != nil && err.Error() != expectedErr.Error() {
					t.Errorf("unexpected error %q, expect %q", err, expectedErr)
				}

				if err := mockRabbit.AllExpectationsDone(); err != nil {
					t.Error(err)
				}
			})
		}
	})

	t.Run("Close", func(t *testing.T) {
		for _, test := range rabbitCloseTests {
			mockRabbit := NewRabbitMQ(t)
			if test.runExpect {
				mockRabbit.ExpectClose(test.mockedErr)
			}

			err := mockRabbit.Close()
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expect error %q, got %q", err, test.expectedErr)
			}

			if err := mockRabbit.AllExpectationsDone(); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("Consume", func(t *testing.T) {
		mockRabbit := NewRabbitMQ(t)
		expectedData := []byte("expected data")
		mockRabbit.ExpectPush(expectedData, nil)
		mockRabbit.ExpectConsume(nil)

		err := mockRabbit.Push(ctx, expectedData)
		if err != nil {
			t.Errorf("not expected error %q", err)
		}

		dataCh, err := mockRabbit.Consume()
		if err != nil {
			t.Errorf("error should be empty, got %q", err)
		}

		data := <-dataCh
		if !reflect.DeepEqual(data.Body, expectedData) {
			t.Errorf("expected %v, got %v", string(expectedData), string(data.Body))
		}
	})

	t.Run("ReadMsg", func(t *testing.T) {
		mockRabbit := NewRabbitMQ(t)
		mockRabbit.ExpectReadMsg()
		mockRabbit.ReadMsgs(ctx, 1, 1*time.Second, func(ctx context.Context, d rabbitmq.RabbitDelivery) error { return nil })
		if err := mockRabbit.AllExpectationsDone(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Check", func(t *testing.T) {
		mockRabbit := NewRabbitMQ(t)
		mockRabbit.ExpectCheck(nil)
		err := mockRabbit.Check(ctx)
		if err != nil {
			t.Errorf("got not expected error %q", err)
		}

		if err := mockRabbit.AllExpectationsDone(); err != nil {
			t.Error(err)
		}
	})
}
