# Meta Image Generator

An open graph meta image generator for social previews.

- Pass the title of your page via URL and get an optimized .jpg returned.
- Easily customize templates via [pug](https://pugjs.org/api/getting-started.html).
- Images are cached and only generated once.
- URLs can be signed via query param to prevent image bombs.

## To run locally:
- Confirm Go is installed: `go version`
- Clone this repo: `git clone git@github.com:robbymilo/go-meta-image-generator.git && cd go-meta-image-generator`
- Run `go mod download`
- Run locally: `make start`
- Visit `localhost:3000/Hello World!.jpg` in your browser.

Then, in your website template you can use:

```
<head>
  <meta property="og:image" content="http://localhost:3000/Title of my Page.jpg" />
</head>
```

### Run locally with docker:

```
docker run -p 3000:3000 --name=meta --rm robbymilo/go-meta-image-generator
```

### Updating the template

Follow the instructions above to run locally, then navigate to `localhost:3000/render/Hello World!.jpg`. Here you can inspect the template in your browser with dev tools.

### Signed URLs

When running in production, you'll want to sign the URLs to prevent a scripter from crashing your server by generating countless images.

Set the environmental variable `SIGNATURE` to a secret key, ex:

```
docker run -p 3000:3000 -e SIGNATURE=milo --name=meta --rm robbymilo/go-meta-image-generator
```

Then, in your HTML template, generate the hash of the URL with your secret key, base64 encode it, then pass it to the `signature` query param.

With [Hugo](https://gohugo.io/):
```
{{ $hashed := base64Encode (sha256 (print .Title "milo")) }}
{{ return (print "http://localhost:3000/" .Title ".jpg?signature=" $hashed) }}
```

Images will not be returned without the correct signature passed to the query param when the `SIGNATURE` env var is set.

## Why not:
https://github.com/vercel/og-image Only designed to run on vercel/serverless, and templates are not flexible and can't be previewed in the browser.
