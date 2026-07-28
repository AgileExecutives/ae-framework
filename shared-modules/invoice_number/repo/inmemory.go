package repo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AgileExecutives/ae-framework/shared-modules/invoice_number/entities"
)

type InMemoryInvoiceRepo struct {
	mu      sync.Mutex
	nextID  uint
	records map[string]*entities.InvoiceNumber // key: tenant_org_year_month
	logs    []*entities.InvoiceNumberLog
}

func NewInMemoryInvoiceRepo() *InMemoryInvoiceRepo {
	return &InMemoryInvoiceRepo{nextID: 1, records: make(map[string]*entities.InvoiceNumber), logs: make([]*entities.InvoiceNumberLog, 0)}
}

func keyFor(tid, oid uint, year, month int) string {
	return fmt.Sprintf("%d_%d_%d_%d", tid, oid, year, month)
}

func (r *InMemoryInvoiceRepo) Find(ctx context.Context, tenantID, organizationID uint, year, month int) (*entities.InvoiceNumber, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := keyFor(tenantID, organizationID, year, month)
	if rec, ok := r.records[k]; ok {
		copy := *rec
		return &copy, nil
	}
	return nil, nil
}

func (r *InMemoryInvoiceRepo) Create(ctx context.Context, rec *entities.InvoiceNumber) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if rec == nil {
		return fmt.Errorf("nil record")
	}
	if rec.ID == 0 {
		rec.ID = r.nextID
		r.nextID++
	}
	k := keyFor(rec.TenantID, rec.OrganizationID, rec.Year, rec.Month)
	copy := *rec
	r.records[k] = &copy
	return nil
}

func (r *InMemoryInvoiceRepo) Update(ctx context.Context, rec *entities.InvoiceNumber) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if rec == nil {
		return fmt.Errorf("nil record")
	}
	k := keyFor(rec.TenantID, rec.OrganizationID, rec.Year, rec.Month)
	if _, ok := r.records[k]; !ok {
		return fmt.Errorf("not found")
	}
	copy := *rec
	r.records[k] = &copy
	return nil
}

func (r *InMemoryInvoiceRepo) CreateLog(ctx context.Context, l *entities.InvoiceNumberLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if l == nil {
		return fmt.Errorf("nil log")
	}
	if l.ID == 0 {
		l.ID = uint(len(r.logs) + 1)
	}
	if l.GeneratedAt.IsZero() {
		l.GeneratedAt = time.Now()
	}
	copy := *l
	r.logs = append(r.logs, &copy)
	return nil
}

func (r *InMemoryInvoiceRepo) GetLogs(ctx context.Context, tenantID, organizationID uint, year, month int, page, pageSize int) ([]entities.InvoiceNumberLog, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entities.InvoiceNumberLog
	for _, l := range r.logs {
		if l.TenantID != tenantID {
			continue
		}
		if organizationID > 0 && l.OrganizationID != organizationID {
			continue
		}
		if year > 0 && l.Year != year {
			continue
		}
		if month > 0 && l.Month != month {
			continue
		}
		out = append(out, *l)
	}
	total := int64(len(out))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (page - 1) * pageSize
	if start >= len(out) {
		return []entities.InvoiceNumberLog{}, total, nil
	}
	end := start + pageSize
	if end > len(out) {
		end = len(out)
	}
	return out[start:end], total, nil
}
