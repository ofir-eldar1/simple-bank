package db

import (
	"context"
	"reflect"
	"testing"

	"github.com/0firE1dar/simple-bank/db/util"
)

func createRandomAmountTransfer(t *testing.T, fromAccount, toAccount Account) Transfer {
	amount := util.RandomInt(0, fromAccount.Balance)
	transfer, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        amount,
	})
	assertUnexpectedError(t, err)

	if transfer.Amount != amount {
		t.Errorf("expected to trasnfer amount: %v, not %v", amount, transfer.Amount)
	}
	if transfer.CreatedAt.IsZero() || transfer.ID == 0 {
		t.Errorf("expected non zero values, created_at: %v, id: %v", transfer.CreatedAt, transfer.ID)
	}
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	_ = createRandomAmountTransfer(t, fromAccount, toAccount)

}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	transfer := createRandomAmountTransfer(t, fromAccount, toAccount)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	assertUnexpectedError(t, err)

	if !reflect.DeepEqual(transfer, transfer2) {
		t.Errorf("got %v, want %v", transfer2, transfer)
	}
}

func TestGetTransfers(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	var transfers []Transfer

	for i := 0; i < 5; i++ {
		transfer := createRandomAmountTransfer(t, fromAccount, toAccount)
		transfers = append(transfers, transfer)
	}

	transfer2, err := testQueries.GetTransfers(context.Background(), GetTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Limit:         5,
		Offset:        0,
	})
	assertUnexpectedError(t, err)

	for i := 0; i < 5; i++ {
		if transfers[i] != transfer2[i] {
			t.Errorf("got %v, expected %v", transfer2[i], transfers[i])
		}
	}
}
