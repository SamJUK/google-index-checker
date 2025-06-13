# Google Index Checker

A small Golang script to help identify accidentally indexed subdomains. Using a [Google Custom Search Engine](https://developers.google.com/custom-search/v1/overview) and domain inclusion and exclusion filters. 

> NOTE: Google Custom Search Engines results can sometimes differ from the actual Google results. Dunno why, but occasionally worth running a manual check in google as well. See the following section #Manual Search Query

## Usage

To use this tool, you need to setup our own google custom search engine and associated API key. See the following [Google document](https://developers.google.com/custom-search/v1/overview) on how to do that.

In the below examples, we are using the following scenario:
- Primary domain is `example.com`
- `www.example.com` and `www2.example.com` are the only domains expected to be indexed


### Authentication

You need to provide the Google CSE ID (Custom Search Engine) and Google API key to the script. This can either be done via environment variables, or through CLI flags.

Recommendation, is setting them via a ENV file.

Name | ENV VAR | CLI FLAG | CLI SHORT FLAG
--- | --- | --- | ---
Custom Search Engine ID | `GOOGLE_CSE_ID` | `--google-cse-id` | `-c`
Google API Key | `GOOGLE_API_KEY` | `--google-api-key` | `-k`


### Docker

Assuming you have your Google tokens setup in a `.env` file.

```sh
docker run --env-file .env samjuk/google-index-checker -d example.com -n www.example.com -n www2.example.com
```

## Manual Search Query

Alternatively, you can run the Google search query directly with the following query.

```
site:example.com -site:www.example.com -site:www2.example.com
```