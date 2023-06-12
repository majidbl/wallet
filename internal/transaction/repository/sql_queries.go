package repository

const (
	createTransactionQuery = `
		INSERT INTO transactions ( wallet_id, amount, type, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, wallet_id, amount, type, created_at
	`

	getByIDQuery = `
        SELECT id, wallet_id, amount, type, created_at
		       FROM transactions  WHERE id = $1`
	getByWalletIDQuery = `
        SELECT id, wallet_id, amount, type, created_at
		       FROM transactions  WHERE wallet_id = $1`
)
