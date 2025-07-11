package logger

const (
	//Create
	LogCreateInit    = "Iniciando criação"
	LogCreateError   = "Erro ao criar"
	LogCreateSuccess = "Sucesso ao criar"

	//Get
	LogGetInit      = "Iniciando busca"
	LogGetError     = "Erro ao buscar"
	LogGetSuccess   = "Sucesso ao buscar"
	LogGetErrorScan = "Erro ao fazer scan"

	//Update
	LogUpdateInit            = "Iniciando Atualização"
	LogUpdateError           = "Erro ao atualizar"
	LogUpdateSuccess         = "Sucesso ao atualizar"
	LogUpdateVersionConflict = "Conflito de versão"

	//Delete
	LogDeleteInit    = "Iniciando exclusão"
	LogDeleteError   = "Erro ao deletar"
	LogDeleteSuccess = "Sucesso ao deletar"

	//Email
	LogEmailInvalid = "email inválido"

	//Password
	LogPasswordInvalid = "erro ao hashear senha"

	//ForeignKey
	LogForeignKeyViolation = "Violação de chave estrangeira"
	LogForeignKeyHasExists = "Relação já existe"

	//HasUserCategoryRelation
	LogCheckInit     = "Iniciando verificação de relação entre usuário e categoria"
	LogCheckNotFound = "Relação entre usuário e categoria não encontrada"
	LogCheckError    = "Erro ao verificar existência de relação entre usuário e categoria"
	LogCheckSuccess  = "Relação entre usuário e categoria encontrada"

	LogValidateError = "Validação falhou"
	LogIterateError  = "Erro ao Iterar"
	LogNotFound      = "Não encontrado"

	LogParseJsonError = "Falha ao fazer parse do JSON"
)
