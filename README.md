# Server



## Run

```
make run
```


## Usage

GET - Get the maximum time period
```
curl -X GET http://localhost:8080/api/v1/get_observation_time/{model_id}
```
GET - Get a report for a given id
```
curl -X GET http://localhost:8080/api/v1/get_report/{report_id}
```
POST - Write data to reports
```
curl -d '{"report_info":"{report_info}"}' -H "Content-Type: application/json" -X POST http://localhost:8080/api/v1/set_report
