package store

import (
	st "github.com/endo-checker/common/store"
	loggerv1 "github.com/endo-checker/logger/internal/gen/logger/v1"
)

func New(uri string) LoggerStore {
	s := st.Connect[*loggerv1.Log](uri)

	return LoggerStore{
		Store: &s,
	}
}
