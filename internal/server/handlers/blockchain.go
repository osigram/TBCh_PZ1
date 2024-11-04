package handlers

import (
	"PZ1/internal/domain"
	"encoding/json"
	"log/slog"
	"net/http"
)

type BlockchainGetter interface {
	Blockchain() domain.Blockchain
}

func BlockchainHandler(logger *slog.Logger, getter BlockchainGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "internal.server.handlers.Blockchain"),
		)

		l.Debug("getting blockchain")
		block := getter.Blockchain()

		blockBytes, err := json.Marshal(block)
		if err != nil {
			l.Error("unable to marshal blockchain", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(200)
		w.Write(blockBytes)
	}
}
