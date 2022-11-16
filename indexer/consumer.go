package indexer

import (
	"os"
	"os/signal"
	"syscall"

	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (a *IndexerApp) consumeEvents(client *tmclient.HTTP) error {
	blockEvents := subscribe(client, "tm.event = 'NewBlockHeader'")
	bondProviderEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgBondProvider'")
	modProviderEvents := subscribe(client, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgModProvider'")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case evt := <-blockEvents:
			data, ok := evt.Data.(tmtypes.EventDataNewBlockHeader)
			if !ok {
				log.Errorf("event not block header: %T", evt.Data)
				continue
			}
			log.Debugf("received block: %d", data.Header.Height)
		case evt := <-bondProviderEvents:
			converted := convertEvent("provider_bond", evt.Events)
			bondProviderEvent, err := parseBondProviderEvent(converted)
			if err != nil {
				log.Errorf("error parsing bondProviderEvent: %+v", err)
				continue
			}
			if err = a.handleBondProviderEvent(bondProviderEvent); err != nil {
				log.Errorf("error handling provider bond event: %+v", err)
				continue
			}
		case evt := <-modProviderEvents:
			converted := convertEvent("provider_mod", evt.Events)
			modProviderEvent, err := parseModProviderEvent(converted)
			if err != nil {
				log.Errorf("error parsing modProviderEvent: %+v", err)
				continue
			}
			if err = a.handleModProviderEvent(modProviderEvent); err != nil {
				log.Errorf("error storing provider bond event: %+v", err)
				continue
			}
			log.Infof("providerModEvent: %#v", modProviderEvent)
		case <-quit:
			log.Infof("received os quit signal")
			return nil
		}
	}
}
