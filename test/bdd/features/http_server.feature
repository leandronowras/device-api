## tags:
##  @create
##  @read
##  @update
##  @delete
##  @domain_validation
Feature: http API
  Rule: 
    Creation time cannot be updated.
    Name and brand properties cannot be updated if the device is in use.
    In use devices cannot be deleted.
    Pagination returns supports page & limit (default page=1, limit=10, max limit=100)

  Background:
    Given the API is running reacheable via http

  @id=2
  Scenario: Create a new device
    When I POST "/devices" with json:
      """
      { "name": "iPhone", "brand": "Apple" }
      """
    Then the response code should be 201
    And the response json at "$.name" should be "iPhone"
    And the response json at "$.brand" should be "Apple"
    And the response json has keys: "id", "state", "creation_time"

  @id=3
  Scenario: Fetch a single device
    Given a device exists with name "iPhone" and brand "Apple"
    When I GET "/devices/{id}"
    Then the response code should be 200
    And the response json at "$.name" should be "iPhone"
    And the response json at "$.brand" should be "Apple"
    And the response json has keys: "id", "state", "creation_time"

  @id=4
  Scenario: Fetch all devices  
    Given the API is running
    And there are more than 10 devices stored
    When I GET "/v1/devices?page=1&limit=10"
    Then the response code should be 200
    And the response json should contain 10 devices
    And the response json should include "next_page" and "previous_page" fields

  @id=5
  Scenario: Fetch devices by brand  
    Given the API is running
    And a device exists with name "iPhone" and brand "Apple"
    And a device exists with name "Galaxy" and brand "Samsung"
    When I GET "/v1/devices?brand=Apple"
    Then the response code should be 200
    And the response json should contain 1 device
    And the response json at "$[0].brand" should be "Apple"

  @id=6
  Scenario: Fetch devices by state  
    Given the API is running
    And a device exists with name "iPhone" and brand "Apple"
    And a device exists with name "Galaxy" and brand "Samsung"
    When I GET "/v1/devices?state=available"
    Then the response code should be 200
    And the response json should contain 2 devices
    And the response json at "$[0].state" should be "available"

  Scenario: Fully and/or partally update an existng device
  Scenario: Delete a single device

##| 7 | Feature: Fully update a device (PUT /v1/devices/{id}) | pending | medium | None | N/A |
##| 8 | Feature: Partially update a device (PATCH /v1/devices/{id}) | pending | medium | None | N/A |
##| 9 | Feature: Delete a device (DELETE /v1/devices/{id}) | pending | medium | None | N/A |
