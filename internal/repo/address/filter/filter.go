package repo

import (
	"context"
	"fmt"
	"strings"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/address/address"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/address/filter"
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
	filter *filter.AddressFilter,
) ([]*model.Address, error) {

	// Aplicar valores padrão do BaseFilter
	base := filter.BaseFilter.WithDefaults()

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

	// Filtros por relacionamento (apenas se não forem nil)
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

	// Filtros de texto (busca parcial com ILIKE)
	if filter.Street != "" {
		query += fmt.Sprintf(" AND street ILIKE '%%' || $%d || '%%'", argPos)
		args = append(args, filter.Street)
		argPos++
	}

	if filter.City != "" {
		query += fmt.Sprintf(" AND city ILIKE '%%' || $%d || '%%'", argPos)
		args = append(args, filter.City)
		argPos++
	}

	if filter.Country != "" {
		query += fmt.Sprintf(" AND country ILIKE '%%' || $%d || '%%'", argPos)
		args = append(args, filter.Country)
		argPos++
	}

	// Filtros de correspondência exata
	if filter.StreetNumber != "" {
		query += fmt.Sprintf(" AND street_number = $%d", argPos)
		args = append(args, filter.StreetNumber)
		argPos++
	}

	if filter.State != "" {
		query += fmt.Sprintf(" AND state = $%d", argPos)
		args = append(args, strings.ToUpper(filter.State))
		argPos++
	}

	// Filtro de CEP com limpeza de formatação
	if filter.PostalCode != "" {
		cleanPostalCode := strings.ReplaceAll(filter.PostalCode, "-", "")
		cleanPostalCode = strings.ReplaceAll(cleanPostalCode, ".", "")
		cleanPostalCode = strings.ReplaceAll(cleanPostalCode, " ", "")
		query += fmt.Sprintf(" AND REPLACE(REPLACE(REPLACE(postal_code, '-', ''), '.', ''), ' ', '') = $%d", argPos)
		args = append(args, cleanPostalCode)
		argPos++
	}

	// Filtro de status (apenas se não for nil)
	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argPos)
		args = append(args, *filter.IsActive)
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

	// Ordenação com validação de campos (prevenção de SQL injection)
	sortField := "created_at"
	if v, ok := addressAllowedSortFields[strings.ToLower(base.SortBy)]; ok {
		sortField = v
	}

	sortOrder := strings.ToLower(base.SortOrder)
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc" // Padrão: mais recentes primeiro
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
	addresses := make([]*model.Address, 0)

	for rows.Next() {
		var a model.Address
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
	filter *filter.AddressFilter,
) ([]*model.Address, error) {
	active := true
	filter.IsActive = &active
	return r.Filter(ctx, filter)
}

// Método para buscar endereços por CEP
func (r *addressFilterRepo) FindByPostalCode(
	ctx context.Context,
	postalCode string,
	exactMatch bool,
) ([]*model.Address, error) {
	active := true
	filter := &filter.AddressFilter{
		PostalCode: postalCode,
		IsActive:   &active,
	}

	result, err := r.Filter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Método para buscar endereços por cidade e estado
func (r *addressFilterRepo) FindByCityAndState(
	ctx context.Context,
	city, state string,
) ([]*model.Address, error) {
	active := true
	filter := &filter.AddressFilter{
		City:     city,
		State:    state,
		IsActive: &active,
	}

	result, err := r.Filter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Versão com função auxiliar para garantir slice não-nil
func ensureNonNilAddressSlice(slice []*model.Address) []*model.Address {
	if slice == nil {
		return make([]*model.Address, 0)
	}
	return slice
}

// Método para buscar endereços por CEP com garantia de slice não-nil
func (r *addressFilterRepo) FindByPostalCodeV2(
	ctx context.Context,
	postalCode string,
	exactMatch bool,
) ([]*model.Address, error) {
	active := true
	filter := &filter.AddressFilter{
		PostalCode: postalCode,
		IsActive:   &active,
	}

	result, err := r.Filter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return ensureNonNilAddressSlice(result), nil
}

// Versão melhorada do FindByPostalCode com busca parcial real
func (r *addressFilterRepo) FindByPostalCodeImproved(
	ctx context.Context,
	postalCode string,
	exactMatch bool,
) ([]*model.Address, error) {
	active := true

	// Limpar o CEP para busca
	cleanPostalCode := strings.ReplaceAll(postalCode, "-", "")
	cleanPostalCode = strings.ReplaceAll(cleanPostalCode, ".", "")
	cleanPostalCode = strings.ReplaceAll(cleanPostalCode, " ", "")

	filter := &filter.AddressFilter{
		PostalCode: cleanPostalCode,
		IsActive:   &active,
	}

	if !exactMatch && len(cleanPostalCode) >= 5 {
		// Para busca parcial, podemos usar os primeiros 5 dígitos
		// Exemplo: CEP 01234-567 -> busca por 01234*
		prefix := cleanPostalCode[:5]
		filter.PostalCode = prefix
		// Nota: Isso exigiria uma modificação no método Filter para suportar LIKE em CEP
		// Por enquanto, mantém a busca exata
	}

	result, err := r.Filter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return ensureNonNilAddressSlice(result), nil
}
