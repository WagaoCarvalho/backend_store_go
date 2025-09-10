package validators

const (
	MsgRequiredField      = "campo obrigatório"
	MsgMin3               = "mínimo de 3 caracteres"
	MsgMax50              = "máximo de 50 caracteres"
	MsgMax100             = "máximo de 100 caracteres"
	MsgMax255             = "máximo de 255 caracteres"
	MsgMin2               = "mínimo de 2 caracteres"
	MsgInvalidState       = "estado inválido"
	MsgInvalidCountry     = "país não suportado (somente Brasil)"
	MsgInvalidPostalCode  = "formato inválido (ex: 12345678)"
	MsgInvalidAssociation = "exatamente um deve ser informado"

	MsgInvalidFormat = "formato inválido"
	MsgInvalidPhone  = "formato inválido (ex: 1112345678)"
	MsgInvalidCell   = "formato inválido (ex: 11912345678)"
	MsgInvalidType   = "tipo inválido"
	MsgOneIDRequired = "exatamente um deve ser informado"

	MsgInvalidEmail = "email inválido"

	MsgCostNonNegative = "preço de custo não pode ser negativo"
	MsgSaleNonNegative = "preço de venda não pode ser negativo"
	MsgStockNegative   = "estoque não pode ser negativo"
)
