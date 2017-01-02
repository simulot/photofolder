package mov

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"time"
)

const epochAdjust = 2082844800

// Created attempts to find a created time from the metadata in a .mov file.
func Created(file io.ReadSeeker) (time.Time, error) {
	c, _, err := times(file)
	return c, err
}

// Modified attempts to find a modification time from the metadata in a .mov
// file.
func Modified(file io.ReadSeeker) (time.Time, error) {
	_, m, err := times(file)
	return m, err
}

// times seeks around in file and finds Created and Modified times.
//
// This was transcibed from http://stackoverflow.com/a/21395803
func times(file io.ReadSeeker) (time.Time, time.Time, error) {
	var c, m time.Time
	var i int64
	var err error

	buf := [8]byte{}
	for {
		_, err := file.Read(buf[:])
		if err != nil {
			return c, m, err
		}
		if bytes.Equal(buf[4:8], []byte("moov")) {
			break
		} else {
			atomSize := binary.BigEndian.Uint32(buf[:])
			file.Seek(int64(atomSize)-8, 1)
		}
	}

	_, err = file.Read(buf[:])
	if err != nil {
		return c, m, err
	}

	s := string(buf[4:8])
	switch s {
	case "mvhd":
		if _, err := file.Seek(4, 1); err != nil {
			return c, m, err
		}

		_, err = file.Read(buf[:4])
		if err != nil {
			return c, m, err
		}

		i = int64(binary.BigEndian.Uint32(buf[:4]))
		c := time.Unix(i-epochAdjust, 0).Local()

		_, err = file.Read(buf[:4])
		if err != nil {
			return c, m, err
		}

		i = int64(binary.BigEndian.Uint32(buf[:4]))
		m := time.Unix(i-epochAdjust, 0).Local()

		return c, m, nil
	case "cmov":
		return c, m, errors.New("moov atom is compressed")
	default:
		return c, m, errors.New("expected to find 'mvhd' header, didn't")
	}
}
