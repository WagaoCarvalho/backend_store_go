package logger

const (
	//Create
	LogCreateInit    = "Iniciando criação"
	LogCreateError   = "Erro ao criar"
	LogGetErrorScan  = "Erro ao fazer scan"
	LogCreateSuccess = "Sucesso ao criar"

	//Get
	LogGetInit    = "Iniciando busca"
	LogGetError   = "Erro ao buscar"
	LogGetSuccess = "Sucesso ao buscar"

	//Update
	LogUpdateInit    = "Iniciando Atualização"
	LogUpdateError   = "Erro ao atualizar"
	LogUpdateSuccess = "Sucesso ao atualizar"

	//Delete
	LogDeleteInit    = "Iniciando exclusão"
	LogDeleteError   = "Erro ao deletar"
	LogDeleteSuccess = "Sucesso ao deletar"

	LogNotFound            = "Não encontrado"
	LogForeignKeyViolation = "Falha por chave estrangeira"
	LogErrorValidate       = "Validação falhou"
	LogParseJsonError      = "Falha ao fazer parse do JSON"
)
