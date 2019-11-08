package utils

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"os"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

type shasumTy struct{}

var SHASum = shasumTy{}

func readerAsString(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	b := buf.Bytes()
	s := *(*string)(unsafe.Pointer(&b))
	return s
}

func (shasumTy) Check(reader io.Reader, expected string) (bool, error) {
	in := readerAsString(reader)

	hash, size, _ := splitSum(expected)

	if size > 0 && int64(len(in)) != size {
		return false, nil
	}

	switch len(hash) {
	case 64:
		return sha256sum(in) == hash, nil
	case 128:
		return sha512sum(in) == hash, nil
	case 40:
		return sha1sum(in) == hash, nil
	case 0:
		return true, nil // if no checksum assume valid
	}

	return false, errors.New("failed to perform sha hash check")
}

func (shasumTy) CheckFile(path string, expected string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, errors.Wrapf(err, "failed to open %s while performing sha checksum", path)
	}
	defer f.Close()
	ok, err := SHASum.Check(f, expected)
	if err != nil {
		return false, errors.Wrapf(err, "unable to perform shasum on %s", path)
	}
	return ok, nil
}

func sha1sum(in string) string {
	h := sha1.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func sha256sum(in string) string {
	h := sha256.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func sha512sum(in string) string {
	h := sha512.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func splitSum(in string) (string, int64, string) {
	var hash string
	var name string
	var size int64

	// the checksum might be split into multiple
	// sections including the file size and name.
	switch strings.Count(in, " ") {
	case 1:
		fmt.Sscanf(in, "%s %s", &hash, &name)
	case 2:
		fmt.Sscanf(in, "%s %d %s", &hash, &size, &name)
	default:
		hash = in
	}

	return hash, size, name
}
