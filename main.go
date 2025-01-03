package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/okteto-community/list-namespace-information/api"
)

func main() {
	token := os.Getenv("OKTETO_TOKEN")
	oktetoURL := os.Getenv("OKTETO_URL")

	logLevel := &slog.LevelVar{} // INFO
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	if token == "" || oktetoURL == "" {
		logger.Error("OKTETO_TOKEN and OKTETO_URL environment variables are required")
		os.Exit(1)
	}

	u, err := url.Parse(oktetoURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Invalid OKTETO_URL %s", err))
		os.Exit(1)
	}

	nsList, err := api.GetDevelopmentNamespaces(u.Host, token, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("There was an error requesting the namespaces: %s", err))
		os.Exit(1)
	}

	currentTime := time.Now()
	filename := fmt.Sprintf("namespaces_%s.csv", currentTime.Format("2006-01-02_15-04-05"))

	file, err := os.Create(filename)
	if err != nil {
		logger.Error(fmt.Sprintf("Error creating CSV file: %s", err))
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	headers := []string{"UUID", "CreationDate", "LastUpdatedDate", "Name", "Persistent", "Personal", "Status", "Type", "DaysSinceLastUpdatedDate"}
	if err := writer.Write(headers); err != nil {
		logger.Error(fmt.Sprintf("Error writing headers to CSV: %s", err))
		os.Exit(1)
	}

	// Initialize counters
	totalNamespaces := 0
	personalNamespaces := 0
	nonPersonalNamespaces := 0
	oldNamespaces := 0

	// Write data rows
	for _, ns := range nsList {
		daysSinceLastUpdate := currentTime.Sub(ns.LastUpdated).Hours() / 24
		row := []string{
			ns.Uuid,
			ns.CreationDate.Format("2006-01-02"),
			ns.LastUpdated.Format("2006-01-02"),
			ns.Name,
			strconv.FormatBool(ns.Persistent),
			strconv.FormatBool(ns.Personal),
			string(ns.Status),
			string(ns.Type),
			fmt.Sprintf("%d", int(daysSinceLastUpdate)),
		}

		if err := writer.Write(row); err != nil {
			logger.Error(fmt.Sprintf("Error writing row to CSV: %s", err))
			os.Exit(1)
		}

		// Update counters
		totalNamespaces++
		if ns.Personal {
			personalNamespaces++
		} else {
			nonPersonalNamespaces++
		}
		if daysSinceLastUpdate > 30 {
			oldNamespaces++
		}
	}

	logger.Info(fmt.Sprintf("Summary:"))
	logger.Info(fmt.Sprintf("Total namespaces: %d", totalNamespaces))
	logger.Info(fmt.Sprintf("Personal namespaces: %d", personalNamespaces))
	logger.Info(fmt.Sprintf("Non-personal namespaces: %d", nonPersonalNamespaces))
	logger.Info(fmt.Sprintf("Namespaces not updated in >30 days: %d", oldNamespaces))
	logger.Info(fmt.Sprintf("CSV file written to %s", filename))

}
