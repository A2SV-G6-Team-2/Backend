-- Users table
CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY,
    name TEXT,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    budgeting_style TEXT NOT NULL,
    default_currency TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    user_id UUID,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Expenses table
CREATE TABLE IF NOT EXISTS expenses (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    amount DECIMAL NOT NULL,
    category_id UUID,
    is_recurring BOOLEAN DEFAULT FALSE,
    recurrence_type TEXT,
    next_due_date DATE,
    reminder_enabled BOOLEAN DEFAULT FALSE,
    reminder_sent_at TIMESTAMP,
    note TEXT,
    expense_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Debts table
CREATE TABLE IF NOT EXISTS debts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type TEXT NOT NULL,
    peer_name TEXT NOT NULL,
    amount DECIMAL NOT NULL,
    due_date DATE NOT NULL,
    reminder_enabled BOOLEAN DEFAULT FALSE,
    remind_at TIMESTAMP,
    sent_at TIMESTAMP,
    status TEXT DEFAULT 'pending',
    note TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Reminders table
CREATE TABLE IF NOT EXISTS reminders (
    id UUID PRIMARY KEY,
    debt_id UUID NOT NULL,
    remind_at TIMESTAMP NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (debt_id) REFERENCES debts(id)
);
