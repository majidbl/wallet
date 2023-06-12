package repository

const (
	createWalletQuery = `
	INSERT INTO wallets (name, mobile, balance, avatar, description) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, mobile, balance, avatar, description, created_at, updated_at;`

	getByIDQuery = `SELECT  id, name, mobile, balance, avatar, description, created_at, updated_at 
                                FROM wallets WHERE id = $1`

	getByMobileQuery = `SELECT  id, name, mobile, balance, avatar, description, created_at, updated_at 
                                FROM wallets WHERE mobile = $1`

	updateBalanceQuery = `
		UPDATE wallets
		SET balance = $1
		WHERE id = $2
	`
)
