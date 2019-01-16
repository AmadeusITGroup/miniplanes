To import data in a local db you can run

```shell
mongoimport -d miniapp -c airports --type csv --file airports.dat --fieldFile=../../../../data/airports_schema.dat
mongoimport -d miniapp -c airlines --type csv --file airlines.dat --fieldFile=../../../../data/airlines_schema.dat
mongoimport -d miniapp -c courses --type csv --file courses.dat --fieldFile=../../../../data/courses_schema.dat
```
