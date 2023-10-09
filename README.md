# ind-download

Started as a copy of gist https://gist.github.com/sug0/214927847bb27e62cc2b3ae6fb1254d2#file-inddownload-go

This is basically just a fix after changes made by indd and make it work on mac.

## How to use

1. Install Go runtime
2. Install wkhtmltopdf. On mac, just run `brew install wkhtmltopdf`
3. Create a file with the list of UUIDs of documents that you would like to download
4. Open terminal and execute `go run main.go <name_of_the_file>`
5. Downloaded documents will be stored in `out/` directory
