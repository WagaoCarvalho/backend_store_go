package logger

const (
	// Login
	LogLoginInit            = "iniciando login"
	LogNotFound             = "não encontrado ou erro ao buscar"
	LogPasswordInvalid      = "senha inválida"
	LogAccountDisabled      = "conta desativada"
	LogTokenGenerationError = "erro ao gerar token"
	LogLoginSuccess         = "login realizado com sucesso"
	LogEmailInvalid         = "email inválido"

	// Logout
	LogLogoutInit           = "iniciando logout"
	LogLogoutSuccess        = "logout realizado com sucesso"
	LogInvalidSigningMethod = "método de assinatura inválido"
	LogTokenValidationFail  = "falha ao validar token"
	LogTokenInvalid         = "token inválido"
	LogClaimsConversionFail = "não foi possível converter claims"
	LogClaimExpInvalid      = "claim 'exp' ausente ou inválida"
	LogTokenAlreadyExpired  = "token já expirado"
	LogBlacklistAddFail     = "erro ao adicionar token à blacklist"

	// Create
	LogCreateInit       = "iniciando criação"
	LogCreateError      = "erro ao criar"
	LogCreateSuccess    = "sucesso ao criar"
	LogMethodNotAllowed = "método não permitido"

	// Get
	LogGetInit      = "iniciando busca"
	LogGetError     = "erro ao buscar"
	LogGetSuccess   = "sucesso ao buscar"
	LogGetErrorScan = "erro ao fazer scan"
	LogInvalidID    = "id inválido"

	// Update
	LogUpdateInit            = "iniciando atualização"
	LogUpdateError           = "erro ao atualizar"
	LogUpdateSuccess         = "sucesso ao atualizar"
	LogUpdateVersionConflict = "conflito de versão"
	LogMissingBodyData       = "dados do usuário são obrigatórios"

	// Delete
	LogDeleteInit    = "iniciando exclusão"
	LogDeleteError   = "erro ao deletar"
	LogDeleteSuccess = "sucesso ao deletar"

	// Disable
	LogDisableInit    = "iniciando a desativação"
	LogDisableError   = "erro ao desativar"
	LogDisableSuccess = "sucesso ao desativar"

	// Enable
	LogEnableInit    = "iniciando a desativação"
	LogEnableError   = "erro ao desativar"
	LogEnableSuccess = "sucesso ao desativar"

	// HasRelation
	LogVerificationInit    = "iniciando verificação"
	LogAlreadyExists       = "relação já existe"
	LogVerificationError   = "erro ao verificar relação"
	LogVerificationSuccess = "verificação concluída com sucesso"

	// ForeignKey
	LogForeignKeyViolation = "violação de chave estrangeira"
	LogForeignKeyHasExists = "relação já existe"

	// HasUserCategoryRelation
	LogCheckInit     = "iniciando verificação de relação entre usuário e categoria"
	LogCheckNotFound = "relação entre usuário e categoria não encontrada"
	LogCheckError    = "erro ao verificar existência de relação entre usuário e categoria"
	LogCheckSuccess  = "relação entre usuário e categoria encontrada"

	// Transactions
	LogTransactionInitError         = "erro ao iniciar transação"
	LogTransactionNull              = "transação retornada é nil"
	LogRollbackError                = "erro ao fazer rollback"
	LogCommitError                  = "erro ao commitar transação"
	LogRollbackErrorAfterCommitFail = "erro ao fazer rollback após commit falhar"
	LogQueryError                   = "erro ao executar query"
	LogScanError                    = "erro ao escanear resultado"

	// Utilitários
	LogValidateError  = "validação falhou"
	LogIterateError   = "erro ao iterar"
	LogParseJSONError = "falha ao fazer parse do json"

	LogVersionConflict = "Conflito de versão"

	//JWT
	LogAuthTokenMissing            = "token ausente"
	LogAuthTokenInvalidFormat      = "formato de token inválido"
	LogAuthTokenRevoked            = "token revogado"
	LogAuthTokenInvalid            = "token inválido"
	LogAuthTokenInvalidParsed      = "token inválido (parse ok, mas inválido)"
	LogAuthTokenExpired            = "token expirado"
	LogAuthTokenExpiredManualCheck = "token expirado (verificação manual)"
	LogAuthInvalidSignature        = "assinatura inválida"
	LogAuthInvalidSigningMethod    = "método de assinatura inválido"
	LogAuthBlacklistError          = "erro ao consultar blacklist"
	LogAuthExpClaimInvalid         = "campo exp ausente ou inválido"
)
