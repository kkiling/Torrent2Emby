package statemachine

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStepContext_GetOptions(t *testing.T) {
	type testOptions struct {
		Title  string
		Amount int
	}

	type wrongOptions struct {
		Name string
	}

	// Успешное получение опций
	t.Run("success get options", func(t *testing.T) {
		completeOptions := testOptions{
			Title:  "some title",
			Amount: 234,
		}

		sc := StepContext[string, interface{}, interface{}, string, string]{
			completeOptionsType: reflect.TypeOf(testOptions{}),
			completeOptions:     completeOptions,
		}

		opts := testOptions{}
		ok, err := sc.GetOptions(&opts)
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, completeOptions, opts)
	})

	// Опции ожидаются, но не установлены
	t.Run("not set options", func(t *testing.T) {
		sc := StepContext[string, interface{}, interface{}, string, string]{
			completeOptionsType: reflect.TypeOf(testOptions{}),
		}

		opts := testOptions{}
		ok, err := sc.GetOptions(&opts)
		require.NoError(t, err)
		require.False(t, ok)
	})

	// Опции не установлены и тип не задан
	t.Run("options type not set", func(t *testing.T) {
		sc := StepContext[string, interface{}, interface{}, string, string]{}

		opts := testOptions{}
		ok, err := sc.GetOptions(&opts)
		require.NoError(t, err)
		require.False(t, ok)
	})

	// Опции установлены, но тип не задан (ошибка)
	t.Run("options set but type not set", func(t *testing.T) {
		sc := StepContext[string, interface{}, interface{}, string, string]{
			completeOptions: testOptions{Title: "test"},
		}

		opts := testOptions{}
		ok, err := sc.GetOptions(&opts)
		require.Error(t, err)
		require.Equal(t, "completeOptions is set but completeOptionsType is not defined", err.Error())
		require.False(t, ok)
	})

	// Несоответствие типов опций
	t.Run("type mismatch", func(t *testing.T) {
		sc := StepContext[string, interface{}, interface{}, string, string]{
			completeOptionsType: reflect.TypeOf(testOptions{}),
			completeOptions:     wrongOptions{Name: "wrong"},
		}

		opts := testOptions{}
		ok, err := sc.GetOptions(&opts)
		require.Error(t, err)
		require.Contains(t, err.Error(), "type mismatch: expected")
		require.False(t, ok)
	})

	// Неверный приемник (не указатель)
	t.Run("invalid receiver (not pointer)", func(t *testing.T) {
		sc := StepContext[string, interface{}, interface{}, string, string]{
			completeOptionsType: reflect.TypeOf(testOptions{}),
			completeOptions:     testOptions{Title: "test"},
		}

		opts := testOptions{}
		ok, err := sc.GetOptions(opts) // передаем по значению
		require.Error(t, err)
		require.Equal(t, "destination must be a non-nil pointer", err.Error())
		require.False(t, ok)
	})

	// Неверный приемник (nil указатель)
	t.Run("invalid receiver (nil pointer)", func(t *testing.T) {
		sc := StepContext[string, interface{}, interface{}, string, string]{
			completeOptionsType: reflect.TypeOf(testOptions{}),
			completeOptions:     testOptions{Title: "test"},
		}

		var opts *testOptions = nil
		ok, err := sc.GetOptions(opts)
		require.Error(t, err)
		require.Equal(t, "destination must be a non-nil pointer", err.Error())
		require.False(t, ok)
	})

	// Несовместимые типы при присваивании
	t.Run("incompatible types", func(t *testing.T) {
		sc := StepContext[string, interface{}, interface{}, string, string]{
			completeOptionsType: reflect.TypeOf(testOptions{}),
			completeOptions:     testOptions{Title: "test"},
		}

		var wrong wrongOptions
		ok, err := sc.GetOptions(&wrong)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot assign")
		require.False(t, ok)
	})
}
