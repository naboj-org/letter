<div align="center">
    <img src="https://user-images.githubusercontent.com/11409143/172726035-c993eaa9-7fac-4700-91c8-857351cfee2c.png" width="128" height="128" />
    <h3>Letter</h3>
    <p>LaTeX document generation over HTTP</p>
</div>

## About

Letter is a simple Go webserver that allows generation of LaTeX documents.
Designed to be run in a Docker container.

## Usage

Letter accepts input files in a `.zip` archive. This archive should contain
all resources needed for your document. Letter builds the document with [Tectonic](https://tectonic-typesetting.github.io/en-US/)
and provides you with the resulting PDF file.

### Synchronous API `POST /sync`

This endpoint builds the PDF synchronously and sends it in the HTTP response.
Expects multipart form request with the following fields:
- `file` - the `.tar` archive containing all input files and resources
- `entrypoint` - name of `.tex` file which will be passed to LuaLaTeX

On success, Letter returns the generated PDF file. When an error happens,
Letter returns JSON response containing `error` (error description)
and `tectonic_output` (Tectonic log).

### Asynchronous API `POST /async`

This endpoint builds the PDF in the background and sends it in callback request.
Expects the same fields as the synchronous API, plus:
- `callback` - URL to which will Letter send the resulting PDF

The callback URL will receive a POST request from Letter with multipart data:
- `file` - the resulting PDF
- `tectonic_output` - Tectonic log
- `error` - error description (if available)

### Authentication

Letter can be configured to require authentication. We expect `X-Token` header
to contain a authentication token equal to the value of `AUTH_TOKEN` environment
variable.

If `AUTH_TOKEN` environment variable is not set, Letter will not require authentication.

## Running

A `Dockerfile` is provided that builds and runs Letter. Letter runs on port 8080 by default,
this can be changed using `PORT` environment variable.

### Requirements

- GoLang 1.21 (or newer) 
- Tectonic
