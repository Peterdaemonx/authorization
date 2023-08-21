package spanner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/option"
)

func NewTestDB(t *testing.T, file string) *spanner.Client {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		os.Setenv("SPANNER_EMULATOR_HOST", "localhost:10010")
	}

	ctx := context.Background()

	dsn := fmt.Sprintf("projects/%s/instances/%s/databases/%s", "cc-acquiring-development", "acquiring-instance", "authorizations")

	client, err := spanner.NewClient(ctx, dsn, option.WithGRPCConnectionPool(1))
	if err != nil {
		t.Fatal(err)
	}

	if file != "" {
		err = seed(ctx, client, fmt.Sprintf("%s_up.sql", file))
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Cleanup(func() {
		if file != "" {
			err = seed(ctx, client, fmt.Sprintf("%s_down.sql", file))
			if err != nil {
				t.Fatal(err)
			}
		}
		client.Close()
	})

	return client
}

func seed(ctx context.Context, client *spanner.Client, file string) error {
	statements, err := getStatements(file)
	if err != nil {
		return fmt.Errorf("failed to fetch statements: %w", err)
	}

	if len(statements) != 0 {
		_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
			_, err := txn.BatchUpdate(ctx, statements)
			if err != nil {
				return fmt.Errorf("failed to execute transactions: %w", err)
			}

			return nil
		})
	}

	return err
}

func getStatements(filename string) ([]spanner.Statement, error) {
	var statements []spanner.Statement
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s, %w", filename, err)
	}
	defer file.Close()

	if err != nil {
		return nil, err
	}
	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)

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
