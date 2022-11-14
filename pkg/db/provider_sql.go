package db

var (
	sqlInsertProvider = `
		insert into providers(pubkey,chain,bond) values ($1,$2,$3) returning id, created, updated
	`
	sqlFindProvider = `
		select id, created, updated, pubkey, chain, bond
		from providers p
		where p.pubkey = $1
		  and p.chain = $2
	`
)
