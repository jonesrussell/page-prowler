# Page Prowler

Page Prowler is a CLI tool designed for web scraping and data extraction from websites, as well as consuming URLs from a Redis set. It provides two main functionalities:

1. **The 'crawl' command**: This command is used to crawl websites and extract information based on specified search terms.
2. **The 'consume' command**: This command fetches URLs from a Redis set.

Page Prowler is designed to be flexible and easy to use, making it a powerful tool for any data extraction needs.

## Installation

To install Page Prowler, clone the repository and build the binary using the following commands:

```bash
git clone https://github.com/jonesrussell/page-prowler.git
cd page-prowler
go build
```

## Usage

You can use Page Prowler from the command line or through its API.

### Command Line

To search for articles from the command line, use the following command:

```bash
./page-prowler articles --url="https://www.example.com" --searchterms="keyword1,keyword2" --crawlsiteid=siteID --maxdepth=1 --debug
```

Replace `"https://www.example.com"` with the URL you want to crawl, `"keyword1,keyword2"` with the search terms you want to look for, `siteID` with your site ID, and `1` with the maximum depth of the crawl.

### API

To start the API server, use the following command:

```bash
./page-prowler api
```

Then, you can send a POST request to start a crawl:

```bash
curl -X POST -H "Content-Type: application/json" -d '{
 "URL": "https://www.example.com",
 "SearchTerms": "keyword1,keyword2",
 "CrawlSiteID": "siteID",
 "MaxDepth": 3,
 "Debug": true
}' http://localhost:3000/articles/start
```

Again, replace `"https://www.example.com"` with the URL you want to crawl, `"keyword1,keyword2"` with the search terms you want to look for, `siteID` with your site ID, and `3` with the maximum depth of the crawl.

## Configuration

Page Prowler uses a `.env` file for configuration. You can specify the Redis host and password in this file. For example:

```bash
REDIS_HOST=localhost
REDIS_PASSWORD=yourpassword
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.