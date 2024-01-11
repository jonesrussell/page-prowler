# Page Prowler

Page Prowler is a tool designed to find and extract links from websites based on specified terms. It allows direct interaction through the command-line interface or the initiation of the Echo web server, which exposes an API. This API utilizes the Asynq library to manage queued crawl jobs.

## Usage

```page-prowler [command]```

## Commands

- **api**: Starts the API server.
- **matchlinks**: Crawls specific websites and extracts matchlinks that match the provided terms. Can be run from the command line or via a POST request to `/v1/matchlinks` on the API server.
- **clearlinks**: Clears the Redis set for a given crawlsiteid.
- **getlinks**: Gets the list of links for a given crawlsiteid.
- **worker**: Starts the Asynq worker.
- **help**: Displays help about any command.

## Building

To install Page Prowler, clone the repository and build the binary using the following commands:

```
git clone https://github.com/jonesrussell/page-prowler.git
cd page-prowler
go build
```

Alternatively, you can use the provided Makefile to build the project:

```make all```

This command will run fmt, lint, test, and build targets defined in the Makefile.

### Command Line

To search for matchlinks from the command line, use the following command:

```bash
./page-prowler matchlinks --url="https://www.example.com" --searchterms="keyword1,keyword2" --crawlsiteid=siteID --maxdepth=1 --debug
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
}' http://localhost:3000/matchlinks
```

Again, replace `"https://www.example.com"` with the URL you want to crawl, `"keyword1,keyword2"` with the search terms you want to look for, `siteID` with your site ID, and `3` with the maximum depth of the crawl.

## Configuration

Page Prowler uses a `.env` file for configuration. You can specify the Redis host and password in this file. For example:

```bash
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_AUTH=yourpassword
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
