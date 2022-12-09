wget https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/6.2.1/openapi-generator-cli-6.2.1.jar -O openapi-generator-cli.jar

go-gin-server

java -jar openapi-generator-cli.jar generate -i videonetics-yojaka-connector-1.0.0-oas3-swagger.yaml \
    -g go-server \
    -o out/openapi \
    --additional-properties=packageName=openapi \
    --additional-properties=router=chi \
    --additional-properties=outputAsLibrary=true

java -jar openapi-generator-cli.jar generate -i videonetics-yojaka-connector-1.0.0-oas3-swagger.yaml ^
    -g go-server ^
    -o out/openapi ^
    --additional-properties=packageName=openapi ^
    --additional-properties=router=chi ^
    --additional-properties=outputAsLibrary=true
