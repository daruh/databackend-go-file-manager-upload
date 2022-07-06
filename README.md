### HTTP sending files

#### HTTP multipart request in RFC and W3

Multipart Content W3C
`https://www.w3.org/Protocols/rfc1341/7_2_Multipart.html`

 Returning Values from Forms: multipart/form-data
`https://datatracker.ietf.org/doc/html/rfc7578`


Content-Disposition Header Field
`https://datatracker.ietf.org/doc/html/rfc2183#section-2.3`

 Multipurpose Internet Mail Extensions(MIME) Part Two: Media Types
`https://datatracker.ietf.org/doc/html/rfc2046`

  MIME (Multipurpose Internet Mail Extensions) Part Three:
              Message Header Extensions for Non-ASCII Text
`https://datatracker.ietf.org/doc/html/rfc2047`

Key points:
- sendig multiple and different data (Content-Types) in one HTTP body through http call constructed in code or even directly from form
- need of extending platform gateway as well as eventcontext in eventbus-go codebase to handle this type of Content-Type (multipart)
- can send multiple files with their metadata (but we probably limit to one file) - we could produce for each file one event to nats
- example transfer of file + metadata to event could be by using Extension part of Event (let's call it "metadata")
- well supported in go by using package multipart , https://pkg.go.dev/mime/multipart


### Go tus
`https://github.com/eventials/go-tus/blob/05d0564bb571e81045012756065a8d002d717caf/upload.go#L49`
`https://tus.io/faq.html`

### Astaxie
`https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.5.html`

### Mohitcare
`https://www.mohitkhare.com/blog/file-upload-golang/`

### Multipart requests
`https://ayada.dev/posts/multipart-requests-in-go/`

`https://github.com/abvarun226/blog-source-code/blob/master/multipart-requests-in-go/multipart-related/server/main.go`

### Upload google drive
`https://freshman.tech/snippets/go/multipart-upload-google-drive/`

### Upload files
`https://gist.github.com/ayoisaiah/1e921f0934f5973b9f83e4060caf865a`

### Zupzup
`https://github.com/zupzup/golang-http-file-upload-download/blob/main/upload.gtpl`

### Medium
`https://medium.com/akatsuki-taiwan-technology/uploading-large-files-in-golang-without-using-buffer-or-pipe-9b5aafacfc16`


### Calling multipartrelated
`go run .\client.go  "C:\SNOW_FILES\test_pdf.pdf" `

