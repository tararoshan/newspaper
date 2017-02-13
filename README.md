# newspaper
simple top level comment filter for HN who's hiring threads

Filters out all internship related posts, and highlights posts in SF (for spring/summer) and DFW (for doing while taking classes).

Install using:
```
go get github.com/cyrusroshan/newspaper
go install github.com/cyrusroshan/newspaper
```

And run with:
```
newspaper
```

# Usage:
newspaper pipes its contents out into your default pager (`echo $PAGER` to find out), so the commands are similar to what you'd use when reading man files.

You can pipe this output out however you want (e.g. `newspaper|cat` as the simplest example).

A few tips on searching:
* `man $PAGER` to find out more info about your default pager
* press q to exit
* type in `/` and a query to search for it (e.g. `BAY AREA` for bay area positions), then use `n` and `N` to navigate backwards and forwards
* use d and u to scroll quickly through the document
* if you're using iterm2, you can press `cmd` and click on links to open them in your default browser
