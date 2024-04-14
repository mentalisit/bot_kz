package slashCommand

import "github.com/bwmarrin/discordgo"

func AddSlashCommandModuleWeaponLocale() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "module",
			Description: "Select the desired module and level",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian:   "модули",
				discordgo.Ukrainian: "модулі",
				discordgo.EnglishUS: "module",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian:   "Выберите нужный модуль и уровень",
				discordgo.Ukrainian: "Виберіть потрібний модуль та рівень",
				discordgo.EnglishUS: "Select the desired module and level",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "module",
					Description: "Select module",
					NameLocalizations: map[discordgo.Locale]string{
						discordgo.Russian:   "модули",
						discordgo.Ukrainian: "модулі",
						discordgo.EnglishUS: "module",
					},
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian:   "Выберите модуль",
						discordgo.Ukrainian: "Виберіть модуль",
						discordgo.EnglishUS: "Select module",
					},
					Required: true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "RSE",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Ингибитор КЗ",
								discordgo.Ukrainian: "Інгібітор ЧЗ",
								discordgo.EnglishUS: "RSE",
							},
							Value: "RSE",
						},
						{
							Name: "Genesis",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Генезис",
								discordgo.Ukrainian: "Генезис",
								discordgo.EnglishUS: "Genesis",
							},
							Value: "GENESIS",
						},
						{
							Name: "Enrich",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Обогатить",
								discordgo.Ukrainian: "Збагатити",
								discordgo.EnglishUS: "Enrich",
							},
							Value: "ENRICH",
						},
						// Добавьте другие модули по мере необходимости
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "level",
					Description: "Select level",
					NameLocalizations: map[discordgo.Locale]string{
						discordgo.Russian:   "уровень",
						discordgo.Ukrainian: "рівень",
						discordgo.EnglishUS: "level",
					},
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian:   "Выберите уровень",
						discordgo.Ukrainian: "Виберіть рівень",
						discordgo.EnglishUS: "Select level",
					},
					Required: true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "Level 0",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 0",
								discordgo.Ukrainian: "Рівень 0",
								discordgo.EnglishUS: "Level 0",
							},
							Value: 0,
						},
						{
							Name: "Level 1",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 1",
								discordgo.Ukrainian: "Рівень 1",
								discordgo.EnglishUS: "Level 1",
							},
							Value: 1,
						}, {
							Name: "Level 2",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 2",
								discordgo.Ukrainian: "Рівень 2",
								discordgo.EnglishUS: "Level 2",
							},
							Value: 2,
						}, {
							Name: "Level 3",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 3",
								discordgo.Ukrainian: "Рівень 3",
								discordgo.EnglishUS: "Level 3",
							},
							Value: 3,
						}, {
							Name: "Level 4",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 4",
								discordgo.Ukrainian: "Рівень 4",
								discordgo.EnglishUS: "Level 4",
							},
							Value: 4,
						}, {
							Name: "Level 5",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 5",
								discordgo.Ukrainian: "Рівень 5",
								discordgo.EnglishUS: "Level 5",
							},
							Value: 5,
						}, {
							Name: "Level 6",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 6",
								discordgo.Ukrainian: "Рівень 6",
								discordgo.EnglishUS: "Level 6",
							},
							Value: 6,
						}, {
							Name: "Level 7",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 7",
								discordgo.Ukrainian: "Рівень 7",
								discordgo.EnglishUS: "Level 7",
							},
							Value: 7,
						}, {
							Name: "Level 8",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 8",
								discordgo.Ukrainian: "Рівень 8",
								discordgo.EnglishUS: "Level 8",
							},
							Value: 8,
						}, {
							Name: "Level 9",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 9",
								discordgo.Ukrainian: "Рівень 9",
								discordgo.EnglishUS: "Level 9",
							},
							Value: 9,
						}, {
							Name: "Level 10",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 10",
								discordgo.Ukrainian: "Рівень 10",
								discordgo.EnglishUS: "Level 10",
							},
							Value: 10,
						}, {
							Name: "Level 11",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 11",
								discordgo.Ukrainian: "Рівень 11",
								discordgo.EnglishUS: "Level 11",
							},
							Value: 11,
						}, {
							Name: "Level 12",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 12",
								discordgo.Ukrainian: "Рівень 12",
								discordgo.EnglishUS: "Level 12",
							},
							Value: 12,
						}, {
							Name: "Level 13",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 13",
								discordgo.Ukrainian: "Рівень 13",
								discordgo.EnglishUS: "Level 13",
							},
							Value: 13,
						}, {
							Name: "Level 14",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 14",
								discordgo.Ukrainian: "Рівень 14",
								discordgo.EnglishUS: "Level 14",
							},
							Value: 14,
						}, {
							Name: "Level 15",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Уровень 15",
								discordgo.Ukrainian: "Рівень 15",
								discordgo.EnglishUS: "Level 15",
							},
							Value: 15,
						},
						// Добавьте другие уровни по мере необходимости
					},
				},
			},
		},
		{
			Name:        "weapon",
			Description: "Select your main weapon",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian:   "оружие",
				discordgo.Ukrainian: "зброя",
				discordgo.EnglishUS: "weapon",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Russian:   "Выберите оружие",
				discordgo.Ukrainian: "Виберіть основну зброю",
				discordgo.EnglishUS: "Select your main weapon",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "weapon",
					Description: "Select weapon",
					NameLocalizations: map[discordgo.Locale]string{
						discordgo.Russian:   "оружие",
						discordgo.Ukrainian: "зброя",
						discordgo.EnglishUS: "weapon",
					},
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Russian:   "Выберите оружие",
						discordgo.Ukrainian: "Виберіть зброю",
						discordgo.EnglishUS: "Select weapon",
					},
					Required: true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "Barrage",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Артобстрел",
								discordgo.Ukrainian: "Артилерія",
								discordgo.EnglishUS: "Barrage",
							},
							Value: "barrage",
						},
						{
							Name: "Laser",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Лазер",
								discordgo.Ukrainian: "Лазер",
								discordgo.EnglishUS: "Laser",
							},
							Value: "laser",
						},
						{
							Name: "Chain ray",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Цепной луч",
								discordgo.Ukrainian: "Ланцюговий промінь",
								discordgo.EnglishUS: "Chain ray",
							},
							Value: "chainray",
						},
						{
							Name: "Battery",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Батарея",
								discordgo.Ukrainian: "Батарея",
								discordgo.EnglishUS: "Battery",
							},
							Value: "battery",
						},
						{
							Name: "Mass battery",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Залповая батарея",
								discordgo.Ukrainian: "Залпова батарея",
								discordgo.EnglishUS: "Mass battery",
							},
							Value: "massbattery",
						},
						{
							Name: "Dart launcher",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Пусковая установка",
								discordgo.Ukrainian: "Пускова установка",
								discordgo.EnglishUS: "Dart launcher",
							},
							Value: "dartlauncher",
						},
						{
							Name: "Rocket launcher",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Ракетная установка",
								discordgo.Ukrainian: "Ракетна установка",
								discordgo.EnglishUS: "Rocket launcher",
							},
							Value: "rocketlauncher",
						},
						{
							Name: "Remove weapon",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.Russian:   "Удалить оружие",
								discordgo.Ukrainian: "Видалити зброю",
								discordgo.EnglishUS: "Remove weapon",
							},
							Value: "Remove",
						},
					},
				},
			},
		},
	}
}

func AddSlashCommandRu() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "модули",
			Description: "Выберите нужный модуль и уровень",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "модуль",
					Description: "Выберите модуль",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Ингибитор КЗ",
							Value: "RSE",
						},
						{
							Name:  "Генезис",
							Value: "GENESIS",
						},
						{
							Name:  "Обогатить",
							Value: "ENRICH",
						},
						// Добавьте другие модули по мере необходимости
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "уровень",
					Description: "Выберите уровень",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Уровень 0",
							Value: 0,
						},
						{
							Name:  "Уровень 1",
							Value: 1,
						}, {
							Name:  "Уровень 2",
							Value: 2,
						}, {
							Name:  "Уровень 3",
							Value: 3,
						}, {
							Name:  "Уровень 4",
							Value: 4,
						}, {
							Name:  "Уровень 5",
							Value: 5,
						}, {
							Name:  "Уровень 6",
							Value: 6,
						}, {
							Name:  "Уровень 7",
							Value: 7,
						}, {
							Name:  "Уровень 8",
							Value: 8,
						}, {
							Name:  "Уровень 9",
							Value: 9,
						}, {
							Name:  "Уровень 10",
							Value: 10,
						}, {
							Name:  "Уровень 11",
							Value: 11,
						}, {
							Name:  "Уровень 12",
							Value: 12,
						}, {
							Name:  "Уровень 13",
							Value: 13,
						}, {
							Name:  "Уровень 14",
							Value: 14,
						}, {
							Name:  "Уровень 15",
							Value: 15,
						},
					},
				},
			},
		},
		{
			Name:        "оружие",
			Description: "Выберите основное оружие",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "оружие",
					Description: "Выберите оружие",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Артобстрел",
							Value: "barrage",
						},
						{
							Name:  "Лазер",
							Value: "laser",
						},
						{
							Name:  "Цепной луч",
							Value: "chainray",
						},
						{
							Name:  "Батарея",
							Value: "battery",
						},
						{
							Name:  "Залповая батарея",
							Value: "massbattery",
						},
						{
							Name:  "Пусковая установка",
							Value: "dartlauncher",
						},
						{
							Name:  "Ракетная установка",
							Value: "rocketlauncher",
						},
						{
							Name:  "Удалить оружие",
							Value: "Remove",
						},
						// Добавьте другие модули по мере необходимости
					},
				},
			},
		},
	}
}
func AddSlashCommandEn() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "module",
			Description: "Select the desired module and level",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "module",
					Description: "Select module",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "RSE",
							Value: "RSE",
						},
						{
							Name:  "Genesis",
							Value: "GENESIS",
						},
						{
							Name:  "Enrich",
							Value: "ENRICH",
						},
						// Добавьте другие модули по мере необходимости
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "level",
					Description: "Select level",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Level 0",
							Value: 0,
						},
						{
							Name:  "Level 1",
							Value: 1,
						}, {
							Name:  "Level 2",
							Value: 2,
						}, {
							Name:  "Level 3",
							Value: 3,
						}, {
							Name:  "Level 4",
							Value: 4,
						}, {
							Name:  "Level 5",
							Value: 5,
						}, {
							Name:  "Level 6",
							Value: 6,
						}, {
							Name:  "Level 7",
							Value: 7,
						}, {
							Name:  "Level 8",
							Value: 8,
						}, {
							Name:  "Level 9",
							Value: 9,
						}, {
							Name:  "Level 10",
							Value: 10,
						}, {
							Name:  "Level 11",
							Value: 11,
						}, {
							Name:  "Level 12",
							Value: 12,
						}, {
							Name:  "Level 13",
							Value: 13,
						}, {
							Name:  "Level 14",
							Value: 14,
						}, {
							Name:  "Level 15",
							Value: 15,
						},
						// Добавьте другие уровни по мере необходимости
					},
				},
			},
		},
		{
			Name:        "weapon",
			Description: "Select your main weapon",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "weapon",
					Description: "Select weapon",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Barrage",
							Value: "barrage",
						},
						{
							Name:  "Laser",
							Value: "laser",
						},
						{
							Name:  "Chain ray",
							Value: "chainray",
						},
						{
							Name:  "Battery",
							Value: "battery",
						},
						{
							Name:  "Mass battery",
							Value: "massbattery",
						},
						{
							Name:  "Dart launcher",
							Value: "dartlauncher",
						},
						{
							Name:  "Rocket launcher",
							Value: "rocketlauncher",
						},
						{
							Name:  "Remove weapon",
							Value: "Remove",
						},
					},
				},
			},
		},
	}
}
func AddSlashCommandUa() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "модулі",
			Description: "Виберіть потрібний модуль та рівень",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "модуль",
					Description: "Виберіть модуль",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Інгібітор ЧЗ",
							Value: "RSE",
						},
						{
							Name:  "Генезис",
							Value: "GENESIS",
						},
						{
							Name:  "Збагатити",
							Value: "ENRICH",
						},
						// Добавьте другие модули по мере необходимости
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "уровень",
					Description: "Выберите уровень",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Рівень 0",
							Value: 0,
						}, {
							Name:  "Рівень 1",
							Value: 1,
						}, {
							Name:  "Рівень 2",
							Value: 2,
						}, {
							Name:  "Рівень 3",
							Value: 3,
						}, {
							Name:  "Рівень 4",
							Value: 4,
						}, {
							Name:  "Рівень 5",
							Value: 5,
						}, {
							Name:  "Рівень 6",
							Value: 6,
						}, {
							Name:  "Рівень 7",
							Value: 7,
						}, {
							Name:  "Рівень 8",
							Value: 8,
						}, {
							Name:  "Рівень 9",
							Value: 9,
						}, {
							Name:  "Рівень 10",
							Value: 10,
						}, {
							Name:  "Рівень 11",
							Value: 11,
						}, {
							Name:  "Рівень 12",
							Value: 12,
						}, {
							Name:  "Рівень 13",
							Value: 13,
						}, {
							Name:  "Рівень 14",
							Value: 14,
						}, {
							Name:  "Рівень 15",
							Value: 15,
						},
					},
				},
			},
		},
		{
			Name:        "зброя",
			Description: "Виберіть основну зброю",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "зброя",
					Description: "Виберіть зброю",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Артилерія",
							Value: "barrage",
						},
						{
							Name:  "Лазер",
							Value: "laser",
						},
						{
							Name:  "Ланцюговий промінь",
							Value: "chainray",
						},
						{
							Name:  "Батарея",
							Value: "battery",
						},
						{
							Name:  "Залпова батарея",
							Value: "massbattery",
						},
						{
							Name:  "Пускова установка",
							Value: "dartlauncher",
						},
						{
							Name:  "Ракетна установка",
							Value: "rocketlauncher",
						},
						{
							Name:  "Видалити зброю",
							Value: "Remove",
						},
						// Добавьте другие модули по мере необходимости
					},
				},
			},
		},
	}
}
