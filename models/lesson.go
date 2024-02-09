package models

import (
	"github.com/google/uuid"
	"time"
)

type Lesson struct {
	ID        uuid.UUID  `json:"id" gorm:"primary_key"`
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
	TimeSpent int        `json:"time_spent"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Question struct {
	ID            uuid.UUID `json:"id" gorm:"primary_key"`
	QuestionText  string    `json:"question_text"`
	QuestionImage string    `json:"question_image"`
	Answers       string    `json:"answers"`
	Type          string    `json:"type"`
	RightAnswer   string    `json:"right_answer"`
	Topic         string    `json:"topic"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LessonID      uuid.UUID `json:"lesson_id"`
}
type Explanation struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key"`
	Explanation string    `json:"explanation"`
	QuestionID  uuid.UUID `json:"question_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Question Question `json:"question"`
}

type LessonWithProperTitle struct {
	ID             uuid.UUID      `json:"id" gorm:"primary_key"`
	Title          string         `json:"title"`
	ProperTitle    string         `json:"proper_title"`
	LessonAnalytic LessonAnalytic `json:"lesson_analytic"`
	Questions      []Question     `json:"questions"`
	CreatedAt      time.Time      `json:"created_at"`
}

var Topics = map[string]string{
	//history
	"pochatok_ukrayinskoyi_revolyuciyi":                        "Початок Української Революції",
	"ukrayina_v_roki_pershoyi_svitovoyi_vijni":                 "Україна в роки Першої Світової Війни",
	"dyrektorija_unr":                                          "Українська революція. Директорія УНР",
	"period_ghetmanatu":                                        "Українська революція. Період Гетьманату",
	"zahidnoukrayinski_zemli_v_mizhvoyennij_period":            "Західноукраїнські землі в міжвоєнний період",
	"utverdzhennya_bilshovickogo_rezhimu_v_ukrayini":           "Утвердження більшовицького тоталітарного режиму в Україні",
	"vstanovlennya_komunistichnogo_rezhimu_v_ukrayini":         "Встановлення комуністичного тоталітарного режиму в Україні",
	"pochatok_drughoyi_svitovoyi_vijny_1939_1941":              "Початок Другої світової війни (1939–1941)",
	"rukh_oporu_1941_1943":                                     "Друга світова війна: Рух Опору (1941–1943)",
	"voyenni_diyi_na_terenakh_ukrayiny_1943_1945":              "Друга світова війна: воєнні дії на теренах України (1943–1945)",
	"ukrayina_v_pershi_povoyenni_roki":                         "Україна в перші повоєнні роки",
	"ukrayina_v_umovah_destalinizaciyi":                        "Україна в умовах десталінізації",
	"ukrayina_v_period_zagostrennya_krizi_radyanskoyi_sistemi": "Україна в період загострення кризи радянської системи",
	"vidnovlennya_nezalezhnosti_ukrayini":                      "Відновлення незалежності України",
	"stanovlennya_ukrayini_yak_nezalezhnoyi_derzhavi":          "Становлення України як незалежної держави",
	"tvorennya_novoyi_ukrayini":                                "Творення нової України",
	//ukrainian
	"sklad_nagolos":                       "Наголоси",
	"cherguvannya_zvukiv":                 "Співвідношення звуків і букв",
	"pravopis_eio":                        "Правопис літер, що позначають ненаголошені голосні е, и, о",
	"prefiksy":                            "Правопис префіксів ",
	"sproshennya_prigolosnih":             "Спрощення в групах приголосних",
	"zmini_prigolosnih_pri_tvorenni_sliv": "Зміни приголосних при творенні слів",
	"apostrof":                            "Апостроф",
	//biology
	"biologhija_jak_nauka_pro_zhyve":                 "Біологія як наука про живе",
	"elementnyj_sklad_klityny_neorghanichni_spoluky": "Елементний склад клітини. Неорганічні сполуки",
	"orghanichni_spoluky_vughlevody_lipidy":          "Органічні сполуки. Вуглеводи, ліпіди",
	"orghanichni_spoluky_bilky":                      "Органічні сполуки. Білки",
	"orghanichni_spoluky_nukleyinovi_kysloty_atf":    "Органічні сполуки. Нуклеїнові кислоти. АТФ",
}
