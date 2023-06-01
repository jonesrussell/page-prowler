
# Streetcode Crawler

Crawl various websites in search of articles for Streetcode


## Authors

- [@jonesrussell](https://www.github.com/jonesrussell)


## License

[MIT](https://choosealicense.com/licenses/mit/)


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`REDIS_HOST`

`REDIS_PORT`

`REDIS_STREAM`


## Run Locally

Clone the project

```bash
git clone https://github.com/jonesrussell/streetcode-crawler.git
```

Go to the project directory

```bash
cd streetcode-crawler
```

### crawler

```bash
cd cmd/crawler
```

Install dependencies

```bash
go install
```

Start the crawler

```bash
go run main.go https://www.sudbury.com
```


## Usage/Examples

```bash
./crawler https://www.sudbury.com
```


## Feedback

If you have any feedback, please reach out to me at https://jonesrussell42.xyz/contact

