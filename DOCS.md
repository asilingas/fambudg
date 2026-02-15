# Fambudg — User Guide

Fambudg is a family budget tracking app. Track income, expenses, budgets, saving goals, and bills across your whole family with role-based access for parents and children.

## Roles

There are three user roles:

- **Admin** — Full access. Manages users, budgets, saving goals, bill reminders, allowances, and sees all family data.
- **Member** — Can create transactions and accounts, view budgets and reports (own data only), pay bills, import/export CSV.
- **Child** — Can view and create own transactions and accounts only. Sees own allowance with spending limit.

The first user to register automatically becomes the admin. Additional users are created by the admin.

## Getting Started

### Register (first user)

Create the admin account by registering with email, password, and name. This is the only self-registration — all other users are created by the admin.

### Login

Log in with email and password. You receive a JWT token used for all subsequent requests.

### Adding Family Members

As admin, go to User Management to create accounts for your spouse (role: member) and children (role: child). Each family member logs in with their own credentials.

## Accounts

Accounts represent where your money lives — bank accounts, wallets, cash, credit cards, etc.

- **Create** an account with a name, type (checking, savings, credit, cash), currency, and starting balance.
- **View** your accounts with current balances. Admin sees all family accounts.
- **Edit** the name, type, or currency of an account.
- **Delete** an account (removes the account record).

Balances update automatically when transactions are created, edited, or deleted.

## Transactions

Transactions are the core of the app — every income and expense you track.

### Creating a Transaction

- **Amount** — in cents. Negative = expense, positive = income. For example, -1999 means a €19.99 expense.
- **Type** — `income`, `expense`, or `transfer`.
- **Account** — which account this transaction belongs to.
- **Category** — what category (groceries, salary, utilities, etc.).
- **Date** — when the transaction occurred.
- **Description** — a note describing the transaction.
- **Tags** — optional labels for extra organization (e.g., "vacation", "birthday").
- **Shared** — mark as a shared family expense (`isShared: true`) or personal.

### Filtering Transactions

Filter your transaction list by:
- Account
- Category
- Type (income/expense/transfer)
- Date range (startDate, endDate)
- Shared or personal (`isShared=true` or `isShared=false`)

### Who Sees What

- **Admin** sees all family transactions.
- **Member** and **Child** see only their own transactions.
- Ownership is enforced — you cannot view, edit, or delete another user's transactions (unless you're admin).

## Categories

Categories organize your transactions (e.g., Groceries, Salary, Rent, Entertainment).

- **Everyone** can view categories.
- **Admin and Member** can create new categories.
- **Only Admin** can edit or delete categories.

## Budgets

Set monthly spending limits per category to stay on track.

- **Create** a budget with a category, month, year, and limit amount (in cents).
- **Summary** shows each budget with the actual amount spent vs. the limit, so you can see if you're over or under budget.
- **Admin** can create, edit, and delete budgets.
- **Member** can view budgets and the summary (read-only).
- **Child** has no access to budgets.

## Reports

### Dashboard

A quick overview of your finances:
- All account balances
- Current month summary (total income, total expenses, net)
- Last 10 recent transactions

Admin sees family-wide data. Member/child sees only their own.

### Monthly Report

Income vs. expense breakdown for a specific month and year. Shows total income, total expenses, and net (income minus expenses).

### Spending by Category

See how much you spent in each category for a given month, with percentage breakdowns. Useful for identifying where most money goes.

### Trends

Month-over-month line chart data showing income, expenses, and net over the last N months (default: 6). Helps you spot patterns — are expenses growing? Is income stable?

### Family Spending Comparison (Admin Only)

See each family member's total income, expenses, and net side by side for a given month. Useful for understanding who's spending what.

## Search

Find transactions across your history with flexible filters:
- **Description** — partial text match (e.g., "groceries")
- **Date range** — startDate and endDate
- **Amount range** — minAmount and maxAmount (in cents)
- **Category** — filter by category ID
- **Account** — filter by account ID
- **Tags** — filter by one or more tags

Admin searches across all family transactions. Others search only their own.

## Saving Goals

Track progress toward saving for something specific (vacation fund, new laptop, emergency fund).

- **Create** a goal with a name, target amount, and optional deadline.
- **Contribute** to a goal by adding an amount. The current amount increases toward the target.
- When the current amount reaches or exceeds the target, the goal is automatically marked as completed.
- **Admin** can create, edit, and contribute to goals.
- **Member** can view goals (read-only).
- **Child** has no access to saving goals.

## Bill Reminders

Never miss a bill payment. Track recurring bills and mark them as paid.

- **Create** a bill reminder with a name, amount, due date, frequency (monthly, weekly, etc.), and optionally link it to a category and account.
- **Upcoming** shows bills due soon, sorted by due date.
- **Mark as Paid** — paying a bill automatically creates a transaction in the linked account with the bill amount, and advances the due date to the next occurrence.
- **Admin** can create, edit, and delete bill reminders.
- **Member** can view reminders and mark them as paid.
- **Child** has no access to bill reminders.

## Transfers

Move money between your own accounts (e.g., from checking to savings).

- Select a source account, destination account, and amount.
- A transfer creates two balance adjustments: the source account decreases, the destination increases.
- You cannot transfer to the same account.

Available to admin and member roles.

## Recurring Transactions

Set up transactions that repeat automatically.

- When creating a transaction, mark it as recurring and set a rule (daily, weekly, monthly with specific day, or yearly).
- **Generate** recurring transactions up to a specified date. The system creates copies of the recurring template for each occurrence.
- Duplicate generation is prevented — if you run generation again, it picks up from where it left off.

Available to admin and member roles.

## CSV Import & Export

### Export

Download all your transactions as a CSV file. The file includes: date, amount, type, description, category ID, account ID, shared flag, and tags.

### Import

Upload a CSV file to bulk-create transactions. The CSV should have columns matching the transaction fields. Useful for migrating data from another app or importing bank statements.

Available to admin and member roles.

## Allowances (Children)

Admin can set a monthly allowance for each child user.

- **Set** an allowance with an amount (in cents) and period start date.
- **Spending** is calculated automatically from the child's expense transactions during the current period — no manual tracking needed.
- **Remaining** shows how much of the allowance is left (allowance amount minus spent).
- **Child** can view their own allowance with spent and remaining amounts.
- **Admin** sees all allowances across all children.

## User Management (Admin Only)

Manage your family's accounts:

- **List** all users with their roles.
- **Create** a new user with email, password, name, and role (admin, member, or child).
- **Edit** a user's name or role.
- **Delete** a user (removes their account).

## Permissions Summary

| Feature | Admin | Member | Child |
|---------|-------|--------|-------|
| Accounts | All family | Own only | Own only |
| Transactions | All family | Own only | Own only |
| Categories | Full CRUD | Read + Create | Read only |
| Budgets | Full CRUD | Read only | No access |
| Reports | All family data | Own data | Own data |
| Family comparison | Yes | No | No |
| Search | All family | Own only | Own only |
| Saving Goals | Full CRUD | Read only | No access |
| Bill Reminders | Full CRUD | Read + Pay | No access |
| Transfers | Yes | Yes | No |
| Recurring | Yes | Yes | No |
| CSV Import/Export | Yes | Yes | No |
| User Management | Yes | No | No |
| Allowances | Manage all | No | View own |
