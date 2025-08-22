package err

import (
	"errors"
)

var (
	ErrUserIDInvalid       = errors.New("ID de usuário inválido")
	ErrUserNotFound        = errors.New("usuário não encontrado")
	ErrUserVersionConflict = errors.New("conflito de versão: os dados foram modificados por outro processo")
	ErrCreateUser          = errors.New("erro ao criar usuário")
	ErrGetUsers            = errors.New("erro ao buscar usuários")
	ErrScanUserRow         = errors.New("erro ao ler os dados do usuário")
	ErrIterateUserRows     = errors.New("erro ao iterar sobre os resultados")
	ErrFetchUser           = errors.New("erro ao buscar usuário")
	ErrUpdateUser          = errors.New("erro ao atualizar usuário")
	ErrDeleteUser          = errors.New("erro ao deletar usuário")
	ErrDisableUser         = errors.New("erro ao desabilitar usuário")
	ErrEnableUser          = errors.New("erro ao ativar usuário")
	ErrFetchUserVersion    = errors.New("erro ao buscar versão do usuário")

	ErrCategoryNotFound  = errors.New("categoria não encontrada")
	ErrGetCategories     = errors.New("erro ao buscar categorias")
	ErrScanCategory      = errors.New("erro ao ler os dados da categoria")
	ErrIterateCategories = errors.New("erro ao iterar sobre os resultados")
	ErrGetCategoryByID   = errors.New("erro ao buscar categoria por ID")
	ErrCreateCategory    = errors.New("erro ao criar categoria")
	ErrUpdateCategory    = errors.New("erro ao atualizar categoria")
	ErrDeleteCategory    = errors.New("erro ao deletar categoria")
)
