package db

var (
	sqlUpsertContract = `
		insert into contracts(provider_id,delegate_pubkey,client_pubkey,contract_type,duration,rate,open_cost)
		values ($1,$2,$3,$4,$5,$6,$7)
		on conflict on constraint pubkey_prov_deleg_uniq
		do update set duration = $5, rate = $6, open_cost = $7, updated = now()
		returning id, created, updated
	`
)
