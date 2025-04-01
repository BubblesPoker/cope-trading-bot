package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserState struct {
	CurrentMenu string
	// 其他用户状态字段
}

var userStates = make(map[int64]*UserState)

func main() {
	//botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	botToken := "7869146013:AAGYL-txb0sGTDOUUQe5B0yDOZy-roTD0EY"
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN环境变量未设置")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("已授权登录账号 %s", bot.Self.UserName)

	// 设置命令菜单
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "开始使用机器人"},
		{Command: "menu", Description: "显示主菜单"},
		{Command: "help", Description: "获取帮助"},
		{Command: "settings", Description: "设置选项"},
	}

	_, err = bot.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		log.Printf("设置命令菜单失败: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, update.Message)
		} else if update.CallbackQuery != nil {
			handleCallback(bot, update.CallbackQuery)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	if _, exists := userStates[message.Chat.ID]; !exists {
		userStates[message.Chat.ID] = &UserState{CurrentMenu: "main"}
	}

	userState := userStates[message.Chat.ID]

	log.Println("text ", message.Text)
	if message.IsCommand() {
		log.Println("is not command")
		switch message.Command() {
		case "start":
			msg.Text = "欢迎使用本机器人!\n\n使用 /menu 查看主菜单"
		case "menu":
			userState.CurrentMenu = "main"
			msg.Text = " \"☀️欢迎使用DT bsc交易机器人☀️\\n🚀当前交易钱包：\\n💳钱包余额：0.000000 BNB\","
			msg.ReplyMarkup = createMainMenu()
		case "help":
			msg.Text = "帮助信息:\n\n" +
				"/start - 开始使用\n" +
				"/menu - 显示菜单\n" +
				"/help - 获取帮助\n" +
				"/settings - 设置选项"
		case "settings":
			userState.CurrentMenu = "settings"
			msg.Text = "设置菜单:"
			msg.ReplyMarkup = createSettingsMenu()
		default:
			msg.Text = "未知命令，使用 /help 查看可用命令"
		}
	} else {
		// 处理普通消息
		switch userState.CurrentMenu {
		case "waiting_for_feedback":
			msg.Text = "感谢您的反馈!"
			userState.CurrentMenu = "main"
		default:
			msg.Text = "您说: " + message.Text + "\n\n使用 /menu 查看选项"
		}
	}

	if _, err := bot.Send(msg); err != nil {
		log.Printf("发送消息失败: %v", err)
	}
}

func handleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	callbackCfg := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackCfg); err != nil {
		log.Printf("回调处理失败: %v", err)
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "")
	userState := userStates[callback.Message.Chat.ID]

	data := strings.Split(callback.Data, "_")
	menu := data[0]
	action := strings.Join(data[1:], "_")

	switch menu {
	case "main":
		switch action {
		case "profile":
			msg.Text = "您的个人资料..."
		case "search":
			msg.Text = "请输入搜索内容:"
			userState.CurrentMenu = "waiting_for_search"
		case "feedback":
			msg.Text = "请输入您的反馈:"
			userState.CurrentMenu = "waiting_for_feedback"

		case "general":
			//msg.Text = "请输入您的反馈:"
			//userState.CurrentMenu = "waiting_for_feedback"

			msg.Text = "返回主菜单"
			msg.ReplyMarkup = createMainMenu()
			userState.CurrentMenu = "main"

		}
	case "settings":
		switch action {
		case "notifications":
			msg.Text = "通知设置..."
		case "privacy":
			msg.Text = "隐私设置..."
		case "back":
			msg.Text = "返回主菜单"
			msg.ReplyMarkup = createMainMenu()
			userState.CurrentMenu = "main"
		}
	}

	if _, err := bot.Send(msg); err != nil {
		log.Printf("发送消息失败: %v", err)
	}
}

func createMainMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰买/卖", "main_profile"),
			tgbotapi.NewInlineKeyboardButtonData("💵资产", "main_search"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳我的钱包", "main_feedback"),
			tgbotapi.NewInlineKeyboardButtonData("⚙️设置", "main_general"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("钱包分析", "main_feedback"),
			tgbotapi.NewInlineKeyboardButtonData("刷新余额", "main_settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎁邀请返佣", "main_feedback"),
			tgbotapi.NewInlineKeyboardButtonData("🌎Language/语言", "main_settings"),
		),
	)
}

func createSettingsMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("通知设置", "settings_notifications"),
			tgbotapi.NewInlineKeyboardButtonData("隐私设置", "settings_privacy"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("返回主菜单", "settings_back"),
		),
	)
}
