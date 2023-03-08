package easylog

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AReader struct {
	max   int
	soFar int
}

func (ar *AReader) Read(b []byte) (int, error) {
	var i int
	for i = 0; i < len(b) && i+ar.soFar < ar.max; i++ {
		b[i] = 'A'
	}
	ar.soFar += i
	if ar.soFar >= ar.max {
		return i, io.EOF
	}
	return i, nil
}

func setUpTest(t *testing.T) {
	// easylog prints to stdout, but that's just noise in the test output,
	// so redirect to DevNull.
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	require.NoError(t, err)
	stdout = devNull
	t.Cleanup(func() {
		stdout = os.Stdout
	})

	logFilePath := filepath.Join(t.TempDir(), fmt.Sprint(t.Name(), ".log"))
	assert.NoError(t, Init(logFilePath))
}

func TestBasic(t *testing.T) {
	setUpTest(t)

	SetLevel(INFO)
	assert.Equal(t, INFO, GetLevel())

	Debug("AAA")
	Debugf("%s", "aaa")
	Info("BBB")
	Infof("%s", "bbb")
	Warn("CCC")
	Warnf("%s", "ccc")
	Error("DDD")
	Errorf("%s", "ddd")
	Fatal("EEE")
	Fatalf("%s", "eee")

	logFile, err := os.Open(filePath)
	require.NoError(t, err)
	var foundInfo, foundWarn, foundError, foundFatal bool
	var foundInfof, foundWarnf, foundErrorf, foundFatalf bool
	for scnr := bufio.NewScanner(logFile); scnr.Scan(); {
		line := scnr.Text()
		assert.NotContains(t, line, "AAA")
		assert.NotContains(t, line, "aaa")
		assert.Contains(t, line, "easylog_test.go")
		if strings.Contains(line, "BBB") {
			foundInfo = true
			assert.Contains(t, line, "Info")
		} else if strings.Contains(line, "CCC") {
			foundWarn = true
			assert.Contains(t, line, "Warn")
		} else if strings.Contains(line, "DDD") {
			foundError = true
			assert.Contains(t, line, "Error")
		} else if strings.Contains(line, "EEE") {
			foundFatal = true
			assert.Contains(t, line, "Fatal")
		} else if strings.Contains(line, "bbb") {
			foundInfof = true
			assert.Contains(t, line, "Info")
		} else if strings.Contains(line, "ccc") {
			foundWarnf = true
			assert.Contains(t, line, "Warn")
		} else if strings.Contains(line, "ddd") {
			foundErrorf = true
			assert.Contains(t, line, "Error")
		} else if strings.Contains(line, "eee") {
			foundFatalf = true
			assert.Contains(t, line, "Fatal")
		}
	}
	assert.True(t, foundInfo)
	assert.True(t, foundWarn)
	assert.True(t, foundError)
	assert.True(t, foundFatal)
	assert.True(t, foundInfof)
	assert.True(t, foundWarnf)
	assert.True(t, foundErrorf)
	assert.True(t, foundFatalf)
	require.NoError(t, logFile.Close())
}

func TestMaxSize(t *testing.T) {
	setUpTest(t)

	SetMaxSize(0)
	assert.Equal(t, defaultMaxSize, maxSize)

	// Check that the log file is rotated when it reaches the max size.
	logFile, err := os.OpenFile(filePath, os.O_WRONLY, 0666)
	require.NoError(t, err)
	const maxSize = 1 << (10 * 2)
	SetMaxSize(maxSize >> (10 * 2))
	written, err := io.Copy(logFile, &AReader{max: maxSize})
	require.NoError(t, err)
	assert.Equal(t, int64(maxSize), written)
	info, err := logFile.Stat()
	require.NoError(t, err)
	assert.Equal(t, int64(maxSize), info.Size())
	require.NoError(t, logFile.Close())

	Info("FFF")

	logFile, err = os.Open(filePath)
	require.NoError(t, err)
	var numLines int
	for scnr := bufio.NewScanner(logFile); scnr.Scan(); {
		line := scnr.Text()
		numLines += 1
		assert.Contains(t, line, "FFF")
	}
	assert.Equal(t, 1, numLines)
}

func TestRace(t *testing.T) {
	setUpTest(t)

	go func() {
		for i := 0; i < 1000; i++ {
			Info("AAA")
		}
	}()
	for i := 0; i < 1000; i++ {
		Info("BBB")
	}
}
