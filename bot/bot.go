package bot

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/bwmarrin/discordgo"
)

// Define a struct to represent each link's information
type Order struct {
	OrderNumber int    `json:"orderNumber"`
	OrderTitle  string `json:"orderTitle"`
	OrderLink   string `json:"orderLink"`
	ReleaseDate string `json:"releaseDate"`
}

var gChannelID = "CHANNEL ID HOLDER"

var orders []Order

var BotToken string

var wg sync.WaitGroup

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func titleInJSON(orderTitle string) bool {
	//orderTitle = strings.TrimSpace(orderTitle) // Trim any leading or trailing spaces
	for _, order := range orders {
		if strings.Compare(order.OrderTitle, orderTitle) == 0 { // Case-insensitive comparison
			fmt.Println("No new entries...")
			fmt.Printf("%s , %s", order.OrderTitle, orderTitle)

			return true // title is in JSON

		}
	}
	return false
}

func alertUser(discord *discordgo.Session) {

	var latestOrder Order

	latestOrder = orders[0] // always first for latest

	// Parse the JSON content into a slice of Order structs
	fmt.Printf("@everyone \nMost recent order: Date: %s\nOrder Number: %d\nOrder Title: %s\n", latestOrder.ReleaseDate, latestOrder.OrderNumber, latestOrder.OrderTitle)

	discordMesssageOrderTitle := "Most recent order:" + latestOrder.OrderTitle

	discord.ChannelMessageSend(gChannelID, latestOrder.ReleaseDate)
	discord.ChannelMessageSend(gChannelID, discordMesssageOrderTitle)
	discord.ChannelMessageSend(gChannelID, latestOrder.OrderLink)
	wg.Add(1) // Add to WaitGroup
	go summarize(discord, latestOrder.OrderLink)
}
func prependJSON(orderNumber int, orderTitle string, orderLink string, releaseDate string) {
	fmt.Println("Append json")

	order := Order{
		OrderNumber: orderNumber,
		OrderTitle:  orderTitle,
		OrderLink:   orderLink,
		ReleaseDate: releaseDate,
	}

	// we know  these values are unique to add them to orders array
	orders = append([]Order{order}, orders...)

	//now that theyre in orders array lets

}
func appendJSON(orderNumber int, orderTitle string, orderLink string, releaseDate string) {
	fmt.Println("Append json")

	order := Order{
		OrderNumber: orderNumber,
		OrderTitle:  orderTitle,
		OrderLink:   orderLink,
		ReleaseDate: releaseDate,
	}

	// we know  these values are unique to add them to orders array
	orders = append(orders, order)
	//now that theyre in orders array lets

}
func viewMostRecent(discord *discordgo.Session, channelID string, orderNum int) {
	fmt.Println("Most recent")

	if len(orders) <= 0 {
		discord.ChannelMessageSend(channelID, "Run !order to initalize the bot!")
	}

	if orderNum < 0 {
		for orderNum, order := range orders {
			orderIndex := orderNum + 1
			sOrderIndex := strconv.Itoa(orderIndex)

			message := "[" + sOrderIndex + "]" + " Order Release Date: " + order.ReleaseDate + "\n" + "Title: " + order.OrderTitle + "\n"
			discord.ChannelMessageSend(channelID, message)
			discord.ChannelMessageSend(channelID, "Link: "+"*"+order.OrderLink+"*")

		}
	} else {
		order := orders[orderNum-1]
		message := "Order Release Date: " + order.ReleaseDate + "\n" + "Title: " + order.OrderTitle + "\n" + order.OrderLink
		discord.ChannelMessageSend(channelID, message)
		wg.Add(1) // Add to WaitGroup
		go summarize(discord, order.OrderLink)
	}

}

func summarize(discord *discordgo.Session, link string) {
	resp, err := soup.Get(link) //right now we'll only get the 10 most recent
	if err != nil {
		os.Exit(1)
	}
	doc := soup.HTMLParse(resp)
	pTags := doc.FindAll("p")

	//TRUNC mode to clear file.
	file, err := os.OpenFile("order.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file for truncation:", err)
		return
	}
	file.Close() // Close the file after clearing its contents

	// APPEND mode to append sentences of order
	file, err = os.OpenFile("order.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening file in append mode:", err)
		return
	}
	defer file.Close() // Ensure the file is closed when done

	for _, pTag := range pTags {
		fmt.Println(pTag.Text())
		para := pTag.HTML() //includes br tags which mess everything up.
		re := regexp.MustCompile(`<.*?>`)
		// Replace the matched content with an empty string
		result := re.ReplaceAllString(para, "")
		fmt.Println(result)

		file.Write([]byte(result + "\n"))
	}

	app := "python3"
	arg0 := "./bruh.py" //dont run this as root and  have bruh.py be some whack thing

	cmd := exec.Command(app, arg0)

	cmd.Run()

	fileContents, err := os.ReadFile("summary.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println(string(fileContents))
	discord.ChannelMessageSend(gChannelID, "(all summaries are based off content contained within the Executive Order) Heres a quick TLDR for you: ")
	discord.ChannelMessageSend(gChannelID, string(fileContents))
}

// Optionally, print the updated JSON to the console

// you cna make this quickly get ALL orders and add to file
func checkWebsite(discord *discordgo.Session) {
	resp, err := soup.Get("https://www.whitehouse.gov/presidential-actions/") //right now we'll only get the 10 most recent
	if err != nil {
		os.Exit(1)
	}

	doc := soup.HTMLParse(resp)

	h2Tags := doc.FindAll("h2")

	releaseTimeTags := doc.FindAll("time")

	for orderNumber, h2 := range h2Tags {
		// Check if the class matches
		releaseTime := releaseTimeTags[orderNumber].Text()

		// Find all <a> tags within this <h2>
		links := h2.FindAll("a")

		// Loop through each <a> tag inside the <h2> and print its text and href attribute
		for _, link := range links {
			fmt.Printf("%d", orderNumber)
			fmt.Println("Order Title: ", link.Text(), "| Link:", link.Attrs()["href"], "Time: ", releaseTime)
			orderTitle := link.Text()
			orderLink := link.Attrs()["href"]

			if !titleInJSON(orderTitle) {
				fmt.Println(("NEW ORDER ALERT"))

				if len(orders) < 10 { //if its still initializing & grabbing data for the first time dont dont print, only most recent ones.
					appendJSON(orderNumber, orderTitle, orderLink, releaseTime)
					//alertUser(discord)
				}
				if len((orders)) == 10 {
					appendJSON(orderNumber, orderTitle, orderLink, releaseTime)
					alertUser(discord)
				}
				if len(orders) > 10 {
					alertUser(discord)
					currentTime := time.Now()
					fmt.Printf("-----------NEW ORDER ALERT (%s)-----------", currentTime.Format("2006-01-02 15:04:05"))
					appendJSON(orderNumber, orderTitle, orderLink, releaseTime)
				}
			}

		}
	}

}

// Worker function that performs a task
func worker(discord *discordgo.Session, interval time.Duration, stopChan <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop() // Ensure the ticker is stopped when the worker exits (execution of line is defered)

	for {
		select {
		case <-ticker.C:
			// Perform the task every interval
			if strings.Compare(gChannelID, "CHANNEL ID HOLDER") == 1 {
				fmt.Println("Waiting for user to designate channel to send data in...")

			}
			fmt.Printf("[+] Worker stopping at %s\n", time.Now().Format(time.RFC3339))
			checkWebsite(discord)
		case <-stopChan:
			// Stop the worker gracefully when a signal is received
			fmt.Printf("[+] Worker stopping at %s\n", time.Now().Format(time.RFC3339))
			return
		}
	}
}

func checkOrders(url string, delay int) {

	url = "https://www.whitehouse.gov/presidential-actions/"

}

func Run() {

	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)
	workerStarted = false
	// add a event handler
	discord.AddHandler(messageInput)

	// open session
	discord.Open()
	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")

	//	wokering
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}

func help(discord *discordgo.Session, channelID string) {
	discord.ChannelMessageSend(channelID, "Presidential Executive Order Notifications!\n\t\"!order\" - specify which channel you'd like alerts sent to.")
	discord.ChannelMessageSend(channelID, "\t\"!view\" - view links,dates, and titles of 10 most recent presidential actions")
	discord.ChannelMessageSend(channelID, "\t\"!ls\" - view listing of 10 most recent orders")

	discord.ChannelMessageSend(channelID, "- zion")
}
func order(discord *discordgo.Session, userchannelID string) {
}

var workerStarted bool

func messageInput(discord *discordgo.Session, message *discordgo.MessageCreate) {

	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	//Make sure to enable message intent in bot pane of discord.
	fmt.Printf("Auhtor ID %s\n", message.Author.ID)
	fmt.Printf("Channel ID %s\n", message.ChannelID)
	fmt.Println("Content %s \n", message.Content)

	// respond to user message if it contains `!help` or `!bye`

	if strings.HasPrefix(message.Content, "!help") {
		help(discord, message.ChannelID)
	} else if strings.HasPrefix(message.Content, "!order") {

		discord.ChannelMessageSend(message.ChannelID, "Channel selected for notifications!")
		discord.ChannelMessageSend(message.ChannelID, "Give me a sec... I'm slow for some reason that I dont fully understand i think its the interval to check for new orders messing with literally everything...")

		gChannelID = message.ChannelID
		//ORDEr CHECKEr
		checkInterval := 20 * time.Second
		stopChan := make(chan struct{})

		if !workerStarted {
			//discord.ChannelMessageSend(message.ChannelID, "Worker started")
			wg.Add(1)
			go worker(discord, checkInterval, stopChan)
			workerStarted = true
		}

	} else if strings.HasPrefix(message.Content, "!view") {
		parts := strings.Split(message.Content, " ")
		if len(parts) > 1 {
			number, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error converting to integer:", err)
				return
			}
			if number <= len(orders) && number >= 0 {
				wg.Add(1)
				go viewMostRecent(discord, message.ChannelID, number)
			}
		}
	} else if strings.HasPrefix(message.Content, "!ls") {
		wg.Add(1)
		go viewMostRecent(discord, message.ChannelID, -1)
	} else if strings.HasPrefix(message.Content, "!ping") {
		discord.ChannelMessageSend(message.ChannelID, "@everyone")
	}
}
