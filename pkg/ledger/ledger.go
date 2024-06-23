package ledger

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// AccountType represents the type of an account.
//
// It affects how the account is treated in the ledger.
//   - Some accounts increase with a debit and decrease with a credit,
//   - While others increase with a credit and decrease with a debit.
type AccountType string

const (
	AccountTypeAsset     AccountType = "Asset"     // Asset accounts represent the resources owned by the business.
	AccountTypeExpense   AccountType = "Expense"   // Expense accounts represent the costs incurred by the business.
	AccountTypeLiability AccountType = "Liability" // Liability accounts represent the obligations of the business.
	AccountTypeEquity    AccountType = "Equity"    // Equity accounts represent the owner's claim on the assets of the business.
	AccountTypeRevenue   AccountType = "Revenue"   // Revenue accounts represent the income earned by the business.
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
	AccountType AccountType
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

// TransactionType represents the type of a transaction.
//   - Regular transactions are the most common transactions.
//   - Closing transactions are used to close the books at the end of an accounting period.
type TransactionType string

const (
	TransactionTypeRegular TransactionType = "Regular"
	TransactionTypeClosing TransactionType = "Closing"
)

// Transaction represents a single transaction in a Ledger.
type Transaction struct {
	// All the entries in the transaction
	// The sum of all the amounts in the entries should be 0
	// meaning that the transaction is balanced.
	Entries         []Entry
	Timestamp       time.Time
	TransactionType TransactionType
}

// AccountBalance represents the balance of an account at a given time.
type AccountBalance struct {
	AccountID   uuid.UUID
	AccountType AccountType
	Balance     int
	Timestamp   time.Time
}

// NewTransaction creates a new regular transaction with the given timestamp.
// It is the same as calling [NewRegularTransaction].
// It is easier to create an empty transaction and add entries later because that are multiple ways to add entries.
func NewTransaction(timestamp time.Time) *Transaction {
	return NewRegularTransaction(timestamp)
}

// NewRegularTransaction creates a new regular transaction with the given timestamp.
// It is the same as calling [NewTransaction].
func NewRegularTransaction(timestamp time.Time) *Transaction {
	return newTransaction(timestamp, TransactionTypeRegular)
}

// NewClosingTransaction creates a new closing transaction with the given timestamp.
func NewClosingTransaction(timestamp time.Time) *Transaction {
	return newTransaction(timestamp, TransactionTypeClosing)
}

// newTransaction creates a new transaction with the given timestamp and type.
// It is an internal function used by the others to create a new transaction.
func newTransaction(timestamp time.Time, t_type TransactionType) *Transaction {
	t := &Transaction{
		Entries:         make([]Entry, 0, 2), // A transaction should have at least two entries.
		Timestamp:       timestamp,
		TransactionType: t_type,
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
func (l *Ledger) AddTransaction(ctx context.Context, t *Transaction) error {
	if balanced, err := t.IsBalanced(); !balanced {
		return err
	}

	return l.storage.SaveTransaction(ctx, t)
}

// GetCurrentAccountBalance returns the current balance of an account at the time of the call.
// it can consider or not closed transactions that may affect the balance using the closed parameter.
func (l *Ledger) GetCurrentAccountBalance(ctx context.Context, accountID uuid.UUID, closed bool) (AccountBalance, error) {
	t := time.Now()
	return l.getAccountBalance(ctx, accountID, t, closed)
}

// GetAccountBalanceAt returns the balance of an account at a given time.
// disconsidering closed transactions that may affect the balance at that time.
func (l *Ledger) GetAccountBalanceAt(ctx context.Context, accountID uuid.UUID, t time.Time) (AccountBalance, error) {
	return l.getAccountBalance(ctx, accountID, t, false)
}

// GetAccountBalanceClosedAt returns the balance of an account at a given time.
// considering closed transactions that may affect the balance at that time.
func (l *Ledger) GetAccountBalanceClosedAt(ctx context.Context, accountID uuid.UUID, t time.Time) (AccountBalance, error) {
	return l.getAccountBalance(ctx, accountID, t, true)
}

// getAccountBalance returns the balance of an account at a given time.
// It is an internal function used by the others to get the balance of an account.
func (l *Ledger) getAccountBalance(ctx context.Context, accountID uuid.UUID, t time.Time, closed bool) (AccountBalance, error) {
	_, _, _, _ = ctx, accountID, t, closed
	return AccountBalance{}, nil
}

// Storage represents a storage engine for the ledger.
type Storage interface {
	SaveTransaction(ctx context.Context, t *Transaction) error
}
