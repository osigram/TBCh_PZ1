package handlers

import (
	"PZ1/internal/domain"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type BlockchainSetter interface {
	AddBlock(block *domain.Block) error
}

type BlockchainGetSetter interface {
	BlockchainGetter
	BlockchainSetter
}

func Transaction(logger *slog.Logger, getSetter BlockchainGetSetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "internal.server.handlers.Transaction"),
		)

		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			l.Debug("unable to read request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var transaction domain.Transaction
		if err := json.Unmarshal(requestBody, &transaction); err != nil {
			l.Debug("unable to decode request body from json", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		blockchain := getSetter.Blockchain()
		block, err := domain.NewBlock(&blockchain, transaction)
		if err != nil {
			l.Debug("unable to create new block", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = getSetter.AddBlock(block)
		if err != nil {
			l.Debug("unable to set block", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
