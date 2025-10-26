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

  Background:
    Given the API is running reacheable via http

  @id=2
  ##| 2 | Feature: Create a new device (POST /v1/devices) | pending | medium | None | N/A |
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
  ##| 3 | Feature: Fetch a single device (GET /v1/devices/{id}) | pending | medium | None | N/A |
  Scenario: Fetch a single device

  @id=4
  ##| 4 | Feature: List devices with pagination (GET /v1/devices) | pending | medium | None | N/A |
  Scenario: Fetch all devices  
    Rule: 
      Pagination

  @id=5
  ##| 5 | Feature: Filter devices by brand | pending | medium | None | N/A |
  Scenario: Fetch devices by brand  
  @id=6
  ##| 6 | Feature: Filter devices by state | pending | medium | None | N/A |
  Scenario: Fetch devices by state  

  Scenario: Fully and/or partally update an existng device
  Scenario: Delete a single device

##| 7 | Feature: Fully update a device (PUT /v1/devices/{id}) | pending | medium | None | N/A |
##| 8 | Feature: Partially update a device (PATCH /v1/devices/{id}) | pending | medium | None | N/A |
##| 9 | Feature: Delete a device (DELETE /v1/devices/{id}) | pending | medium | None | N/A |
