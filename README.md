## WME Architecture PoC

The main point of this repository is to create working example of WME target architecture.


### Getting started:

First of all you need `go`, `docker` and `docker-compose` installed on your machine.

1. First setup is to setup infrastructure on your local machine by running:

```bash
docker-compose up
```

2. After that you need to create bucket inside the `minio` console. You should be able to access that by going to [http://localhost:9200/](http://localhost:9200/). Login is `admin` and password is `password`. Go to `buckets` page and create one with a name of `wme-data-bk`.

3. After that's done you can start event bridge by running:

```bash
go run bridge/main.go
```

4. While bridge is running you can start collection data in your bucket by running:

```bash
go run store/main.go
```

5. To run streaming examples you need to start by creating the streams by running:

```bash
go run example/create/main.go 
```

6. If streams are successfully created you can run:

```bash
go run example/query/main.go
```

7. If you want to change the streams or add the new one you can `delete` all of the streams by running:

```bash
go run example/delete/main.go
```