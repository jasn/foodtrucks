# Foodtrucks

Hosted at http://ec2-18-219-139-157.us-east-2.compute.amazonaws.com:8080/index

This app is an example of how one can implement finding nearby Points of Interest (POI).
The POI in this example are food trucks parked around San Francisco.
Here we use the data from [DataSF](https://datasf.org/), and in particular the [Food Trucks](https://data.sfgov.org/Economy-and-Community/Mobile-Food-Facility-Permit/rqzj-sfat).

In this exercise I primarily focus on the back-end with a minimal front-end (single html page with a little bit of javascript in there for the communication with the back-end).

The project is in three parts:
* [Foodtrucks](https://github.com/jasn/foodtrucks) -- this repository
* [goors](https://github.com/jasn/goors)
* [gorasp](https://github.com/jasn/gorasp)

I have chosen to make three repositories as I think both goors and gorasp are reusable in other contexts.
They are however all part of this exercise!

## Overview of architecture
The architecture is a classic server-client architecture via the HTTP protocol.
The server is responsible for storing locations of all the food trucks, and answering queries, which is 'give me all food trucks within this rectangle'.
The client is responsible for sending and receiving queries as well as drawing some nearby foodtrucks to the query location.

### Client side
The client code looks at where the user has dragged the map, and issues a query for this rectangle (once the button is clicked).
The query is a GET request to the server with two lat/lng pairs that describe the rectangle.
The query result is a json list where each entry is an object with lat/lng position of a food truck along with a name.
When the client receives the query result, it draws the up to 25 nearest foodtrucks to the (current) center of the map.
The front-end leverages html/javascript(jQuery)/css to accomplish these tasks.
I chose these technologies over, say developing an android app, because I am familiar with these technologies.
I used to work in web-development several years ago. My experience building web-apps was the PHP/MySQL/HTML/CSS/Javascript combination running on an Apache server.
My experience is from before CSS3 and before many of the new javascript frameworks, and before HTML5 became the standard.
However I do not think my solution here suffered much on that behalf, as it required a fairly low amount of code.


### Server side
The server is made of several components.
First there is the actual server, that listens for requests and serves the appriopriate response. This what server.go handles.
But the challenge is how do we answer the queries? I.e. return all food trucks in a rectangle?
This is an instance of the 2D-orthogonal range searching problem.
In our case the data is static (no updates), so we can use a [range tree](https://en.wikipedia.org/wiki/Range_tree).
However out of the box, a range tree has O(log^2 n) queries and O(nlog n) space. This can be improved using [fractional cascading](https://en.wikipedia.org/wiki/Fractional_cascading).
And the running time of the structure I implemented in [goors](https://github.com/jasn/goors) (Go Orthogonal Range Search), has the same running time / space usage.
The goors project was implemented as part of the foodtrucks project. There is a README in the goors project that describes the solution in greater detail.
Further, the goors project depends on the gorasp project, which of course was also made because of this project.
The gorasp project implements a [succinct rank/select](https://en.wikipedia.org/wiki/Succinct_data_structure) data structure. See the README in that project for more information.


# Remarks
First, why this choice of architecture? Honestly it was the first thing that came to my mind.
I am sure there are other perfectly good solutions too, and I will be happy to discuss advantages/disadvantages of them.

I chose to use golang for the implementation.
My primary reason is "Learn something new". I have no experience with golang and this was my first project using it.
The app on its own I think was not enough to fall under the "learn something new", as I do have some prior experience with web-apps.
I also have prior experience with algorithm engineering, which is why the food trucks task was my choice.

The solution currently is a bit weird, in that the front-end sends a query that
describes a box and the server returns every foodtruck in that box. Wouldn't it be more natural to just send the center of the map
(or some other query point) and just return the 25 nearest?
Yes. Yes it would be more natural. I can give reasons why I did not do that, but I do not think any of them are any good. So I would like to change that if I were to write it again.

# Setup
Assuming go is intalled, we can get up and running following these steps:

* `go get github.com/jasn/gorasp`
* `go get github.com/jasn/goors`
* `go get github.com/jasn/foodtrucks`
* `cd $GOPATH/src/github.com/jasn/foodtrucks`
* `go build`

Edit "config" and put in the ip:port to listen on.
Edit "googleapikey" and put in a working google api key (remember to configure it properly via the google interface).

run `./foodtrucks`
