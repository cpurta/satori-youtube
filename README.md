# satori-crawler

A web crawler that will can be specified to crawler certain websites based on a regex and to a specified depth.
This is to be used for crawling over YouTube videos and requesting the video info via GoogleAPIs which provides
data for YouTube videos.

## Setup

It is recommended that you have the most recent version of Go installed, which can be found [here](https://golang.org/dl/).
Also you will need to have an account to GoogleAPIs and enable the YoutubeAPI. You will also have to set up an API key. There
is no need to setup OAuth since we are not going to be accessing user data.

## Environment Variable Config

Currently environment variables are used to configure the necessary satori
variables to establish a connection to the correct channel. Also this is used
for the connecting to the YouTube API to get video statistics and snippets of
data pertaining to the video. 

It is recommended that you create a environment file that uses bash to export
the needed environment variables.

*example.env*
```
export SATORI_APP_KEY=<your_satori_app_key>
export SATORI_ENDPOINT=<your_satori_endpoint>
export SATORI_ROLE=<your_satori_role>
export SATORI_CHANNEL=<your_satori_channel>
export SATORI_SECRET=<your_satori_secret>
export YOUTUBE_AUTHORIZATION=<your_youtube_api_key>
```

After you have your env file established you can then source the file so that
those variable are now available in your shell. 

```
$ source example.env
```

## Building

You will have to have this project as a part of your `$GOPATH` in order for the
build to work. I recommend creating the following folders in your `GOPATH` under
the `github.com` folder:

`$GOPATH/src/github.com/cpurta/satori`

That should be the file structure that you now have and you can clone/move the
project under the `satori` directory.

If you `cd` in the `satori-youtube` directory and run `pwd` you should have the
following path: `$GOPATH/src/github.com/cpurta/satori/satori-youtube`. With the
exception of the `$GOPATH` being the gopath that you have set up. 

From the project root you can build doing the following:

```
$ mkdir -p ./bin
$ cd ./cmd/go-crawler
$ go build -o ../../bin/crawler
```

If you had a successful build you should be able to change back into the project
root and perform an `ls -l ./bin` and see that the crawler binary has been built
and is now ready to run. 

## Running

The crawler is a pretty versatile tool that can be customized to select the seed
endpoint, the urls to be crawled, the number of crawlers that will concurrently
crawl pages for urls and push data to a satori publisher struct. I suggest that
you use the following command to get started with the using the crawler:

```
$ ./bin/crawler -crawlers 20 -limit 50000 -seed-url \
'https://www.youtube.com/feed/trending' -search='^.*\/watch\?v=[a-zA-Z0-9\-]+$'
```

If you need to familiarize yourself with the flags that can be used, you can use
the `--help` flag to see the usage of each. 
