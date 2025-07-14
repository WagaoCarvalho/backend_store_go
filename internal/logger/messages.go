package logger

const (
	//Create
	LogCreateInit       = "Iniciando criação"
	LogCreateError      = "Erro ao criar"
	LogCreateSuccess    = "Sucesso ao criar"
	LogMethodNotAllowed = "Método nao permitido"

	//Get
	LogGetInit      = "Iniciando busca"
	LogGetError     = "Erro ao buscar"
	LogGetSuccess   = "Sucesso ao buscar"
	LogGetErrorScan = "Erro ao fazer scan"
	LogInvalidID    = "ID inválido"

	//Update
	LogUpdateInit            = "Iniciando Atualização"
	LogUpdateError           = "Erro ao atualizar"
	LogUpdateSuccess         = "Sucesso ao atualizar"
	LogUpdateVersionConflict = "Conflito de versão"
	LogMissingBodyData       = "Dados do usuário são obrigatórios"

	//Delete
	LogDeleteInit    = "Iniciando exclusão"
	LogDeleteError   = "Erro ao deletar"
	LogDeleteSuccess = "Sucesso ao deletar"

	//HasRelation
	LogVerificationInit    = "Iniciando verificação"
	LogAlreadyExists       = "Relação já existe"
	LogVerificationError   = "Erro ao verificar relação"
	LogVerificationSuccess = "Verificação concluída com sucesso"

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
