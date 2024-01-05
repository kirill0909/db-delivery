package repository

const (
	queryUserActivation = `
	UPDATE bot.users SET active = TRUE, tg_chat_id = $1, tg_name = $3 
	WHERE tg_uuid = $2 AND active IS FALSE
	`
)
