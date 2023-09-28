# Test task for Effective mobile
description:
Implement a service that will receive a stream of full names (FIO), enrich the response with the most probable age, gender, and nationality from open APIs, 
and save the data in a database. Upon request, provide information about the found individuals. The following needs to be implemented:
  - The service listens to a Kafka queue called "FIO," where information with full names arrives in the following format:
      {
        "name": "Dmitriy",
        "surname": "Ushakov",
        "patronymic": "Vasilevich" // optional
      }
  - In case of an incorrect message, enrich it with an error reason (missing required field, incorrect format, etc.) and send it to the "FIO_FAILED" Kafka queue.
  - Enrich a correct message with:
      - Age - using https://api.agify.io/?name=Dmitriy
      - Gender - using https://api.genderize.io/?name=Dmitriy
      - Nationality - using https://api.nationalize.io/?name=Dmitriy
  - Store the enriched message in a PostgreSQL database (the database structure should be created through migrations).

  - Expose REST endpoints for:
      - Retrieving data with various filters and pagination.
      - Adding new individuals.
      - Deleting by identifier.
      - Modifying entities.
  - Provide GraphQL endpoints similar to REST.
  - Implement data caching in Redis.
  - Cover the code with logs.
  - Cover the business logic with unit tests.
  - Store all configuration data in a .env file.


# Comments for checking:
Solution consists from two sides: agent and server for kafka's processing local check.
Agent can randomly generate messsages for kafka.

Docker:
  - docker-compose up --build
  (*server may not start from the first time due to postgreSQL needs some time to create 'effective_mobile_db')

Postgres:
  - implemented 3 linked tables;
  - request to database doesn't depend on the case of letters;
  - crud operations are made as single transactions;
  - select data is allowed by keys: id, name, surname, patronymic, age, gender, country;
  - supported pagination (page_num, page_size query keys);
  - added migrations.

Redis:
  - implemented read-through cache check (for select requests);
  - in case of other request cache is flushing.

GraphQL: REST:
  - insert request receives json body with person data. Added validation of input data;
  - delete request receives person id through URL query;
  - update request receives person id through URL query and json data in body;
  - select request receives filters, pagination through URL query. Added validation of input data.

GraphQL:
  - available on "/graphql/" route.

Kafka:
  - kafka's messages parsing is implemented on controller conveyor (msgbroker->extradata->save);
  - all these controllers are independent and are communicating through channels;
  - added errors controller for collection error's messages from conveyor and send it back in fail topic.

Postman collection is located in root repository:
- https://github.com/ERupshis/effective_mobile/blob/main/Effective%20Mobile%20test.postman_collection.json

GraphQL queries examples are located in root repository:
- https://github.com/ERupshis/effective_mobile/blob/main/GraphQL%20query%20examples.txt

  Added scripts:
    - for test coverage check: "./test_coverage.sh";
    - documentation: "./generate_docs.sh".
