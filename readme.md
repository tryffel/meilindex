# Meilindex

[![License](https://img.shields.io/github/license/tryffel/mailindex.svg)](LICENSE)
[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/tryffel.net/go/meilindex?tab=overview)

Email indexing and extremely fast full-text-search with Meilisearch. 

Features:
* Index mail from Imap or Mbox-file (tested with Thunderbird), store to Meilisearch
* Multiple configurations for different mailboxes
* Configure Meilisearch: stop words, ranking rules order
* Query Meilisearch instance either with CLI or with terminal gui
* Open selected mail in Thunderbird

This is a work-in-progress.

Default config file: ~/meilindex.yaml

# Build
```
go get tryffel.net/go/meilindex
```

# Run
1: Make sure meilisearch is running and accessible
try with e.g. ```curl http://localhost:7700```

2: Create & fill config file

```
meilindex
```
This should create new config file, which is by default at ~/.meilindex.yaml
You can always override config file with '--config'.
Edit config file to suit your needs, insert at least imap and meilisearch settings.

3: (Optional) Customize Meilisearch index before parsing emails, see below

4: Index mail from imap
```
# index INBOX
meilindex index imap 

# index any other folder, e.g. Archive
meilindex index imap --folder Archive
```

5: Query with cli
```
meilindex query my message
meilindex query --folder inbox --subject "item received" my message

```

6: Terminal ui for viewing & queying mail
```
meilindex
```

Gui shortcuts:
* Move between tabs with TAB
* Move up/down list: Key-Up/Key-Down or J/K
* Enter mail with Enter
* Open selected mail in thunderbird with F2 (requires 'thunderlink' add-on)
* Close application with Ctrl-C


# Customize Meilisearch index
Meilisearch features various optimizations and customizations for tailoring search results, 
see [docs](https://docs.meilisearch.com/references/settings.html) for more info. Meilindex supports 
modifying some of them, which hopefully makes the search experience better.

## Stop words
Stop words are irrelevant words in regard to searching content. 
You can create custom stop word lists and let Meilindex 
push them to Meilisearch. Assets-directory contains some example files for stop word lists. These files were 
produced from NLTK language database. You can enable them by calling:
```
# View current stop words
meilindex settings stopwords get

# Push new stop word list
meilindex settings stopwords set assets/stopwords-en.json
```
Do note that only one list (file) can be enabled at a time. If you want to use multiple files, 
you need to combine the stop_word lists into a single file, for now.

## Ranking rules
Ranking is based on a set of rules. Meilisearch provides default set, which you can change to see more relevant
messages first. Creating new rules is not supported yet. See asset/ranking-default.json for format.

```
# View current rules
meilindex settings ranking get

# Reorder rules
meilindex settings ranking set assets/ranking-default.json
```

