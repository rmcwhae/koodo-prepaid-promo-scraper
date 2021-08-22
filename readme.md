# Koodo Scraper

A simple scraper to check for promotions in [Koodo prepaid plans](https://www.koodomobile.com/prepaid-plans), where data and minute boosters are occasionally doubled or tripled (only a few times per year in my experience). This will check HTML divs with class `.add-on-title` and see if "bonus" is found. (Hopefully this is a good strategy—it’s untested at the moment as I’m waiting for another promo such as this screenshot from 2017-12-07 to check the HTML.)

Based on: https://aaronvb.com/articles/simple-website-text-scraping-with-go-and-aws-lambda.

![Bonus promotion on 2017-12-07](assets/Koodo-bonus-2017-12-07.png)

## Build

```
$ GOOS=linux GOARCH=amd64 go build -o main scraper.go
$ zip main.zip main
```

## Deploy

Upload the zip file and ensure the handler is set to `main`.

### Environment Variables

Set these in AWS Lambda:

```
RECIPIENT=
SENDER=
SES_LOCATION=
```

Ensure sender/recipient addresses exist SES and are verified. The lambda execution role will also required SES permission.

### Automation

Using AWS Cloudwatch Events, I set an events rule to run this function every day. It should not exceed the free tier of Lambda functions, but the Cloudwatch custom event might cost a few cents.

## Why?

First, see my article [How to Minimize Your Phone Bill](https://russellmcwhae.ca/journal/minimize-cell-phone).

The idea here is to stock up on prepaid add-ons when they are double or triple their standard amounts. Automating this via a scraper means not having to manually check the page at (ir)regular intervals. The actual dollar savings in doing this likely will not exceed a hundred dollars a year, but it was a fun project to learn Go and AWS Lambda functions.
