Feature: dp-sitemap generates English version sitemap

    Scenario: Generate local sitemap
        Given I index the following URLs:
        | http://example.com/1 | 2022-01-01 |
        | http://example.com/2 | 2023-02-02 |
        When I generate a local sitemap
        Then the content of the resulting sitemap should be
        """
        <?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url>
          <loc>http://example.com/1</loc>
          <lastmod>2022-01-01</lastmod>
        </url>
        <url>
          <loc>http://example.com/2</loc>
          <lastmod>2023-02-02</lastmod>
        </url></urlset>
        """

    Scenario: Generate S3 sitemap
        Given I index the following URLs:
        | http://example.com/1 | 2022-01-01 |
        | http://example.com/2 | 2023-02-02 |
        When I generate S3 sitemap
        Then the content of the S3 sitemap should be
        """
        <?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url>
          <loc>http://example.com/1</loc>
          <lastmod>2022-01-01</lastmod>
        </url>
        <url>
          <loc>http://example.com/2</loc>
          <lastmod>2023-02-02</lastmod>
        </url></urlset>
        """
