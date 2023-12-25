# Page Prowler

Page Prowler is a tool for finding articles from websites where the URL matches provided terms. It provides functionalities for crawling specific websites, extracting articles that match the provided terms, and consuming URLs from a Redis set.

## Usage

```page-prowler [command]```

## Commands

- api: Starts the API server.
- articles: Crawls specific websites and extracts articles that match the provided terms.
- consume: Consumes URLs from a Redis set.
- help: Displays help about any command.

## Installation

To install Page Prowler, clone the repository and build the binary using the following commands:

```bash
git clone https://github.com/jonesrussell/page-prowler.git
cd page-prowler
go build
```

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