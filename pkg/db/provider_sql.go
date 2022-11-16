package db

var (
	sqlInsertProvider = `
		insert into providers(pubkey,chain,bond) values ($1,$2,$3) returning id, created, updated
	`
	sqlUpdateProvider = `
		update providers
		set bond = $3,
		    metadata_uri = $4,
				metadata_nonce = $5,
				status = $6,
				min_contract_duration = $7,
				max_contract_duration = $8,
				subscription_rate = $9,
				paygo_rate = $10,
				updated = now()
		where pubkey = $1
		  and chain = $2
		returning id, created, updated
	`
	sqlFindProvider = `
		select id,
					 created,
					 updated,
					 pubkey,
					 chain,
					 bond,
					 metadata_uri,
					 metadata_nonce,
					 status,
					 min_contract_duration,
					 max_contract_duration,
					 subscription_rate,
					 paygo_rate
		from providers p
		where p.pubkey = $1
		  and p.chain = $2
	`
	sqlInsertBondProviderEvent = `
		insert into provider_bond_events(provider_id,bond_rel,bond_abs) values ($1,$2,$3) returning id, created, updated
	`
	sqlInsertModProviderEvent = `
		insert into provider_mod_events(provider_id,metadata_uri,metadata_nonce,status,min_contract_duration,max_contract_duration,subscription_rate,paygo_rate) values ($1,$2,$3,$4,$5,$6,$7,$8) returning id, created, updated
	`
)
