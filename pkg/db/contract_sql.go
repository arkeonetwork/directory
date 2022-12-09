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
	sqlFindContractByPubKeys = `
	select
	c.id,
	c.created,
	c.updated,
	c.provider_id,
	c.delegate_pubkey,
	c.client_pubkey,
	c.contract_type,
	c.duration,
	c.rate,
	c.open_cost
	from providers p join contracts c on p.id = c.provider_id
	where p.chain = $1 and p.pubkey = $2 and c.delegate_pubkey = $3
	`

	sqlUpsertContract = `
		insert into contracts(provider_id,delegate_pubkey,client_pubkey,contract_type,duration,rate,open_cost,height)
		values ($1,$2,$3,$4,$5,$6,$7,$8)
		on conflict on constraint pubkey_prov_dlgt_uniq
		do update set contract_type = $4, duration = $5, rate = $6, open_cost = $7, height = $8, updated = now()
		where contracts.provider_id = $1
		  and contracts.delegate_pubkey = $2
		returning id, created, updated
	`
	sqlUpsertOpenContractEvent = `
	insert into open_contract_events(contract_id,client_pubkey,contract_type,height,txid,duration,rate,open_cost)
	values ($1,$2,$3,$4,$5,$6,$7,$8)
	on conflict on constraint open_contract_events_txid_unq
	do update set updated = now()
	where open_contract_events.txid = $5
	returning id, created, updated
	`
	sqlUpsertCloseContractEvent = `
	insert into close_contract_events(contract_id,client_pubkey,delegate_pubkey,height,txid)
	values ($1,$2,$3,$4,$5)
	on conflict on constraint close_contract_events_txid_key
	do update set updated = now()
	where close_contract_events.txid = $5
	returning id, created, updated
	`

	sqlUpsertContractSettlementEvent = `
	insert into contract_settlement_events(contract_id,txid,client_pubkey,height,nonce,paid,reserve)
	values ($1,$2,$3,$4,$5,$6,$7)
	on conflict on constraint contract_settlement_events_txid_key
	do update set updated = now()
	where contract_settlement_events.txid = $2
	returning id, created, updated
`
)
