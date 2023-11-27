package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	fmt.Printf("[DEBUG] Created DB object: store: %v\n", store)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	fmt.Printf("Balance 1:%v\nBalance 2:%v\n", fromAccount.Balance, toAccount.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// test transfers
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//  check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)

		require.Equal(t, transfer.FromAccountID, fromAccount.ID)
		require.Equal(t, transfer.ToAccountID, toAccount.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NoError(t, err)
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, fromAccount.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		ToEntry := result.ToEntry
		require.NoError(t, err)
		require.NotEmpty(t, ToEntry)
		require.Equal(t, ToEntry.AccountID, toAccount.ID)
		require.Equal(t, ToEntry.Amount, amount)
		require.NotZero(t, ToEntry.ID)
		require.NotZero(t, ToEntry.CreatedAt)

		// check accounts
		retFromAccount := result.FromAccount
		require.NotEmpty(t, retFromAccount)
		require.Equal(t, retFromAccount.ID, fromAccount.ID)

		retToAccount := result.ToAccount
		require.NotEmpty(t, retToAccount)
		require.Equal(t, retToAccount.ID, toAccount.ID)

		// check accounts balance
		amountDiff1 := fromAccount.Balance - retFromAccount.Balance
		amountDiff2 := retToAccount.Balance - toAccount.Balance
		require.Equal(t, amountDiff1, amountDiff2)
		require.True(t, amountDiff1 > 0)
		require.True(t, amountDiff1%amount == 0)

		k := int(amountDiff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	// check the final updated account balances
	updatedFromAccount, err := store.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := store.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance-int64(n)*amount, updatedFromAccount.Balance)
	fmt.Printf("Updated account balance: %v\nAmount: %v", updatedToAccount.Balance, amount)
	require.Equal(t, toAccount.Balance+int64(n)*amount, updatedToAccount.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccount := account1
		toAccount := account2

		if i%2 == 1 {
			toAccount = account1
			fromAccount = account2
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated account balances
	updatedFromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedToAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedFromAccount.Balance)
	require.Equal(t, account2.Balance, updatedToAccount.Balance)
}
