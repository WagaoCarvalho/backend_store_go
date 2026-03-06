package repo

import (
	"context"
	"fmt"
	"strings"

	address "github.com/WagaoCarvalho/backend_store_go/internal/model/address/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

var addressAllowedSortFields = map[string]string{
	"id":            "id",
	"user_id":       "user_id",
	"client_cpf_id": "client_cpf_id",
	"supplier_id":   "supplier_id",
	"street":        "street",
	"street_number": "street_number",
	"city":          "city",
	"state":         "state",
	"country":       "country",
	"postal_code":   "postal_code",
	"is_active":     "is_active",
	"created_at":    "created_at",
	"updated_at":    "updated_at",
}

func (r *addressFilterRepo) Filter(
	ctx context.Context,
	filter *address.Address,
) ([]*address.Address, error) {

	query := `
		SELECT
			id,
			user_id,
			client_cpf_id,
			supplier_id,
			street,
			street_number,
			complement,
			city,
			state,
			country,
			postal_code,
			is_active,
			created_at,
			updated_at
		FROM addresses
		WHERE 1=1
	`

	args := []any{}
	argPos := 1

	// Filtros por relacionamento
	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, *filter.UserID)
		argPos++
	}

	if filter.ClientCpfID != nil {
		query += fmt.Sprintf(" AND client_cpf_id = $%d", argPos)
		args = append(args, *filter.ClientCpfID)
		argPos++
	}

	if filter.SupplierID != nil {
		query += fmt.Sprintf(" AND supplier_id = $%d", argPos)
		args = append(args, *filter.SupplierID)
		argPos++
	}

	// Filtros por endereço
	if filter.City != "" {
		query += fmt.Sprintf(" AND city ILIKE '%%' || $%d || '%%'", argPos)
		args = append(args, filter.City)
		argPos++
	}

	if filter.State != "" {
		query += fmt.Sprintf(" AND state = $%d", argPos)
		args = append(args, strings.ToUpper(filter.State))
		argPos++
	}

	if filter.PostalCode != "" {
		// Remove caracteres não numéricos para busca
		cleanPostalCode := strings.ReplaceAll(filter.PostalCode, "-", "")
		cleanPostalCode = strings.ReplaceAll(cleanPostalCode, ".", "")
		query += fmt.Sprintf(" AND REPLACE(REPLACE(postal_code, '-', ''), '.', '') = $%d", argPos)
		args = append(args, cleanPostalCode)
		argPos++
	}

	if filter.Street != "" {
		query += fmt.Sprintf(" AND street ILIKE '%%' || $%d || '%%'", argPos)
		args = append(args, filter.Street)
		argPos++
	}

	if filter.StreetNumber != "" {
		query += fmt.Sprintf(" AND street_number = $%d", argPos)
		args = append(args, filter.StreetNumber)
		argPos++
	}

	if filter.Country != "" {
		query += fmt.Sprintf(" AND country ILIKE '%%' || $%d || '%%'", argPos)
		args = append(args, filter.Country)
		argPos++
	}

	// Filtro de status
	query += fmt.Sprintf(" AND is_active = $%d", argPos)
	args = append(args, filter.IsActive)
	argPos++

	// Ordenação padrão
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, 100, 0) // Limite padrão e offset zero

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	// Inicializa a slice vazia
	addresses := make([]*address.Address, 0)

	for rows.Next() {
		var a address.Address
		var complement *string // complement pode ser NULL

		if err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.ClientCpfID,
			&a.SupplierID,
			&a.Street,
			&a.StreetNumber,
			&complement,
			&a.City,
			&a.State,
			&a.Country,
			&a.PostalCode,
			&a.IsActive,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}

		// Trata complemento nulo
		if complement != nil {
			a.Complement = *complement
		}

		addresses = append(addresses, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	return addresses, nil
}

// Método para buscar endereços ativos
func (r *addressFilterRepo) FindActive(
	ctx context.Context,
	filter *address.Address,
) ([]*address.Address, error) {
	filter.IsActive = true
	result, err := r.Filter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Método para buscar endereços por CEP
func (r *addressFilterRepo) FindByPostalCode(
	ctx context.Context,
	postalCode string,
	exactMatch bool,
) ([]*address.Address, error) {
	filter := &address.Address{
		PostalCode: postalCode,
		IsActive:   true,
	}

	var result []*address.Address
	var err error

	if !exactMatch {
		// Para busca parcial, você precisaria modificar a lógica
		// Este é um placeholder - ajuste conforme necessário
		result, err = r.Filter(ctx, filter)
	} else {
		result, err = r.Filter(ctx, filter)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Método para buscar endereços por cidade e estado
func (r *addressFilterRepo) FindByCityAndState(
	ctx context.Context,
	city, state string,
) ([]*address.Address, error) {
	filter := &address.Address{
		City:     city,
		State:    state,
		IsActive: true,
	}
	result, err := r.Filter(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Versão mais concisa com função auxiliar
func ensureNonNilSlice(slice []*address.Address) []*address.Address {
	if slice == nil {
		return make([]*address.Address, 0)
	}
	return slice
}

// Versão alternativa usando função auxiliar
func (r *addressFilterRepo) FindByPostalCodeV2(
	ctx context.Context,
	postalCode string,
	exactMatch bool,
) ([]*address.Address, error) {
	filter := &address.Address{
		PostalCode: postalCode,
		IsActive:   true,
	}

	result, err := r.Filter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return ensureNonNilSlice(result), nil
}
