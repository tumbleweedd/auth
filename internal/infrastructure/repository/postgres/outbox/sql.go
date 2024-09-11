package outbox

const (
	queryCreateOutbox = `
	INSERT INTO outbox (id, event_type, payload)
	VALUES ($1, $2, $3);
`

	queryDeleteOutbox = `
	DELETE FROM outbox
	WHERE id = ANY($1);
	`

	queryGetOutbox = `
	SELECT id, event_type, payload
	FROM outbox;
	`
)
