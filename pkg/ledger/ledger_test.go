package ledger_test

import (
	"haya/pkg/ledger"
	"testing"
	"time"

	"github.com/google/uuid"
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

	if ok, err := transaction.IsBalanced(); !ok || err == nil {
		t.Error("transaction should be balanced but returning error")
	}

	transaction.AddEntry(debit)
	if ok, err := transaction.IsBalanced(); ok || err == nil {
		t.Error("transaction should be unbalanced and returning error")
	}

	transaction.AddEntries([]ledger.Entry{credit})
	if ok, err := transaction.IsBalanced(); !ok || err != nil {
		t.Error("transaction should be balanced and not returning error")
	}

	transaction.AddEntry(plus)
	if ok, err := transaction.IsBalanced(); ok || err == nil {
		t.Error("transaction should be unbalanced and returning error if the sum is not 0")
	}
}
