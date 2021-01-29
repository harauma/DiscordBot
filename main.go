package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	// BotName "Bot"という接頭辞がないと401unauthorizedエラーが起きます
	BotName           = ""
	Token             = ""
	stopBot           = make(chan bool)
	vcsession         *discordgo.VoiceConnection
	HelloWorld        = "!helloworld"
	ChannelVoiceJoin  = "!vcjoin"
	ChannelVoiceLeave = "!vcleave"
)

func init() {
	err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		// .env読めなかった場合の処理
		fmt.Println(err)
	}
	BotName = os.Getenv("BotName")
	Token = os.Getenv("Token")
}

func main() {
	//Discordのセッションを作成
	discord, err := discordgo.New()
	fmt.Println(BotName, Token)
	discord.Token = Token
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
	}

	discord.AddHandler(onMessageCreate) //全てのWSAPIイベントが発生した時のイベントハンドラを追加
	// websocketを開いてlistening開始
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}
	defer discord.Close()

	fmt.Println("Listening...")
	<-stopBot //プログラムが終了しないようロック
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	switch {
	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, HelloWorld)):
		//Bot宛に!helloworld コマンドが実行された時
		fmt.Println("Hello worldやで")
		sendMessage(s, m.ChannelID, "Hello world！")

	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, ChannelVoiceJoin)):
		fmt.Println("voicejoinやで")
		//今いるサーバーのチャンネル情報の一覧を喋らせる処理を書いておきますね
		c, err := s.State.Channel(m.ChannelID) //チャンネル取得
		if err != nil {
			log.Println("Error getting channel: ", err)
			return
		}
		//guildChannels, _ := s.GuildChannels(c.GuildID)
		//var sendText string
		//for _, a := range guildChannels{
		//sendText += fmt.Sprintf("%vチャンネルの%v(IDは%v)\n", a.Type, a.Name, a.ID)
		//}
		//sendMessage(s, c, sendText) チャンネルの名前、ID、タイプ(通話orテキスト)をBOTが話す

		//VOICE CHANNEL IDには、botを参加させたい通話チャンネルのIDを代入してください
		//コメントアウトされた上記の処理を使うことでチャンネルIDを確認できます
		vcsession, _ = s.ChannelVoiceJoin(c.GuildID, "802755175089307654", false, false)
		vcsession.AddHandler(onVoiceReceived) //音声受信時のイベントハンドラ

	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, ChannelVoiceLeave)):
		fmt.Println("voiceleaveやで")
		vcsession.Disconnect() //今いる通話チャンネルから抜ける

	case strings.HasPrefix(m.Content, fmt.Sprintf("%s", BotName)):
		sendMessage(s, m.ChannelID, strings.Trim(m.Content, BotName))
	}
}

//メッセージを受信した時の、声の初めと終わりにPrintされるようだ
func onVoiceReceived(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
	log.Print("しゃべったあああああ")
}

//メッセージを送信する関数
func sendMessage(s *discordgo.Session, channelID string, msg string) {
	_, err := s.ChannelMessageSend(channelID, msg)

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}
