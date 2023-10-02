# dp-sitemap-cli
This holds all information related to using sitemap-cli tool

# Pre-Requisite
Set the opensearch signer value using "export OPENSEARCH_SIGNER=true" when running on the sandbox/prod.

# Global Flags
      --api-url string                    elastic search api url (default "http://localhost")
      --fake-scroll                       enable fake scroll (default true)
      --robots-file-path string           path to robots file that will be generated (default "test_robots.txt")
      --robots-file-path-reader string    path to robots files that we are reading from (default "./assets/robot/")
      --scroll-size int                   OPENSEARCH_SCROLL_SIZE (default 10)
      --scroll-timeout string             OPENSEARCH_SCROLL_TIMEOUT (default "2000")
      --sitemap-file-path string          path to sitemap file (default "test_sitemap")
      --sitemap-file-path-reader string   path to sitemap files that we are reading from (default "./sitemap/static/")
      --sitemap-index string              OPENSEARCH_SITEMAP_INDEX (default "1")
      --zebedee-url string                zebedee url (default "http://localhost:8082")
