package services

import (
	"context"
	"testing"

	repo "github.com/AgileExecutives/shared-modules/invoice_number/repo"
	"github.com/stretchr/testify/require"
)

func TestInvoiceNumberService_WithInMemoryRepo(t *testing.T) {
	r := repo.NewInMemoryInvoiceRepo()
	svc := NewInvoiceNumberServiceWithRepo(r)

	ctx := context.Background()
	cfg := DefaultInvoiceConfig()
	resp1, err := svc.GenerateInvoiceNumber(ctx, 1, 0, cfg)
	require.NoError(t, err)
	require.Equal(t, 1, resp1.Sequence)

	resp2, err := svc.GenerateInvoiceNumber(ctx, 1, 0, cfg)
	require.NoError(t, err)
	require.Equal(t, 2, resp2.Sequence)

	// ensure logs recorded
	logs, total, err := r.GetLogs(ctx, 1, 0, 0, 0, 1, 10)
	require.NoError(t, err)
	require.Equal(t, int64(2), total)
	require.Len(t, logs, 2)
	require.Equal(t, resp2.InvoiceNumber, logs[1].InvoiceNumber)
}
