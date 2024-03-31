## DistriFS Server

The DistriFS server is a component of the DistriFS network. The server is responsible for the following tasks:
- File discovery (suggesting new files to the indexer)
- Speed benchmarks (the server self-reports speed metrics to the indexer)
- Setting up user passports (see PASSPORT.md in the main folder for more information)
- Serving files

Most responses, excluding file downloads, are expected to have sub-1ms times.

## Running a Server

You can run the DistriFS indexer through the binaries provided in GitHub Actions. After the first run, a config file will be generated for you to edit. Most config settings are self-explanatory; you can learn more about the more complex aspects of the server in the below sections.

## Building

If you don't see your CPU architecture and operating system in the binaries provided on GitHub Actions, or simply want to verify the code for yourself, you can build the server from source.

Required components: Go (this repo is built against Go 1.22)

1. Clone the repository: `git clone https://github.com/JIBSIL/DistriFS`
2. `cd` into the server directory: `cd server`
3. Install modules: `go install`
4. Build for the current system: `go build`; or build for a different system using the `GOARCH` and `GOOS` parameters. A nonexhaustive list of these parameters can be found [here](https://gist.github.com/lizkes/975ab2d1b5f9d5fdee5d3fa665bcfde6)
5. The server binary should be in the same directory that you ran `go build` from (presumably DistriFS/server)

## Authenticating
Users can authenticate with a public key by signing a message provided by the server. This authentication method requires no personal information such as email addresses, passwords, etc. Here is the flow that clients go through when authenticating:
1. Client calls `passport/authenticate` route with no arguments and gets back a key (identifier for that session) and value (value to sign)
2. Client signs the message with the integrated crypto functions in the DistriFS client
3. Client calls the `passport/getKey` route with a public key, signed message and "message key" (the identifier provided by the server in the last request)
4. The server provides the client with a key to use in further communication with that specific server
5. The client calls `passport/verify` with their new key to verify that the signup was successful

After the key is acquired, no further cryptographic functions are needed, as the public key of the client has been verified.

## Anonymous Downloads
All files not specified in the `passportprotectedfiles` option in `config.yaml` are considered to be unprotected. Proactive users may choose to protect all files with `passport` by setting the `allfilespassportlocked` option to true.

## ReadDir Route
When calling the ReadDir route, all files which the user does not have Passport access to are hidden. There is currently no option to prevent this, except perhaps creating a specific directory where passport-protected files reside, such as the `passport` folder in the example directory.

## Indexer Crawling
Indexers and users may crawl your server to get all accessable files. This is a feature of the DistriFS network, but if you want to use DistriFS within your own network (to not allow access from clients outside of your network), you may firewall the server instance or set `allfilespassportlocked: true` in config.yaml.