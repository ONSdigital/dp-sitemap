Feature: dp-sitemap write robot file

    Scenario: Write simple robots file
        Given i have the following robot.json:
        """
            {
                "GoogleBot" :
                    {
                        "AllowList": ["/allow1", "/allow2"],
                        "DenyList": ["/deny1", "/deny2"]
                    }
            }
        """
        When i invoke writejson with the sitemaps "www.site1.com/sitemap1,www.site2.com/sitemap2"
        Then the content of the resulting robots file must be
        """

User-agent: GoogleBot
Allow: /allow1
Allow: /allow2
Disallow: /deny1
Disallow: /deny2

sitemap: www.site1.com/sitemap1
sitemap: www.site2.com/sitemap2

        """
