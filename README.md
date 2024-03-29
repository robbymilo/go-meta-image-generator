[![Build Status](https://ci.rmilo.dev/api/badges/robbymilo/go-meta-image-generator/status.svg)](https://ci.rmilo.dev/robbymilo/go-meta-image-generator)

# Meta Image Generator

An open graph meta image generator for social previews.

- Pass the title of your page via URL and get an optimized .jpg returned.
- Customize templates via Golang templates.
- Images are cached and only generated once.
- URLs can be signed via a query param to prevent image bombs.

## To run locally:

- Confirm Go is installed: `go version`
- Clone this repo: `git clone git@github.com:robbymilo/go-meta-image-generator.git && cd go-meta-image-generator`
- Run locally: `make start`
- Visit `localhost:3000/Hello world.jpg` in your browser.

Then, in your website template you can use:

```
<head>
  <meta property="og:image" content="http://localhost:3000/Title of my Page.jpg" />
</head>
```

### Run locally with docker:

```
docker run -p 3000:3000 --name=meta --rm robbymilo/go-meta-image-generator:latest
```

### Updating the template

Follow the instructions above to run locally, then navigate to `localhost:3000/render/Hello world.jpg`. Here you can inspect the template in your browser with dev tools.

### Signed URLs

When running in production, you can sign the URLs to prevent a scripter from crashing your server by generating countless images.

Set the environmental variable `SIGNATURE` to a secret key, ex:

```
docker run -p 3000:3000 -e SIGNATURE=milo --name=meta --rm robbymilo/go-meta-image-generator
```

then visit http://localhost:3000/Hello%20world.jpg?signature=Yjg3Nzc1N2FmZjIxNGU2M2MxMjJkNGM0YmU4ZGM5NTE0NDFjZmJhNmExNzgwOTBjOWZlOTcxOGU5ZTEyYmNhZA==

In your HTML template, you can generate the hash of the URL with your secret key, base64 encode it, then pass it to the `signature` query param.

With [Hugo](https://gohugo.io/):

```
{{ $title := "Hello world" }}
{{ $hashed := base64Encode (sha256 (print $title "milo")) }}
{{ return (print "http://localhost:3000/" $title ".jpg?signature=" $hashed) }}
```

Images will not be returned without the correct signature passed to the query param when the `SIGNATURE` env var is set.

## Why not:

- https://github.com/vercel/og-image Only designed to run on vercel/serverless, and templates are not flexible and cannot be previewed in the browser.

## Todo:

- [ ] Add tests
