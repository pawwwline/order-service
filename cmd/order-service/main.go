package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/pressly/goose/v3"
)

func main() {
	var _ = dummy()
}

func dummy() int {
	return 0
}
