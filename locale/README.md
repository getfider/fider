# How to localize Fider

English files `en/client.json` and `en/server.json` are the only files that should be modified during development.

All other locale files are sourced from [Crowdin](https://crowdin.com/project/fider). If you're looking to contribute with a new translation, a correction or improvement, please use Crowdin.

# How to sync Git with Crowdin

_Note:_ This section is a how-to process for Fider admins only.
_Note 2:_ Would be great to automate this...

## Source Strings

1. Visit https://crowdin.com/project/fider/content/files
2. For each file, click Update and upload the respective file

## Sync English from Git to Crowdin

1. Visit https://crowdin.com/project/fider/content/files
2. For each file, select Update

## Sync other locales from Crowdin to Git

1. Review and approve translations
2. Visit https://crowdin.com/project/fider
3. Select build and download
4. Copy files to locale folder
5. Update percentages on `locales.ts` based on Crowdin numbers
6. Create PR
