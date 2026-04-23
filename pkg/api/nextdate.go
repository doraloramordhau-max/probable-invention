package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

var errInvalidRepeat = errors.New("invalid repeat format")

func afterNow(date, now time.Time) bool {
	dy, dm, dd := date.Date()
	ny, nm, nd := now.Date()
	return time.Date(dy, dm, dd, 0, 0, 0, 0, date.Location()).After(
		time.Date(ny, nm, nd, 0, 0, 0, 0, now.Location()),
	)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errInvalidRepeat
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errInvalidRepeat
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errInvalidRepeat
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", errInvalidRepeat
		}

		date = date.AddDate(0, 0, interval)
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, interval)
		}

		return date.Format(dateFormat), nil
	case "y":
		if len(parts) != 1 {
			return "", errInvalidRepeat
		}

		date = date.AddDate(1, 0, 0)
		for !afterNow(date, now) {
			date = date.AddDate(1, 0, 0)
		}

		return date.Format(dateFormat), nil
	default:
		return "", errInvalidRepeat
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now := time.Now()
	var err error
	if nowStr != "" {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
