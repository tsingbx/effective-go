package main

import "time"

type reqTimeout time.Duration

func toReqTimeout(p *time.Duration) *reqTimeout {
	return (*reqTimeout)(p)
}

func (t *reqTimeout) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*t = (reqTimeout)(v)
	return nil
}

func (t *reqTimeout) String() string {
	v := time.Duration(*t)
	return v.String()
}
