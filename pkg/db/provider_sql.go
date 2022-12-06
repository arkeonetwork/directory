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
		select 
			id,
			created,
			updated,
			pubkey,
			chain,
			coalesce(bond,0) as bond,
			coalesce(metadata_uri,'') as metadata_uri,
			coalesce(metadata_nonce,0) as metadata_nonce,
			coalesce(status,'Offline') as status,
			coalesce(min_contract_duration,-1) as min_contract_duration,
			coalesce(max_contract_duration,-1) as max_contract_duration,
			coalesce(subscription_rate,-1) as subscription_rate,
			coalesce(paygo_rate,-1) as paygo_rate
		from providers p
		where p.pubkey = $1
		  and p.chain = $2
	`
	sqlInsertBondProviderEvent = `
		insert into provider_bond_events(provider_id,height,txid,bond_rel,bond_abs) values ($1,$2,$3,$4,$5) returning id, created, updated
	`
	sqlInsertModProviderEvent = `
		insert into provider_mod_events(provider_id,height,txid,metadata_uri,metadata_nonce,status,min_contract_duration,max_contract_duration,subscription_rate,paygo_rate)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning id, created, updated
	`
	sqlUpsertProviderMetadata = `
		insert into provider_metadata(provider_id,version,moniker,website,description,location,port,proxy_host,source_chain,event_stream_host,claim_store_location,
			free_rate_limit,free_rate_limit_duration,subscribe_rate_limit,subscribe_rate_limit_duration,paygo_rate_limit,paygo_rate_limit_duration)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		on conflict on constraint prov_version_uniq
		do nothing
		returning id, created, updated
	`
)
