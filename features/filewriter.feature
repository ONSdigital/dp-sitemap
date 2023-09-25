Feature: dp-sitemap write robot file

    Scenario: Write simple robots file
        Given i have my robots config files in the folder "./features/steps/robot/"
        When i invoke writejson with the sitemap "www.site1.com/sitemap1"
        Then the content of the resulting robots file must be
        """

User-agent: GoogleBot
Allow: /allow1
Allow: /allow2
Disallow: /deny1
Disallow: /deny2

sitemap: www.site1.com/sitemap1

        """
