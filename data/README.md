#Disclaimer

Source of this data is https://openflights.org/data.html no Amadeus or Amadeus customer data has been used.

Data has been downloaed and saved here on 30th of August 2018.

## airports database

AirportID	Unique OpenFlights identifier for this airport.
Name	Name of airport. May or may not contain the City name.
City	Main city served by airport. May be spelled differently from Name.
Country	Country or territory where airport is located. See countries.dat to cross-reference to ISO 3166-1 codes.
IATA	3-letter IATA code. Null if not assigned/unknown.
ICAO	4-letter ICAO code. Null if not assigned.
Latitude	Decimal degrees, usually to six significant digits. Negative is South, positive is North.
Longitude	Decimal degrees, usually to six significant digits. Negative is West, positive is East.
Altitude	In feet.
Timezone	Hours offset from UTC. Fractional hours are expressed as decimals, eg. India is 5.5.
DST	Daylight savings time. One of E (Europe), A (US/Canada), S (South America), O (Australia), Z (New Zealand), N (None) or U (Unknown). See also: Help: Time
Tz database time zone	Timezone in "tz" (Olson) format, eg. "America/Los_Angeles".
Type	Type of the airport. Value "airport" for air terminals, "station" for train stations, "port" for ferry terminals and "unknown" if not known. In airports.csv, only type=airport is included.
Source	Source of this data. "OurAirports" for data sourced from OurAirports, "Legacy" for old data not matched to OurAirports (mostly DAFIF), "User" for unverified user contributions. In airports.csv, only source=OurAirports is included.


## airlines database

AirlineID	Unique OpenFlights identifier for this airline.
Name	Name of the airline.
Alias	Alias of the airline. For example, All Nippon Airways is commonly known as "ANA".
IATA	2-letter IATA code, if available.
ICAO	3-letter ICAO code, if available.
Callsign	Airline callsign.
Country	Country or territory where airline is incorporated.
Active	"Y" if the airline is or has until recently been operational, "N" if it is defunct. This field is not reliable: in particular, major airlines that stopped flying long ago, but have not had their IATA code reassigned (eg. Ansett/AN), will incorrectly show as "Y".



## Routes database

Airline	2-letter (IATA) or 3-letter (ICAO) code of the airline.
AirlineID	Unique OpenFlights identifier for airline (see Airline).
SourceAirport	3-letter (IATA) or 4-letter (ICAO) code of the source airport.
SourceAirportID	Unique OpenFlights identifier for source airport (see Airport)
DestinationAirport	3-letter (IATA) or 4-letter (ICAO) code of the destination airport.
DestinationAirportID	Unique OpenFlights identifier for destination airport (see Airport)
Codeshare	"Y" if this flight is a codeshare (that is, not operated by Airline, but another carrier), empty otherwise.
Stops	Number of stops on this flight ("0" for direct)
Equipment	3-letter codes for plane type(s) generally used on this flight, separated by spaces



In mongo we use miniapp DB. That it should be initialized.

Mongo structure should be something like
miniapp.airports
miniapp.airlines
miniapp.routes

To import data we do


## Local monogdb

cleanup openflight airports file:

```shell
$ sed -i "s/\\\\\"/'/g" data/airports.dat
$ mongoimport -d miniapp -c airports --type csv --file data/airports.dat --fieldFile=data/airports_schema.dat
```

same for _airlines_ and _routes_ but no clean-up needed.

```shell
$ mongoimport -d miniapp -c airlines --type csv --file data/airlines.dat --fieldFile=data/airlines_schema.dat
$ mongoimport -d miniapp -c routes --type csv --file data/routes.dat --fieldFile=data/routes_schema.dat
```

## Scheduling

First idea of scheduling comes from [https://www.jetairways.com/en/fr/planyourtravel/flight-schedules.aspx] but it has been largely modified.


## Docker mongodb

In case your mongo is running as a docker image for example via

```shell
$docker run --name mongodb bitnami/mongodb:latest
```

then to load your data you should something like

```shell
$ mongoContainer=$(docker ps -aqf "name=mongodb")
$ mongoIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $mongoContainer)
$ sed -i "s/\\\\\"/'/g" data/airports.dat
$ mongoimport -h $mongoIP -d miniapp -c airports --type csv --file data/airports.dat --fieldFile=data/airports_schema.dat
$ mongoimport -h $mongoIP -d miniapp -c airlines --type csv --file data/airlines.dat --fieldFile=data/airlines_schema.dat
$ mongoimport -h $mongoIP -d miniapp -c routes --type csv --file data/routes.dat --fieldFile=data/routes_schema.dat
```
