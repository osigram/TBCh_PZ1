package handlers

import (
	"PZ1/internal/domain"
	"encoding/json"
	"log/slog"
	"net/http"
)

func Account(logger *slog.Logger, getter BlockchainGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			slog.String("op", "internal.server.handlers.Account"),
		)

		clientData := r.URL.Query().Get("key")
		if clientData == "" {
			l.Error("unable to get client data")
			w.WriteHeader(http.StatusBadRequest)

			return
		}
		client := &domain.Client{}
		err := client.UnmarshalJSON([]byte(clientData))
		if err != nil {
			l.Error("unable to get client", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		l.Debug("getting blockchain")
		block := getter.Blockchain()

		account, err := domain.Account(block, *client)
		if err != nil {
			l.Error("unable to get amount", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		accBytes, err := json.Marshal(account)
		if err != nil {
			l.Error("unable to marshal account", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(200)
		w.Write(accBytes)
	}
}
