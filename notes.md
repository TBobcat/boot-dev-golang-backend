### **Nov. 20th**  
- count handler, reset handler, and middleware to increment req count

~~use middleware to just increment count,~~    
~~handler to write back current request~~  
~~manually tested stateful(save value in memory using struct type) handler endpoints~~

### **Nov.22nd**  
middleware takes a http.handler and returns a http.handler, while adding some logic to it (sorta like a wrapper)  
handlers are basically routers that takes req and send resp to route endpoints

### **Nov.23rd**  
allow http methods only to GET for certain end points
when forbidden methods hit, return 405 (Method not allowed)  

    A type that implements the interface, and methods of the interface, just means such type can use those methods in that interface.
    The implementing type then can be consumed as the interface type as arguments for other functions.
    ref: https://gobyexample.com/interfaces

### **Nov.25th**
- `curl -I http://example.com` puts a HEAD method in request, to use GET explicitly use `-X GET` 

- myRouter.Mount() can mount either http.Handler, or a chi.Router, code spec doesn't say but documentation does...

### **Nov.26th**
use info here: https://www.reddit.com/r/golang/comments/16htjkw/how_to_serve_html_files/

### **Nov.27th**
used chatgpt for code to rendered a html page as a string, and made sure url path was mounted correctly 

### **Dec.4th**
- write route
- write handler(s?) to decode json from request and send a json as response

### **Dec.5th**
- wrote the handler for checking json input and responde accordingly
- right now server is written mainly in 
    - `/api` endpoint with its sub urls mounted to it
    - `/admin/metrics` also mounted
    - `/app` serves `index.html`
    - `/app/assets/logo.png` serves the chirpy picture


### **Dec.6th**
- add code checking logic in json valid checking condition, logic being a helper function
- write helper function to check the json string


### **Dec.8th**
- got inpatient with chatgpt's help, missed a word to censor in code and spent LONG time to figure out regex through chatgpt
- didn't need chatgpt, the punctuation mask can be simply solved by matching exact words, once splitting words by space
- discovered Golang doesn't have built in In function to check if a thing is in the slice of things
- most importantly, try to understand as much as possible the full picture of code, and what the coding task is asking for PATIENTLY on ANY coding task

### **Dec.9th**
- finished reading 1st assignment of Storage chapter
  - write POST and GET endpints to API, which saves data persistently
  - learn proper way to code, to save data to disk (code is essentially database package design) through the link given, and the example Lane provided


### **Dec.13th**
- added internal package (of the root module):
    - making a golang module in `internal` directory,  `go mod init module_name`
    - in root `go.mod` add `require` & `replace` line
    - in root dir's `main.go` add `import internal/module_name`

### *Dec.15th*
- requests and responses both have headers
- all requests methods have response(and response headers) seems

### *Dec.17th*
- self tested chap 5. challenge 1, not sure website tests are not passing
- made `returnVals` and the `jsonState` slice global variables
    - slice of golang is dynamic, don't need to worry about memory bound