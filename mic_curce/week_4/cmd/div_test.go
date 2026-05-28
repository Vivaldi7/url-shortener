package main

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiv(t *testing.T) {
	a := gofakeit.Float64()
	b := gofakeit.Float64()

	require.Equal(t, a/b, Div(a, b))

	// Тестирование assert.Equal то же амое что и require.Equal но при первой ошибке теста  проверка не завершается
	//	assert.Equal(t, a/b, Div(a, b))

	//Строки ниже заменяет пакет github.com/stretchr/testify/require
	//	if a/b != Div(a, b) {
	//		t.Errorf("Div(%f, %f) = %f; Want %f", a, b, Div(a, b), a/b)
	//	}
}

func TestDiv2(t *testing.T) {
	t.Run("10/5", func(t *testing.T) {
		a := 10.0
		b := 5.0

		expected := 2.0

		require.Equal(t, expected, Div(a, b))
	})

	t.Run("10/0", func(t *testing.T) {
		a := 10.0
		b := 0.0

		expected := 0.0

		require.Equal(t, expected, Div(a, b))
	})
}
