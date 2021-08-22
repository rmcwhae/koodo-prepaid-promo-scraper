# Koodo Scraper

A simple scraper to check for promotions in [Koodo prepaid plans](https://www.koodomobile.com/prepaid-plans), where data and minute boosters are occasionally doubled or tripled (only a few times per year in my experience). This will check HTML divs with class `.add-on-title` and see if "bonus" is found.

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

## Why?

See my article [How to Minimize Your Phone Bill](https://russellmcwhae.ca/journal/minimize-cell-phone).
