package sqlite

var migrate = []string{
	`
		CREATE TABLE IF NOT EXISTS projects(
			id            INTEGER           PRIMARY KEY AUTOINCREMENT,
			name          varchar(50)       DEFAULT NULL,
			error         boolean           NOT NULL DEFAULT false,
			notes					text       				DEFAULT NULL,
			created_at    datetime(6)       NOT NULL
		);
	`,
}

var drop = []string{
	`drop table if exists projects cascade`,
}
