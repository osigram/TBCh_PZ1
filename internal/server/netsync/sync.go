package netsync

import (
	"PZ1/internal/domain"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type BlockchainGetter interface {
	Blockchain() domain.Blockchain
}

type BlockchainSetter interface {
	SetBlockchain(blockchain domain.Blockchain)
}

type BlockchainGetSetter interface {
	BlockchainGetter
	BlockchainSetter
}

func StartSynchronizationRoutine(logger *slog.Logger, getSetter BlockchainGetSetter, urls []string, sleepDuration time.Duration) {
	go func() {
		for {
			_ = SynchronizeBlockchain(logger, getSetter, urls)

			time.Sleep(sleepDuration)
		}

	}()
}

func SynchronizeBlockchain(logger *slog.Logger, getSetter BlockchainGetSetter, urls []string) error {
	var err error
	for _, url := range urls {
		localErr := syncBlockchain(logger, getSetter, url)

		errors.Join(err, localErr)
	}

	return err
}

func syncBlockchain(logger *slog.Logger, getSetter BlockchainGetSetter, url string) error {
	l := logger.With(
		slog.String("op", "internal.server.netsync.syncBlockchain"),
		slog.String("url", url),
	)

	resp, err := http.Get(url + "/blockchain")
	if err != nil {
		l.Info("unable to make request to sync blockchain", slog.String("err", err.Error()))

		return errors.New("unable to make request to sync blockchain")
	}
	defer resp.Body.Close()

	foreignBlockBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Error("unable to read response from body", slog.String("err", err.Error()))

		return errors.New("unable to read response from body")
	}

	var foreignBlockchain domain.Blockchain
	if err := json.Unmarshal(foreignBlockBytes, &foreignBlockchain); err != nil {
		l.Debug("unable to unmarshal response body from json", slog.String("err", err.Error()))

		return errors.New("unable to unmarshal response body from json")
	}

	blockchain := getSetter.Blockchain()

	if len(blockchain) > 0 && len(foreignBlockchain) > 0 &&
		foreignBlockchain[len(foreignBlockchain)-1].ID == blockchain[len(blockchain)-1].ID {
		return nil
	}

	if len(foreignBlockchain) > len(blockchain) {
		for num, block := range foreignBlockchain {
			if err := block.Validate(foreignBlockchain[:num]); err != nil {
				l.Debug("unable to validate foreign blockchain", slog.String("err", err.Error()))

				return err
			}
		}

		getSetter.SetBlockchain(foreignBlockchain)
		return nil
	}

	l.Debug("own blockchain is newer")
	return errors.New("own blockchain is newer")
}
