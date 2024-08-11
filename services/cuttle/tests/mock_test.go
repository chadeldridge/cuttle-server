package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockNewMockTest(t *testing.T) {
	require := require.New(t)

	t.Run("pass", func(t *testing.T) {
		test := NewMockTest(false)
		require.Equal("Mock Test", test.Name, "NewMockTest() Name is not 'Mock Test'")
		require.True(test.MustSucceed, "NewMockTest() MustSucceed is not true")
		require.False(test.Tester.(*MockTest).fail, "NewMockTest() fail is not false")
	})

	t.Run("fail", func(t *testing.T) {
		test := NewMockTest(true)
		require.Equal("Mock Test", test.Name, "NewMockTest() Name is not 'Mock Test'")
		require.True(test.MustSucceed, "NewMockTest() MustSucceed is not true")
		require.True(test.Tester.(*MockTest).fail, "NewMockTest() fail is not true")
	})
}

func TestMockRun(t *testing.T) {
	require := require.New(t)
	server := testServerSetup(t)
	test := MockTest{fail: false}

	t.Run("pass", func(t *testing.T) {
		require.NoError(test.Run(server), "MockTest.Run() returned an error")
	})

	t.Run("fail", func(t *testing.T) {
		test.fail = true
		require.Equal(ErrTestFailed, test.Run(server), "MockTest.Run() did not return ErrTestFailed")
	})
}
