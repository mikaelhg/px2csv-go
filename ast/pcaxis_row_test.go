package ast_test

import (
	"github.com/mikaelhg/gpcaxis/ast"

	"testing"

	"github.com/alecthomas/participle/v2"
)

var (
	rowParser = participle.MustBuild[ast.PxRow](
		participle.Lexer(ast.PxLexer),
		participle.Unquote("String"),
		participle.Elide("whitespace", "EOL"),
	)
)

func TestPxRowWithLang(t *testing.T) {
	text := `SUBJECT-AREA[sv]="Besiktningar av personbilar";`
	sv := "sv"
	expected := ast.PxRow{
		Keyword: ast.PxKeyword{
			Keyword:  "SUBJECT-AREA",
			Language: &sv,
		},
		Value: ast.PxValue{
			List: &[]ast.PxStringVal{
				{
					Strings: []string{"Besiktningar av personbilar"},
				},
			},
		},
	}
	parseRow(t, expected, text)
}

func TestPxRow(t *testing.T) {
	text := `SUBJECT-AREA="Besiktningar av personbilar";`
	expected := ast.PxRow{
		Keyword: ast.PxKeyword{
			Keyword: "SUBJECT-AREA",
		},
		Value: ast.PxValue{
			List: &[]ast.PxStringVal{
				{
					Strings: []string{"Besiktningar av personbilar"},
				},
			},
		},
	}
	parseRow(t, expected, text)
}

func TestMultilineNote(t *testing.T) {
	text := `NOTE="<A HREF='https://www.stat.fi/til/vtp/rev.html' TARGET=_blank>Tietojen tarkentuminen</A>#<A HREF='https://www.stat.fi/til/vtp/meta.html' TARGET=_blank>Tilaston kuvaus</A>#<A HREF='https://www.stat.fi/til/vtp/laa.html' TARGET=_blank>Laatuselosteet</A>#"
"<A HREF='https://www.stat.fi/til/vtp/men.html' TARGET=_blank>Menetelmäseloste</A>#<A HREF='https://www.stat.fi/til/vtp/kas.html' TARGET=_blank>Käsitteet ja määritelmät</A>#<A HREF='https://www.stat.fi/til/vtp/uut.html' TARGET=_blank>Muutoksia tässä "
"tilastossa</A>##... tieto on salassapitosäännön alainen#... tieto on salassapitosäännön alainen##Poiketen muista kansantalouden tilinpidon taulukoista, tässä on esitetty (pl. P51K Kiinteän pääoman bruttomuodostus)  volyymisarjat sekä perusvuoden 2010 "
"hintaisina (kantaindeksi; voidaan summata alasarjoista) että viitevuoden 2010 sarjoina (ketjuindeksi; ei voida summata alasarjoista). Viitevuoden 2010 volyymisarja on muodostettu ketjuttamalla edellisen vuoden hintaiset sarjat.";`
	expected := ast.PxRow{
		Keyword: ast.PxKeyword{
			Keyword: "NOTE",
		},
		Value: ast.PxValue{
			List: &[]ast.PxStringVal{
				{
					Strings: []string{
						`<A HREF='https://www.stat.fi/til/vtp/rev.html' TARGET=_blank>Tietojen tarkentuminen</A>#<A HREF='https://www.stat.fi/til/vtp/meta.html' TARGET=_blank>Tilaston kuvaus</A>#<A HREF='https://www.stat.fi/til/vtp/laa.html' TARGET=_blank>Laatuselosteet</A>#`,
						`<A HREF='https://www.stat.fi/til/vtp/men.html' TARGET=_blank>Menetelmäseloste</A>#<A HREF='https://www.stat.fi/til/vtp/kas.html' TARGET=_blank>Käsitteet ja määritelmät</A>#<A HREF='https://www.stat.fi/til/vtp/uut.html' TARGET=_blank>Muutoksia tässä `,
						`tilastossa</A>##... tieto on salassapitosäännön alainen#... tieto on salassapitosäännön alainen##Poiketen muista kansantalouden tilinpidon taulukoista, tässä on esitetty (pl. P51K Kiinteän pääoman bruttomuodostus)  volyymisarjat sekä perusvuoden 2010 `,
						`hintaisina (kantaindeksi; voidaan summata alasarjoista) että viitevuoden 2010 sarjoina (ketjuindeksi; ei voida summata alasarjoista). Viitevuoden 2010 volyymisarja on muodostettu ketjuttamalla edellisen vuoden hintaiset sarjat.`,
					},
				},
			},
		},
	}
	parseRow(t, expected, text)
}
