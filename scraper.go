package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)


func checkForAddOnPromo(url string, words []string) (bool, string, error) {
	possiblePromo := false

	resp, err := http.Get(url)
	if err != nil {
		return possiblePromo, "", err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return possiblePromo, "", err
	}

	addOnText := "Koodo Prepaid Add-ons:\n"
	doc.Find(".add-on-title").Each(func(i int, s *goquery.Selection) {
		addOnText += s.Text() + "\n"
	})

	if contains(words, strings.ToLower(addOnText)) {
		possiblePromo = true
	}
		
	return possiblePromo, addOnText, nil
}

func notifyError(err error, url string) {
	emailSubject := fmt.Sprintf("Error parsing %s", url)
	emailTextBody := fmt.Sprintf("There was an error parsing: %s\n\n", err)

	// Send email
	sendEmail(emailSubject, emailTextBody)
}

func sendEmail(emailSubject string, emailTextBody string) {
	// Get sender/recipient from ENV
	sender := os.Getenv("SENDER")
	recipient := os.Getenv("RECIPIENT")

	sesLocation := os.Getenv("SES_LOCATION")

	// The character encoding for the email.
	emailCharSet := "UTF-8"

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(sesLocation)},
	)

	// Create SES session
	svc := ses.New(sess)

	// Build email
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(emailCharSet),
					Data:    aws.String(emailTextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(emailCharSet),
				Data:    aws.String(emailSubject),
			},
		},
		Source: aws.String(sender),
	}

	// Attempt to send the email
	result, err := svc.SendEmail(input)

	// Print any error messages
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println("Email Sent to address: " + recipient)
	fmt.Println(result)
}

func start(event scrapeData) {
	urlToScrape := event.Url
	containsWords := event.Words
	containsWordsArray := strings.Split(containsWords, ",")

	possiblePromo, results, error := checkForAddOnPromo(urlToScrape, containsWordsArray)

	if (error != nil) {
		notifyError(error, urlToScrape)
	}

	if (possiblePromo) {
		sendEmail("Possible Koodo Prepaid Promotions!", results)
		fmt.Println("Possible Koodo Promotion")

	} else {
		sendEmail("No Koodo Prepaid Promotion", results)
		fmt.Println("No Koodo Promotion")

	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.Contains(e, a) {
			return true
		}
	}
	return false
}

type scrapeData struct {
	Url  string "json:url"
	Words string "json:words"
}

func main() {
	lambda.Start(start)
}