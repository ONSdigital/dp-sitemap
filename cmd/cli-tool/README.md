# dp-sitemap-cli

This holds all information related to using sitemap-cli tool

## Prerequisites

Set the opensearch signer value using "export OPENSEARCH_SIGNER=true" when running on the sandbox/prod.

## Flags

    --elasticsearch-url string          elastic search api url (default "http://localhost")
    --fake-scroll                       enable fake scroll (default true)
    --robots-file-path string           path to robots file that will be generated (default "test_robots.txt")
    --robots-file-path-reader string    path to robots files that we are reading from (default "./assets/robot/")
    --scroll-size int                   OPENSEARCH_SCROLL_SIZE (default 10)
    --scroll-timeout string             OPENSEARCH_SCROLL_TIMEOUT (default "2000")
    --sitemap-file-path string          path to sitemap file (default "test_sitemap")
    --sitemap-file-path-reader string   path to sitemap files that we are reading from (default "./sitemap/static/")
    --sitemap-index string              OPENSEARCH_SITEMAP_INDEX (default "1")
    --zebedee-url string                zebedee url (default "http://localhost:8082")

## Build Commands

To build for Remote envirionment (sandbox/prod..etc):

```sh
    make build-cli-remote
```

This will create a Linux build

To build for Local envirionment :

```sh
    make build-cli
```

## To run in a remote environment

[Build the tool](#build-commands) for remote environment

Ship to remote:

```sh
    dp scp <env> <mount> ./build/dp-sitemap-cli-remote ./dp-sitemap
```

Remote onto the box.

Now run the tool:

```sh
    export OPENSEARCH_SIGNER=true
    ./dp-sitemap generate --fake-scroll=false --elasticsearch-url=<ElasticSearchURL> --zebedee-url=http://localhost:<ZebedeePort> --sitemap-index="ons"
```

ElasticSearchURL can be obtained from the configs for dp-search-data-importer
ZebedeePort can be obtained from dp-setup

Now copy the sitemap to your machine:

```sh
    dp scp <env> <mount> --pull ./test_sitemap_en.xml .
```
