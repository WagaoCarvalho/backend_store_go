package repo

import (
	"context"
	"fmt"
	"strings"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	filterModel "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

var clientCpfAllowedSortFields = map[string]string{
	"id":          "id",
	"name":        "name",
	"email":       "email",
	"cpf":         "cpf",
	"description": "description",
	"status":      "status",
	"version":     "version",
	"created_at":  "created_at",
	"updated_at":  "updated_at",
}

func (r *clientCpfFilterRepo) Filter(
	ctx context.Context,
	filter *filterModel.ClientCpfFilter,
) ([]*model.ClientCpf, error) {

	base := filter.BaseFilter.WithDefaults()

	query := `
		SELECT
			id,
			name,
			email,
			cpf,
			description,
			status,
			version,
			created_at,
			updated_at
		FROM clients_cpf
		WHERE 1=1
	`

	args := []any{}
	argPos := 1

	// Filtros de texto (busca parcial com ILIKE)
	if filter.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE '%%' || $%d || '%%'", argPos)
		args = append(args, filter.Name)
		argPos++
	}

	if filter.Email != "" {
		query += fmt.Sprintf(" AND email ILIKE '%%' || $%d || '%%'", argPos)
		args = append(args, filter.Email)
		argPos++
	}

	// Filtros exatos
	if filter.CPF != "" {
		// Remove formatação do CPF para busca consistente
		cleanCPF := strings.ReplaceAll(filter.CPF, ".", "")
		cleanCPF = strings.ReplaceAll(cleanCPF, "-", "")
		query += fmt.Sprintf(" AND REPLACE(REPLACE(cpf, '.', ''), '-', '') = $%d", argPos)
		args = append(args, cleanCPF)
		argPos++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, *filter.Status)
		argPos++
	}

	if filter.Version != nil {
		query += fmt.Sprintf(" AND version = $%d", argPos)
		args = append(args, *filter.Version)
		argPos++
	}

	// Filtros de data
	if filter.CreatedFrom != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argPos)
		args = append(args, *filter.CreatedFrom)
		argPos++
	}

	if filter.CreatedTo != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argPos)
		args = append(args, *filter.CreatedTo)
		argPos++
	}

	if filter.UpdatedFrom != nil {
		query += fmt.Sprintf(" AND updated_at >= $%d", argPos)
		args = append(args, *filter.UpdatedFrom)
		argPos++
	}

	if filter.UpdatedTo != nil {
		query += fmt.Sprintf(" AND updated_at <= $%d", argPos)
		args = append(args, *filter.UpdatedTo)
		argPos++
	}

	// Ordenação com validação de campos
	sortField := "created_at"
	if v, ok := clientCpfAllowedSortFields[strings.ToLower(base.SortBy)]; ok {
		sortField = v
	}

	sortOrder := strings.ToLower(base.SortOrder)
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	query += fmt.Sprintf(
		" ORDER BY %s %s LIMIT $%d OFFSET $%d",
		sortField,
		sortOrder,
		argPos,
		argPos+1,
	)

	args = append(args, base.Limit, base.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	// Inicializa a slice vazia (garante que nunca retorna nil)
	clients := make([]*model.ClientCpf, 0)

	for rows.Next() {
		var c model.ClientCpf
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&c.CPF,
			&c.Description,
			&c.Status,
			&c.Version,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		clients = append(clients, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	return clients, nil
}

// Método para buscar clientes ativos (status = true)
func (r *clientCpfFilterRepo) FindActive(
	ctx context.Context,
	filter *filterModel.ClientCpfFilter,
) ([]*model.ClientCpf, error) {
	active := true
	filter.Status = &active
	return r.Filter(ctx, filter)
}

// Método para buscar clientes por CPF (busca exata ou parcial)
func (r *clientCpfFilterRepo) FindByCPF(
	ctx context.Context,
	cpf string,
	exactMatch bool,
) ([]*model.ClientCpf, error) {
	filter := &filterModel.ClientCpfFilter{
		CPF: cpf,
	}

	result, err := r.Filter(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Se não for busca exata, mantém o resultado (a query já fez a limpeza)
	// Se for busca exata, a query já garante a igualdade
	return result, nil
}

// Método para buscar clientes por nome (busca parcial)
func (r *clientCpfFilterRepo) FindByName(
	ctx context.Context,
	name string,
) ([]*model.ClientCpf, error) {
	filter := &filterModel.ClientCpfFilter{
		Name: name,
	}
	return r.Filter(ctx, filter)
}

// Função auxiliar para garantir slice não-nil (útil para outros métodos)
func ensureNonNilClientSlice(slice []*model.ClientCpf) []*model.ClientCpf {
	if slice == nil {
		return make([]*model.ClientCpf, 0)
	}
	return slice
}
