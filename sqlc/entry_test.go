package db

import (
	"context"
	"reflect"
	"testing"

	"github.com/0firE1dar/simple-bank/db/util"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	params := CreateEntryParams{
		AccountID: account.ID,
		Amount:    int64(util.RandomMoney()),
	}
	entry, err := testQueries.CreateEntry(context.Background(), params)
	assertUnexpectedError(t, err)

	if entry.AccountID != account.ID {
		t.Errorf("got %v, expected %v", entry.AccountID, account.ID)
	}
	if entry.Amount != params.Amount {
		t.Errorf("got %v, expected %v", entry.Amount, params.Amount)
	}
	if entry.CreatedAt.IsZero() || entry.ID == 0 {
		t.Errorf("entry values should not be empty, %v, %v", entry.CreatedAt, entry.ID)
	}

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)

}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account)

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	assertUnexpectedError(t, err)

	if !reflect.DeepEqual(entry, entry2) {
		t.Errorf("expected equal entries, got %v, want %v", entry2, entry)
	}
}

func TestGetEntries(t *testing.T) {
	account := createRandomAccount(t)
	var entires []Entry

	for i := 0; i < 5; i++ {
		entry := createRandomEntry(t, account)
		entires = append(entires, entry)
	}

	entries2, err := testQueries.GetEntries(context.Background(), GetEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    0,
	})
	assertUnexpectedError(t, err)

	for i := 0; i < 5; i++ {
		if !reflect.DeepEqual(entires[i], entries2[i]) {
			t.Errorf("got %v, expected %v", entries2[i], entires[i])
		}
	}
}
