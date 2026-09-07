package entity

import "time"

// Categoria, Severidade e Status são enums fechados. Os valores espelham
// exatamente frontend/src/lib/constants.ts — qualquer divergência quebra
// o contrato da API.
type (
	Categoria  string
	Severidade string
	Status     string
)

const (
	CategoriaPhishing       Categoria = "phishing"
	CategoriaMalware        Categoria = "malware"
	CategoriaAcessoIndevido Categoria = "acesso_indevido"
	CategoriaVazamentoDados Categoria = "vazamento_dados"
	CategoriaOutro          Categoria = "outro"

	SeveridadeBaixa   Severidade = "baixa"
	SeveridadeMedia   Severidade = "media"
	SeveridadeAlta    Severidade = "alta"
	SeveridadeCritica Severidade = "critica"

	StatusAberto    Status = "aberto"
	StatusEmAnalise Status = "em_analise"
	StatusResolvido Status = "resolvido"
)

// CategoriasValidas etc. são as listas canônicas usadas pela validação.
var (
	CategoriasValidas  = []Categoria{CategoriaPhishing, CategoriaMalware, CategoriaAcessoIndevido, CategoriaVazamentoDados, CategoriaOutro}
	SeveridadesValidas = []Severidade{SeveridadeBaixa, SeveridadeMedia, SeveridadeAlta, SeveridadeCritica}
	StatusValidos      = []Status{StatusAberto, StatusEmAnalise, StatusResolvido}
)

// Valida indica se o valor pertence ao conjunto fechado.
func (c Categoria) Valida() bool  { return contains(CategoriasValidas, c) }
func (s Severidade) Valida() bool { return contains(SeveridadesValidas, s) }
func (s Status) Valida() bool     { return contains(StatusValidos, s) }

func contains[T comparable](set []T, v T) bool {
	for _, x := range set {
		if x == v {
			return true
		}
	}
	return false
}

// Incident é o registro de incidente de segurança. OwnerID amarra o
// recurso ao usuário dono — todo acesso é filtrado por ele (OWASP A01).
type Incident struct {
	ID         string     `json:"id"`
	Titulo     string     `json:"titulo"`
	Descricao  string     `json:"descricao"`
	Categoria  Categoria  `json:"categoria"`
	Severidade Severidade `json:"severidade"`
	Status     Status     `json:"status"`
	OwnerID    string     `json:"owner_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
