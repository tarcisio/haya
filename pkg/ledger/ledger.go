package ledger

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// AccountTypeType represents the type of an account.
//
// It affects how the account is treated in the ledger.
//   - Some accounts increase with a debit and decrease with a credit,
//   - While others increase with a credit and decrease with a debit.
type AccountTypeType string

const (
	AccountTypeAsset     AccountTypeType = "Asset"     // Asset accounts represent the resources owned by the business.
	AccountTypeExpense   AccountTypeType = "Expense"   // Expense accounts represent the costs incurred by the business.
	AccountTypeLiability AccountTypeType = "Liability" // Liability accounts represent the obligations of the business.
	AccountTypeEquity    AccountTypeType = "Equity"    // Equity accounts represent the owner's claim on the assets of the business.
	AccountTypeRevenue   AccountTypeType = "Revenue"   // Revenue accounts represent the income earned by the business.
)

// Account represents a single account in a Ledger.
//
// An account can be a parent account, a child account or both.
//   - if the account is a top-level account, your **ParentID** should be the zero value of uuid.UUID.
//   - A parent account can have multiple child accounts.
//   - A child account can have only one parent account.
type Account struct {
	ID          uuid.UUID
	ParentID    uuid.UUID
	Name        string
	AccountType AccountTypeType
}

// Entry represents a single entry in a Transaction.
//
// It represents the amount of money that is moved from one account to another.
// The amount can be positive or negative.
//   - A positive amount represents an increase in the account balance.
//   - A negative amount represents a decrease in the account balance.
//
// It not specify if it is a debit or a credit for a number of reasons:
//   - The account type determines if the amount is a debit or a credit.
//   - Sistematicaly it is better to use positive and negative numbers to avoid confusion or complexity.
//   - It is easier to calculate the total amount of increases and decreases in a transaction.
type Entry struct {
	Account uuid.UUID
	Amount  int // The amount can be positive or negative.
}

// Transaction represents a single transaction in a Ledger.
type Transaction struct {
	// All the entries in the transaction
	// The sum of all the amounts in the entries should be 0
	// meaning that the transaction is balanced.
	Entries   []Entry
	Timestamp time.Time
}

// NewTransaction creates a new transaction with the given timestamp.
//
// It is easier to create an empty transaction and add entries later because that are multiple ways to add entries.
func NewTransaction(timestamp time.Time) *Transaction {
	t := &Transaction{
		Entries:   make([]Entry, 0),
		Timestamp: timestamp,
	}

	return t
}

// IsBalanced returns true if the transaction is balanced.
//
//   - If the transaction has no entries, it is considered balanced but returns an error.
//   - If the transaction has only one entry, it is considered unbalanced.
//   - A transaction is balanced if the sum of all the amounts in the entries is 0.
func (t *Transaction) IsBalanced() (bool, error) {
	if len(t.Entries) == 0 {
		return true, errors.New("transaction has no entries")
	}

	if len(t.Entries) == 1 {
		return false, errors.New("transaction has only one entry")
	}

	{
		// Check if the transaction is balanced summing all the amounts in the entries
		// if it is not 0, return false and an error.
		var sum int
		for _, entry := range t.Entries {
			sum += entry.Amount
		}
		if sum != 0 {
			return false, errors.New("transaction is unbalanced")
		}
	}

	return true, nil
}

// TotalIncreases returns the total amount of all the increases in the transaction.
// knowing the type of account it is possible to know if the amount is a debit or a credit.
func (t *Transaction) TotalIncreases() int {
	var total int
	for _, entry := range t.Entries {
		if entry.Amount > 0 {
			total += entry.Amount
		}
	}
	return total
}

// TotalDecreases returns the total amount of all the decreases in the transaction.
// knowing the type of account it is possible to know if the amount is a debit or a credit.
func (t *Transaction) TotalDecreases() int {
	var total int
	for _, entry := range t.Entries {
		if entry.Amount < 0 {
			total += entry.Amount
		}
	}
	return total
}

// AddEntry adds an entry to the transaction.
func (t *Transaction) AddEntry(e Entry) {
	t.Entries = append(t.Entries, e)
}

// AddEntries adds multiple entries to the transaction.
func (t *Transaction) AddEntries(entries []Entry) {
	t.Entries = append(t.Entries, entries...)
}

// Ledger represents a collection of transactions.
type Ledger struct {
	storage Storage
}

// NewLedger creates a new ledger with the given [Storage].
func NewLedger(storage Storage) *Ledger {
	return &Ledger{
		storage: storage,
	}
}

// AddTransaction adds a transaction to the ledger.
func (l *Ledger) AddTransaction(t *Transaction) error {
	if balanced, err := t.IsBalanced(); !balanced {
		return err
	}

	return l.storage.SaveTransaction(t)
}

// Storage represents a storage engine for the ledger.
type Storage interface {
	SaveTransaction(t *Transaction) error
}
