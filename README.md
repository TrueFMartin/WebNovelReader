# WebNovelReader
Reads a web-novel and writes to an epub file. "The Wandering Inn" is set as default.

See go.mod for required imports

To run make sure you have Go installed.

$sudo apt install golang-go

after cloning repo in command line: 

$go mod tidy

//this will install required packages

$go run main.go

//this will run the main file.

The book is currently very  long, it may take a while.
The epub will be in the same directory as main.go

To change to another ebook: 
1) modify the "table of contents" url on line 39 in main.go to your new webnovel's Table of Contents
2) modify the "doc.Find(...)" on line 40
    1. Go to the page that holds the table of contents for the web-novel
    2. View source code of webpage.
    3. Before the first listing in the talbe of contents look for a class name, likely in a \<div\>
    4. In "The Wandering Inn" table of contents it was \<div class="entry-content"\>
    5. in doc.Find(...) replace "..." with ".[class-name]" ( as a string, and with the '.' )
3) if needed, change the "selection.Find("p") on line 67
    1. In "The Wandering Inn" on each chapter's page, the first instance of actual story content begins with a \<p\> element. So the .Find() searched for all "p"
4) And of course change the "Title" and "Author" on line 79, 80

TO-DO:

Offer user input for filling in "table of contents url" and class name for each ToC entry
