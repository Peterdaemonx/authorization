package fixture

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/spanner"
)

func Seed(ctx context.Context, client *spanner.Client) error {
	statements, err := getStatements()
	if err != nil {
		return fmt.Errorf("failed to fetch statements: %w", err)
	}

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		rowCounts, err := txn.BatchUpdate(ctx, statements)
		if err != nil {
			return fmt.Errorf("failed to execute transactions: %w", err)
		}

		fmt.Fprintf(os.Stdout, "Executed %d SQL statements using Batch DML.\n", len(rowCounts))
		return nil
	})

	return err
}

//go:embed seed.sql
var f embed.FS

func getStatements() ([]spanner.Statement, error) {
	var statements []spanner.Statement

	file, err := f.ReadFile("seed.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to open seed.sql, %w", err)
	}

	// Start reading from the file with a reader.
	reader := bufio.NewReader(bytes.NewReader(file))

	var line string
	// read the line until the semicolon and then append it to the statements
	for {
		line, err = reader.ReadString(';')
		if err != nil {
			break
		}
		statements = append(statements, spanner.Statement{SQL: line})

	}

	if err != io.EOF {
		return nil, err
	}

	return statements, nil
}
