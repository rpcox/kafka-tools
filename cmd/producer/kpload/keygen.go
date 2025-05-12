package main

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

type KeyGener interface {
	KeyGen() string
}

type Uuid struct {
	Type string
}

func (u Uuid) KeyGen() string {
	return uuid.New().String()
}

type Epoch struct {
	Type string
}

func (e Epoch) KeyGen() string {
	epoch := time.Now().Unix()
	return strconv.FormatInt(epoch, 10)
}

var enumeratorCount int64 = 1

type Enumerator struct {
	Type string
}

func (e Enumerator) KeyGen() string {
	return strconv.FormatInt(enumeratorCount, 10)
}

func GetKeyGenerator(s string) KeyGener {
	switch s {
	case "enum":
		return Enumerator{Type: "enum"}
	case "epoch":
		return Epoch{Type: "epoch"}
	case "uuid":
		return Uuid{Type: "uuid"}
	}

	return nil
}
