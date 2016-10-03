# point

Collaborative drawing and messaging web-app written in Golang.
Developed this app in my spare time to get familiarized with golang, websockets, socket programming, web app development. This was done in Fall 2014 during my time in Zoho.

**Head over to [http://sbv-point.herokuapp.com](http://sbv-point.herokuapp.com) to see the app live.**

Server is written entierly in Go. Uses websockets for quick communication.
HTML, JavaScript and HTML Canvas is used in the clientside for drawing and communication.
CSS framework used here is Semantic UI.

## Installation

Install golang. 
    [https://golang.org/doc/install](https://golang.org/doc/install) 
  
In a terminal, run 

    go get github.com/sanjaybv/point

## Running

Once installed, you should be able to run the server by running 

    cd $GOPATH/src/github.com/sanjaybv/point; point
