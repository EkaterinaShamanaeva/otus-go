package main

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestConnectWithoutErrors(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockTelnetClient := NewMockTelnetClient(controller)

	mockTelnetClient.EXPECT().Connect().Return(nil)
	mockTelnetClient.EXPECT().Send().Return(nil)
	mockTelnetClient.EXPECT().Receive().Return(nil)
	mockTelnetClient.EXPECT().Close().Return(nil)
	require.NoError(t, runTelnet(mockTelnetClient))
}

func TestConnectWithError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockTelnetClient := NewMockTelnetClient(controller)
	mockTelnetClient.EXPECT().Connect().Return(errors.New("error"))
	require.Error(t, runTelnet(mockTelnetClient))
}
