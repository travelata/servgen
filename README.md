## servgen 

`servgen` tool allows to generate a new stub service, that should be adapted for your needs.

Result service contains `sample` CRUD service which is fully functional and can be used 
for figuring out the concept as well as a good start for new service.

#### Installation from source

````
make build
````

#### How to generate a new service 

Go to `./templates/svc/service.json` and modify parameters of newly generating service.
Take the latest libs' tags from the gitlab.

run the command
````
./servgen -d [data-file-path] -t [template-folder-path] -o [output-folder-path]
````
`data-file-path` is a path to `service.json` (file you've modified on the prev step)

`template-folder-path` is a path to template folder

`output-folder-path` specifies a target folder where a new service is generated

example: 
````
./servgen -d ./templates/svc/service.json -t ./templates/svc/_template -o ../../dev/local/good
````

#### How to proceed with a generated service

- go to new service folder

- setup `TRAVELATAROOT` env variable to your project root (parent folder of your service's folder) in `.env` or globally in your environment

- build your service

````
make build
````

- run unit tests
````
make test
````

- initialize a new DB schema (make sure all infrastructure components are up and running)

````
make db-init-schema
````

- run your service

````
make run
````

- run integration tests

````
make test-integration
````

So, it works now and you can go deeper in code and play around.

**Important:** while generated service is fully functional it could be used either for learning or as a starting point.

