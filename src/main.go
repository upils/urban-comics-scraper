package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const (
	urlBase      = "http://www.urban-comics.com"
	urlChecklist = urlBase + "/a-paraitre"
)

// Issue represent one of the monthly issue
type Issue struct {
	Link       string
	Title      string
	Collection string
	Date       string
	Price      string
}

// Configuration represent the configuration
type Configuration struct {
	MailSender    string
	MailPassword  string
	MailRecipient string
	MailServer    string
	MailPort      string
}

// fetchConf go fetch the configuration from a file and put it in Configuration object
func fetchConf() (configuration Configuration, err error) {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)

	configuration = Configuration{}
	err = decoder.Decode(&configuration)

	return
}

// getIssueslist fetch the monthly issues list
func getIssuesLinklist() (issuesList []string, err error) {

	log.Printf("Get the main comics list from %s\n", urlChecklist)

	resp, err := http.Get(urlChecklist)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err.Error())
	}

	doc.Find(".comics-container").Each(func(i int, s *goquery.Selection) {

		attr, _ := s.Find("a").Attr("href")

		re := regexp.MustCompile("www.urban-comics.com")

		if re.MatchString(attr) {
			issuesList = append(issuesList, attr)
			// log.Printf("Found the link %s\n", attr)
		}
	})

	return issuesList, err
}

func getDetails(issueLink string, c chan Issue, wg *sync.WaitGroup) {
	defer wg.Done()
	issue := Issue{}
	issue.Link = issueLink

	// log.Printf("Get details from %s\n", issueLink)

	resp, err := http.Get(issueLink)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err.Error())
	}

	issue.Title = getTitle(doc)
	issue.Price = getPrice(doc)

	doc.Find("li").Each(func(i int, s *goquery.Selection) {

		if collection := getCollection(s); collection != "" {
			issue.Collection = collection
		}

		if date := getDate(s); date != "" {
			issue.Date = date
		}
	})
	c <- issue
}

// getDetailsIssuesList fetch the issue details such as price, type, availlaibility, date
func getDetailsIssueList(issuesList []string) ([]Issue, error) {

	var detailsIssueslist []Issue
	c := make(chan Issue, len(issuesList))
	var wg sync.WaitGroup
	wg.Add(len(issuesList))

	for _, issueLink := range issuesList {

		go getDetails(issueLink, c, &wg)
	}
	wg.Wait()
	close(c)

	for issue := range c {
		detailsIssueslist = append(detailsIssueslist, issue)
	}

	return detailsIssueslist, nil
}

// getTitle look for the issue title in a <h1> node
func getTitle(doc *goquery.Document) (issueTitle string) {
	issueTitle = doc.Find("h1[id='titre-album']").Text()
	//log.Printf("Title %s\n", issueTitle)
	return
}

// getCollection look for the issue collection in a <li> node
func getCollection(s *goquery.Selection) (issueCollection string) {
	issueCollection = ""
	re := regexp.MustCompile("Collection :")
	content := s.Text()
	if re.MatchString(content) {
		issueCollection = s.Find("a").Text()
		//log.Printf("Collection %s\n", issueCollection)
	}
	return
}

// getDate look for the availlaibility date in a <li> node
func getDate(s *goquery.Selection) (issueDate string) {
	re := regexp.MustCompile("Date de sortie : ")
	content := s.Text()
	if re.MatchString(content) {
		issueDate = strings.Split(content, "Date de sortie : ")[1]
		//log.Printf("Date de sortie %s\n", issueDate)
	}
	return
}

// getPrice look for the issue price in a <div id="prix"> node
func getPrice(doc *goquery.Document) (issuePrice string) {
	content := doc.Find("div[id='prix']").Text()
	issuePrice = strings.Split(content, "Prix : ")[1]
	//log.Printf("Price %s\n", issuePrice)
	return
}

// sendMail build and send the mail
func sendMail(detailsIssues string, configuration Configuration) (err error) {
	from := mail.Address{"", configuration.MailSender}
	to := mail.Address{"", configuration.MailRecipient}
	subj := "Checklist UrbanComics"
	body := detailsIssues

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := configuration.MailServer + ":" + configuration.MailPort

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", configuration.MailSender, configuration.MailPassword, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	c, err := smtp.Dial(servername)
	if err != nil {
		log.Panic(err)
	}

	c.StartTLS(tlsconfig)

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()

	return
}

func renderIssues(detailsIssuelist []Issue) (mailContent string, err error) {
	template := template.Must(template.ParseFiles("./templates/mail.html"))

	var tpl bytes.Buffer
	err = template.Execute(&tpl, detailsIssuelist)
	if err != nil {
		log.Fatal(err.Error())
	}

	mailContent = tpl.String()
	//log.Printf("Resulting mail content: %s\n", mailContent)

	return
}

func main() {
	configuration, err := fetchConf()
	if err != nil {
		log.Fatal(err.Error())
	}

	issuesList, err := getIssuesLinklist()
	if err != nil {
		log.Fatal(err.Error())
	}

	detailsIssuelist, err := getDetailsIssueList(issuesList)
	if err != nil {
		log.Fatal(err.Error())
	}

	detailsIssues, err := renderIssues(detailsIssuelist)
	if err != nil {
		log.Fatal(err)
	}

	err = sendMail(detailsIssues, configuration)
	if err != nil {
		log.Fatal(err)
	}
}
