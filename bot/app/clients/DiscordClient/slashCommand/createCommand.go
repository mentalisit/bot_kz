package slashCommand

import "github.com/bwmarrin/discordgo"

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
