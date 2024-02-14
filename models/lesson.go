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
	"sklad_nagolos":                           "Наголоси",
	"cherguvannya_zvukiv":                     "Співвідношення звуків і букв",
	"pravopis_eio":                            "Правопис літер, що позначають ненаголошені голосні е, и, о",
	"prefiksy":                                "Правопис префіксів ",
	"sproshennya_prigolosnih":                 "Спрощення в групах приголосних",
	"zmini_prigolosnih_pri_tvorenni_sliv":     "Зміни приголосних при творенні слів",
	"apostrof":                                "Апостроф",
	"m_jakyj_znak":                            "Мʼякий знак",
	"spoluchennya_jo":                         "Сполучення йо, ьо",
	"podvoyennya_liter":                       "Подвоєння літер",
	"pravopys_sliv_inshomovnogo_pohodzhennya": "Правопис слів іншомовного походження",
	"velika_litera_ta_lapki_u_vlasnih_nazvah": "Велика буква та лапки у власних назвах",
	"skladni_slova":                           "Складні слова",
	"pravopis_ne_z_riznimi_chastinami_movi":   "Правопис не з різними частинами мови",
	"osnovni_vipadki_cherguvannya_uv_ij":      "Основні випадки чергування у–в, і–й",
	"leksychne_znachennja_slova":              "Лексичне значення слова",
	"synonimy":                                "Синоніми",
	"leksychna_pomylka":                       "Лексична помилка",
	"frazeolohiia":                            "Фразеологія",
	"budova_slova_slovotvir":                  "Будова слова. Словотвір",
	"klasifikaciya_chastin_movi":              "Класифікація частин мови",
	"imennyk_rid_i_chyslo_imennykiv":          "Іменник. Рід і число іменників",
	"vidminjuvannja_imennykiv":                "Відмінювання іменників",
	"napysannja_imen_po_batkovi":              "Написання імен по батькові",
	"prikmetnik":                              "Прикметник",
	"chislivnik":                              "Числівник",
	"zajmennik":                               "Займенник",
	"diyeslovo":                               "Дієслово",
	"diyeprikmetnik":                          "Дієприкметник",
	"diyeprislivnik":                          "Дієприслівник",
	"prislivnik":                              "Прислівник",
	"prijmennik":                              "Прийменник",
	"spoluchnik":                              "Сполучник",
	"chastka":                                 "Частка",
	"slovospoluchennya":                       "Словосполучення",
	"klasifikaciya_rechen":                    "Класифікація речень",
	"golovni_chleni_rechennya":                "Головні члени речення",
	"drugoryadni_chleni_rechennya":            "Другорядні члени речення",
	"odnoskladni_rechennya":                   "Односкладне речення",
	"odnoridni_chleny_rechennja":              "Однорідні члени речення",
	"zvertannja_vstavni_slova_slovospoluchennja_rechennja": "Звертання. Вставні слова. Словосполучення речення",
	"vidokremlene_oznachennja":                             "Відокремлене означення",
	"vidokremlena_obstavyna":                               "Відокремлена обставина",
	"vidokremlenyj_dodatok":                                "Відокремлений додаток",
	"skladnosurjadne_rechennja":                            "Складносурядне речення",
	"skladnopidrjadne_rechennja":                           "Складнопідрядне речення",
	"skladne_bezspoluchnykove_rechennja":                   "Складне безсполучникове речення",
	"skladne_rechennja_z_riznymy_vydamy_zvjazku":           "Складне речення з різними видами зв'язку",
	"chuzhe_movlennya":                                     "Чуже мовлення",

	//biology
	"biologhija_jak_nauka_pro_zhyve":                                       "Біологія як наука про живе",
	"elementnyj_sklad_klityny_neorghanichni_spoluky":                       "Елементний склад клітини. Неорганічні сполуки",
	"orghanichni_spoluky_vughlevody_lipidy":                                "Органічні сполуки. Вуглеводи, ліпіди",
	"orghanichni_spoluky_bilky":                                            "Органічні сполуки. Білки",
	"orghanichni_spoluky_nukleyinovi_kysloty_atf":                          "Органічні сполуки. Нуклеїнові кислоти. АТФ",
	"metody_doslidzhennja_klityny_eukariotychna_klityna_klitynni_membrany": "Методи дослідження клітини. Еукаріотична клітина. Клітинні мембрани",
	"komponenty_cytoplazmy_eukariotychnoyi_klityny":                        "Компоненти цитоплазми еукаріотичної клітини",
	"osoblyvosti_orghanizaciyi_klityn_eukariotiv":                          "Особливості організації клітин еукаріотів",
	"jadro_khromosomy_ponjattja_pro_kariotyp":                              "Ядро. Хромосоми. Поняття про каріотип",
	"obmin_rechovin":                                                         "Обмін речовин",
	"zberezhennja_spadkovoyi_informaciyi":                                    "Збереження спадкової інформації",
	"realizacija_spadkovoyi_informaciyi":                                     "Реалізація спадкової інформації",
	"rozmnozhennja_ta_indyvidualnyj_rozvytok_orghanizmiv":                    "Розмноження та індивідуальний розвиток організмів",
	"ghenetyka_jak_nauka_zakonomirnosti_spadkovosti_orghanizmiv":             "Генетика як наука. Закономірності спадковості організмів",
	"zakonomirnosti_minlivosti":                                              "Закономірності мінливості",
	"selekciya_organizmiv":                                                   "Селекція організмів",
	"virusi_viroyidi_prioni":                                                 "Віруси, віроїди, пріони",
	"prokariotichni_organizmi":                                               "Прокаріотичні організми",
	"vodorosti":                                                              "Водорості",
	"roslini_vegetativni_organi":                                             "Рослини. Вегетативні органи",
	"generativni_organi_roslin":                                              "Генеративні органи рослин",
	"riznomanitnist_i_rozmnozhennja_mokhiv_paporotej_khvoshhiv_plauniv":      "Різноманітність і розмноження мохів, папоротей, хвощів, плаунів",
	"rizmnomanitnist_rozmnozhennja_gholonasinnykh_i_pokrytonasinnykh_roslyn": "Різноманітність і розмноження голонасінних і покритонасінних рослин",
	"ghryby_lyshajnyky":                                                      "Гриби, лишайники",
	"odnoklitinni_organizmi":                                                 "Одноклітинні організми",
	"ghubky_spravzhni_baghatoklitynni_tvaryny_budova_i_zhyttyedijalnist":     "Губки, справжні багатоклітинні тварини. Будова і життєдіяльність",
	"povedinka_tvarin":                                                       "Поведінка тварин",
	"kyshkovoporozhnynni_ploski_chervy_krughli_chervy":                       "Кишковопорожнинні, плоскі, круглі черви",
	"kilchati_chervy_moljusky":                                               "Кільчаті черви, молюски",
	"chlenystonoghi_yikhnja_riznomanitnist":                                  "Членистоногі. Укриття. Різноманітність",
	"khordovi_ryby_amfibiyi_reptyliyi":                                       "Хордові, риби, амфібії, рептилії",
	"ptakhy_yikhnja_riznomanitnist":                                          "Птахи. Укриття. Різноманітність",
	"ssavci_yikhnja_riznomanitnist":                                          "Ссавці. Укриття. Різноманітність",
	"budova_tila_lyudini":                                                    "Будова тіла людини",
	"nervova_reghuljacija_nervova_systema_vyshha_nervova_dijalnist_ljudyny":  "Нервова регуляція. Нервова система. Вища нервова діяльність людини",
	"gumoralna_regulyaciya_endokrinna_sistema":                               "Гуморальна регуляція. Ендокринна система",
	"krov_limfa":                         "Кров, лімфа",
	"imunitet_imunna_sistema":            "Імунітет. Імунна система",
	"dihalna_sistema_lyudini":            "Дихальна система людини",
	"travna_sistema_lyudini":             "Травна система людини",
	"obmin_rechovin_v_organizmi_lyudini": "Обмін речовин в організмі людини",
	"sechovidilna_sistema_lyudini":       "Сечовидільна система людини",
	"shkira_termoregulyaciya":            "Шкіра. Терморегуляція",
	"oporno_ruhova_sistema_lyudini":      "Опорно-рухова система людини",
	"sensorni_sistemi_lyudini":           "Сенсорні системи людини",
	"reprodukciya_lyudini":               "Репродукція людини",
	"ekologichni_chinniki_populyaciya":   "Екологічні чинники. Популяція",
	"ekosistemi":                         "Екосистеми",
	"biosfera_yak_ekosistema":            "Біосфера як екосистема",
	"adaptaciya_biosistem":               "Адаптація біосистем",
	"osnovy_evoljucijnogho_vchennja":     "Основи еволюційного вчення",
}
