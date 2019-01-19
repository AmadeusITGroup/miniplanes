
# How to generate schedules data

As soon you've `storage` && `mongo` that is running in your env you can import data.


## locally

```
$MYPORT=9999
$kubectl port-forward pod/$(kubectl get pods -lrole=mongo -o=jsonpath='{.items[0].metadata.name}') ${MYPORT}:27017
```

and in another shell you can execute mongoimporter

```shell
$ mongoimport --port=$MYPORT -d miniapp -c airports --type csv --file airports.dat --fieldFile=../../../../data/airports_schema.dat
$ mongoimport --port=$MYPORT -d miniapp -c airlines --type csv --file airlines.dat --fieldFile=../../../../data/airlines_schema.dat
$ mongoimport --port=$MYPORT -d miniapp -c courses --type csv --file courses.dat --fieldFile=../../../../data/courses_schema.dat
```

Now you can create the schedules running these commands

```shell
$cd .../miniapp/schedules-generator/cmd
$go run main.go --csv-file-name=schedules.csv
```

which generates a file `schedules.csv` that you can import in the usual way

```shell
$ mongoimport --port=$MYPORT -d miniapp -c schedules --type csv --file schedules.dat --fieldFile=../../../../data/schedules_schema.dat
```

if you don't supply `--csv-file-name` file `schedules-generator` will insert directly values in the `schedules` mongo port.

```shell
$ mongo --port=$MYPORT
> use miniapp
switched to db miniapp
> show collections
airlines
airports
courses
schedules
> db.schedules.find()
{ "_id" : ObjectId("5c423da22dc4a51ba4a73dbe"), "origin" : 1382, "destination" : 3797, "flightNumber" : "AF001", "operatingCarrier" : "AF", "daysOperated" : "1234567", "departure" : "1020", "arrival" : "1130", "arriveNextDay" : false }
...
```
you should be able to see `schedules` collection populated

If your storage pod is running you can equally test

# Troubleshooting

Using `kubectl port-foward` you can easily debug `miniapp`

for `storage`

```shell
$ MYPORT=9999
$ kubectl port-forward pod/$(kubectl get pods -lapp=storage -o=jsonpath='{.items[0].metadata.name}') ${MYPORT}:33775
```

and in another shell you can for example just injected schedules like this

```shell
$ http http://localhost:9999/schedules
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 18 Jan 2019 23:20:24 GMT
Transfer-Encoding: chunked

[
    {
        "Arrival": "0001-01-01T00:00:00.000Z",
        "ArriveNextDay": false,
        "DaysOperated": "1234567",
        "Departure": "0001-01-01T00:00:00.000Z",
        "Destination": 3797,
        "FlightNumber": "AF001",
        "OperatingCarrier": "AF",
        "Origin": 1382,
        "ScheduleID": null
    },

```


```shell
$MYPORT=9999
$kubectl port-forward pod/$(kubectl get pods -lapp=itineraries-server -o=jsonpath='{.items[0].metadata.name}') ${MYPORT}:4838
```
and in another shell

```shell

```


