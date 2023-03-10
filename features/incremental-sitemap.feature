Feature: URLs can be added to an existing sitemap

    Scenario: Add URL to an existing sitemap
        Given Sitemap "A" looks like the following:
        """
        <?xml version="1.0" encoding="UTF-8"?>
          <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
          <url>
            <loc>http://example.com/1</loc>
            <lastmod>2022-01-01</lastmod>
          </url>
          <url>
            <loc>http://example.com/2</loc>
            <lastmod>2023-02-02</lastmod>
          </url>
        </urlset>
        """
        When I add a URL "http://example.com/3" dated "2024-03-03" to sitemap "A"
        Then the new content of the sitemap "A" should be
        """
        <?xml version="1.0" encoding="UTF-8"?>
        <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
          <url>
            <loc>http://example.com/1</loc>
            <lastmod>2022-01-01</lastmod>
          </url>
          <url>
            <loc>http://example.com/2</loc>
            <lastmod>2023-02-02</lastmod>
          </url>
          <url>
            <loc>http://example.com/3</loc>
            <lastmod>2024-03-03</lastmod>
          </url>
        </urlset>
        """

    Scenario: Add URL to a non-existing sitemap file
        Given Sitemap "B" doesn't exist yet
        When I add a URL "http://example.com/1" dated "2022-01-01" to sitemap "B"
        Then the new content of the sitemap "B" should be
        """
        <?xml version="1.0" encoding="UTF-8"?>
        <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
          <url>
            <loc>http://example.com/1</loc>
            <lastmod>2022-01-01</lastmod>
          </url>
        </urlset>
        """
