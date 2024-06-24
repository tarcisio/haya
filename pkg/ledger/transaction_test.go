package ledger_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tarcisio/haya/pkg/ledger"
)

func Test_Transactions(t *testing.T) {

	now := time.Now()

	debit := ledger.Entry{
		Account: uuid.New(),
		Amount:  -100,
	}
	credit := ledger.Entry{
		Account: uuid.New(),
		Amount:  100,
	}

	plus := ledger.Entry{
		Account: uuid.New(),
		Amount:  100,
	}

	transaction := ledger.NewTransaction(now)

	// test if transaction is balanced at the beginning
	if ok, err := transaction.IsBalanced(); !ok || err == nil {
		t.Error("transaction should be balanced but returning error")
	}

	// test if transaction is balanced after adding a one single entry
	transaction.AddEntry(debit)
	if ok, err := transaction.IsBalanced(); ok || err == nil {
		t.Error("transaction should be unbalanced and returning error")
	}

	// test if transaction is balanced after adding a second entry
	transaction.AddEntries([]ledger.Entry{credit})
	if ok, err := transaction.IsBalanced(); !ok || err != nil {
		t.Error("transaction should be balanced and not returning error")
	}

	// test if transaction is balanced after adding a third entry that makes the sum different than 0
	transaction.AddEntry(plus)
	if ok, err := transaction.IsBalanced(); ok || err == nil {
		t.Error("transaction should be unbalanced and returning error if the sum is not 0")
	}

	// test NewClosingTransaction
	close_transaction := ledger.NewClosingTransaction(now)
	close_transaction.AddEntry(debit)
	close_transaction.AddEntry(credit)
	if ok, err := close_transaction.IsBalanced(); !ok || err != nil {
		t.Error("transaction should be balanced and not returning error")
	}

	// test if total increases are correct
	if total := close_transaction.TotalIncreases(); total != credit.Amount {
		t.Errorf("total increases should be %d but got %d", credit.Amount, total)
	}

	// test if total decreases are correct
	if total := close_transaction.TotalDecreases(); total != debit.Amount {
		t.Errorf("total decreases should be %d but got %d", debit.Amount, total)
	}
}
