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

	ErrRelationNotFound       = errors.New("relação usuário-categoria não encontrada")
	ErrRelationExists         = errors.New("relação já existe")
	ErrCreateRelation         = errors.New("erro ao criar relação")
	ErrGetRelationsByUser     = errors.New("erro ao buscar relações por usuário")
	ErrGetRelationsByCategory = errors.New("erro ao buscar relações por categoria")
	ErrScanRelation           = errors.New("erro ao ler relação")
	ErrIterateRelations       = errors.New("erro após ler relações")
	ErrDeleteRelation         = errors.New("erro ao deletar relação")
	ErrDeleteAllUserRelations = errors.New("erro ao deletar todas as relações do usuário")
	ErrInvalidForeignKey      = errors.New("usuário ou categoria inválido")
	ErrCheckRelationExists    = errors.New("erro ao verificar se a relação entre usuário e categoria existe")

	ErrGetUser     = errors.New("erro ao buscar usuário")
	ErrGetAllUsers = errors.New("erro ao buscar usuários")

	ErrFetchCategories     = errors.New("erro ao buscar categorias")
	ErrFetchCategory       = errors.New("erro ao buscar categoria")
	ErrInvalidCategoryName = errors.New("o nome da categoria é obrigatório")
	ErrInvalidCategory     = errors.New("categoria: objeto inválido")
	ErrCheckBeforeUpdate   = errors.New("erro ao verificar dados antes da atualização")

	ErrInvalidUserID              = errors.New("ID do usuário inválido")
	ErrInvalidCategoryID          = errors.New("ID da categoria inválido")
	ErrCheckExistingRelation      = errors.New("erro ao verificar relação existente")
	ErrFetchUserRelations         = errors.New("erro ao buscar relações do usuário")
	ErrFetchCategoryRelations     = errors.New("erro ao buscar relações da categoria")
	ErrCheckUserCategoryRelations = errors.New("erro ao verificar relações do usuário")
)
