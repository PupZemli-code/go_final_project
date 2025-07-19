package tests

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/PupZemli-code/go-final-project/go_final_project/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getURL(path string) string {
	port := Port
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {
		if eport, err := strconv.ParseInt(envPort, 10, 32); err == nil {
			port = int(eport)
		}
	}
	path = strings.TrimPrefix(strings.ReplaceAll(path, `\`, `/`), `../web/`)
	return fmt.Sprintf("http://localhost:%d/%s", port, path)
}

func getBody(path string) ([]byte, error) {
	resp, err := http.Get(getURL(path))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, err
}

func walkDir(path string, f func(fname string) error) error {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, v := range dirs {
		fname := filepath.Join(path, v.Name())
		if v.IsDir() {
			if err = walkDir(fname, f); err != nil {
				return err
			}
			continue
		}
		if err = f(fname); err != nil {
			return err
		}
	}
	return nil
}

func TestApp(t *testing.T) {
	cmp := func(fname string) error {
		fbody, err := os.ReadFile(fname)
		if err != nil {
			return err
		}
		body, err := getBody(fname)
		if err != nil {
			return err
		}
		assert.Equal(t, len(fbody), len(body), `сервер возвращает для %s данные другого размера`, fname)
		return nil
	}
	assert.NoError(t, walkDir("../web", cmp))
}

func TestGetAddr(t *testing.T) {
	// Очищаем переменные
	os.Unsetenv("TODO_HOST")
	os.Unsetenv("TODO_PORT")
	if err := os.Setenv("TODO_HOST", "test"); err != nil {
		log.Fatal("ошибка установки TODO_HOST", err)
	}
	if err := os.Setenv("TODO_PORT", "test"); err != nil {
		log.Fatal("ошибка установки TODO_PORT", err)
	}
	require.Equal(t, server.GetAddr(), "test:test")
}
