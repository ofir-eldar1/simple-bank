package db

import (
	"context"
	"reflect"
	"testing"

	"github.com/0firE1dar/simple-bank/db/util"
)

func createRandomAccount(t *testing.T) Account {
	params := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)

	assertUnexpectedError(t, err)

	if account.Balance != params.Balance || account.Owner != params.Owner || account.Currency != params.Currency {
		t.Errorf("got %v, want %v", account, params)
	}

	if account.ID == 0 {
		t.Errorf("expected non zero account id, got %v", account.ID)
	}
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}
func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	got, err := testQueries.GetAccount(context.Background(), account.ID)
	assertUnexpectedError(t, err)

	if !reflect.DeepEqual(account, got) {
		t.Errorf("Got %v, expected %v", got, account)
	}
}
func TestGetAccounts(t *testing.T) {
	// type GetAccountsParams struct {
	// 	Limit  int32 `json:"limit"`
	// 	Offset int32 `json:"offset"`
	// }
	numOfAccounts := 10
	for i := 0; i < numOfAccounts; i++ {
		createRandomAccount(t)
	}

	accountsParams := GetAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.GetAccounts(context.Background(), accountsParams)
	assertUnexpectedError(t, err)
	expectedNumOfAccounts := numOfAccounts - int(accountsParams.Limit)
	if len(accounts) != expectedNumOfAccounts {
		t.Errorf("Expected %v accounts, got %v accounts", expectedNumOfAccounts, len(accounts))
	}

	for _, account := range accounts {
		assertNotEmptyAccount(t, account)
	}

}
func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	assertUnexpectedError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assertExpectedError(t, err)
	if account2 != (Account{}) {
		t.Errorf("Got %v, want: %v", account2, (Account{}))
	}

}
func TestUpdateAccount(t *testing.T) {
	// type UpdateAccountParams struct {
	// 	ID      int64 `json:"id"`
	// 	Balance int64 `json:"balance"`
	// }
	account1 := createRandomAccount(t)
	wantedBalance := util.RandomMoney()
	err := testQueries.UpdateAccount(context.Background(), UpdateAccountParams{ID: account1.ID, Balance: wantedBalance})
	assertUnexpectedError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assertUnexpectedError(t, err)
	if account2.Balance == account1.Balance {
		t.Errorf("Got %v, want %v", account2.Balance, account1.Balance)
	}
}

func assertUnexpectedError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Got an error but didn't expect one\nError:%v", err)
	}
}

func assertExpectedError(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected an error but didn't get one")
	}
}

func assertNotEmptyAccount(t testing.TB, account Account) {
	t.Helper()
	if account == (Account{}) {
		t.Errorf("Expected non empty account, got %v", account)
	}
}
