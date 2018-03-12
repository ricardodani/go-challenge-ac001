# citytravel - ac project

Below you will find the instructions about how to develop and create your project.

## Requirements

You need to create a RESTful service to manage a list of cities and its roads. If there is a road between city A and city B, we say that **A borders B**.

* A unique ID which will be created by your endpoint. This ID MUST be a valid integer (in the 64 bit range).
* A name which might be duplicated
* A list of borders cities. If A is bordered by B, a traveler can go from A to B and from B to A.

The service **MUST** listen for incomming HTTP connnections on port :3000.

**IMPORTANT**: provde the name of the package which contains the binary which will listen on port :3000. The binary **MUST** run without any parameters. Eg.:

```

go build my/package/with/the/binary

./binary

```

# Endpoints

Consider that you have a map of cities with the form

Cities:

* ID: 1 / Name: City 1 / Borders: 3
* ID: 2 / Name: City 2 / Borders: 4
* ID: 3 / Name: City 3 / Borders: 1, 4
* ID: 4 / Name: City 4 / Borders: 2, 3

That data gives the following connections:

1 -> 3 -> 4 -> 2

Here you have a sample of the endpoints required:

* GET /city/1 -> used to get a JSON object with the following format _{"name":"city 1","id":1,"borders":[3]}_
* PUT /city/1 -> used to update data about the city, the input format is the same as the one used to read data about the city. As a response you MUST return the same body as _GET /city/1_
* DELETE /city/1 -> used to delete the city. After deleted a city cannot appear as a _border_ of other city or on _path_ responses.
* POST /city -> can be used to create a new city. The input format is _{"name":"city name","borders":[2, 3]}_. As a response you MUST return _{"name":"city name","id":4,"borders":[2, 3]}_. Also you should return the location of the created city, in the correct http header, ie, the "/city/4".
* GET /cities -> return the list of all cities in the system. As a response you MUST return _{"cities"[{"name":"city 1","id":1,"borders":[3]}, {"name":"city 2","id":2,"borders":[4]}, ...]}_
* DELETE /cities -> used to delete **ALL** the data (all cities and all borders)
* GET /city/1/travel/2 -> it should return one valid path between city 1 and city 2 in the format _{"path":[1, 3, 4, 2]}_. There is no requirement to return the shortest path, as long as the response is valid. If there are no valid paths between city 1 and city 2, you should return a valid status code.
* GET /city/1/travel/2?by=3&by=4 -> it should return one valid path between city 1 and city 2, with stops at the _by_ cities (up to 3 stops). A valid response would be _{"path":[1, 3, 4, 2]}_.

The endpoints need to follow the REST recommendations for HTTP verbs.

## Evaluation

Your code will be evaluated for correctness (all endpoints behave as specified) and for overall organization.

The tests will be executed by _go test_ with _race detector_ on.

You are free to use third-party libraries, but if you do, use the _dep_ to ensure reproducible builds. If you use _goroutines_ their use will be evaluated to. All tests will be executed against the _race detector_ from go.

For storage you can use:

* SQLite (https://github.com/mattn/go-sqlite3 - cgo) or ql (https://github.com/cznic/ql - pure go)
* In memory data structures

Criteria:

* Correctness is more important than completeness. Make sure that your endpoints work as expected even if you don't implement all of them.
* Overall code quality (docs, code formatting, tests)
* Code organization
* Performance

## Concurrency

To measure how good your code runs the endpoints will be called under the following load:

* 1 request per second against your endpoint from the same host
* 10 request per second against your endpoint from the same host
* 100 request per second against your endpoint from the same host

The metrics are evaluated relative to your own code. This means that we'll compare how your system changes from 1 to 100.

Also your _/city/1/travel/2_ endpoint need to return in under than 10 seconds for a database of 100 cities.