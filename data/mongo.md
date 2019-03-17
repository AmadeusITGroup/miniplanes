

How to print schema

```
$ mongo
> use miniplanes;
> function printSchema(obj) { for (var key in obj) { print(key, typeof obj[key]); } };
>  schemaObj=db.airlines.findOne()
> printSchema(schemaObj)
_id object
airlineID number
name string
alias string
IATA string
ICAO string
callsign string
country string
active string
```
