package db

var (
	sqlFindContract = `
	select
	id,
	created,
	updated,
	provider_id,
	delegate_pubkey,
	client_pubkey,
	contract_type,
	duration,
	rate,
	open_cost
	from contracts c
		where c.provider_id = $1
		  and c.delegate_pubkey = $2
	`
	sqlUpsertContract = `
		insert into contracts(provider_id,delegate_pubkey,client_pubkey,contract_type,duration,rate,open_cost,height)
		values ($1,$2,$3,$4,$5,$6,$7,$8)
		on conflict on constraint pubkey_prov_dlgt_uniq
		do update set duration = $5, rate = $6, open_cost = $7, height = $8, updated = now()
		returning id, created, updated
	`
	sqlInsertOpenContractEvent = `
	insert into open_contract_events(contract_id,client_pubkey,
		contract_type,height,txid,duration,rate,open_cost)
	values ($1,$2,$3,$4,$5,$6,$7,$8) returning id, created, updated
	`
)
