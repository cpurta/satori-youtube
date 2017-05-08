# satori-crawler

A web crawler that will can be specified to crawler certain websites based on a regex and to a specified depth.
This is to be used for crawling over YouTube videos and requesting the video info via GoogleAPIs which provides
data for YouTube videos.

## Setup

It is recommended that you have the most recent version of Go installed, which can be found [here](https://golang.org/dl/).
Also you will need to have an account to GoogleAPIs and enable the YoutubeAPI. You will also have to set up an API key. There
is no need to setup OAuth since we are not going to be accessing user data. 
