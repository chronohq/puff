// Copyright Chrono Technologies LLC
// SPDX-License-Identifier: MIT

package main

import (
	"time"

	"github.com/google/uuid"
)

func generateUUIDV4() (uuid.UUID, error) {
	return uuid.NewRandom()
}

func generateUUIDV7(customTime string) (uuid.UUID, error) {
	var err error
	var ret uuid.UUID
	var timestamp time.Time

	if ret, err = uuid.NewV7(); err != nil {
		return uuid.UUID{}, err
	}

	if len(customTime) == 0 {
		return ret, nil
	}

	if timestamp, err = parseTimeInput(customTime); err != nil {
		return uuid.UUID{}, err
	}

	// replace the first 6-bytes with the custom timestamp
	msec := timestamp.UnixMilli()

	ret[0] = byte(msec >> 40)
	ret[1] = byte(msec >> 32)
	ret[2] = byte(msec >> 24)
	ret[3] = byte(msec >> 16)
	ret[4] = byte(msec >> 8)
	ret[5] = byte(msec)

	return ret, nil
}
