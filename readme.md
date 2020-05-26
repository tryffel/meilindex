# Meilindex

[![License](https://img.shields.io/github/license/tryffel/mailindex.svg)](LICENSE)


* Index mail from imap server or from files
* Store to meilisearch
* Query indexed mails


Default config file: ~/meilindex.yaml

# Build
```
go get tryffel.net/go/meilindex
```

# Run
1. Make sure meilisearch is running and accessible
2. Create & fill config file
```
meilindex
```
This should create new config file, which is by default at ~/.meilindex.yaml
You can always override config file with '--config'.
Edit that, insert at least imap and meilisearch settings.

3. Index mail from imap
```
# index INBOX
meilindex index imap 

# index any other folder, e.g. Archive
meilindex index imap --folder Archive
```

4. Query with cli
```
meilindex query my message
meilindex query --folder inbox --subject "item received" my message

```

5. Terminal ui for viewing & queying mail
```
meilindex
```

Movement inside gui:
* Move between tabs with TAB
* Move up/down list: Key-Up/Key-Down or J/K
* Enter mail with Enter
* Close application with Ctrl-C

