package ledger

import (
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

// AccountBalance represents the balance of an account at a given time.
type AccountBalance struct {
	AccountID   uuid.UUID
	AccountType AccountType
	Balance     int
	Timestamp   time.Time
}
