# Job scheduler

Schedule jobs to run on given intervals

## job types

**webhook** - job will send a http request scheduled interval

## example job

Send a get request to example.com every min,

```json
{
  "name": "My first Job",
  "description": "A description of my first job",
  "version": "0.0.1",
  "author": "Nicholas Terrell <email@example.com>",
  "schedule": {
    "timezone": "UTC",
    "dayOfWeek": [],
    "month": [],
    "day": [],
    "hour": [],
    "minute": []
  },
  "type": "webhook",
  "allowOverlap": false,
  "data": {
    "method": "GET",
    "url": "https://example.com",
    "timeout": 10000,
    "expectedResponseStatus": 200,
    "expectedResponseContentType": "text/html",
    "expectedResponse": "ok"
  }
}

```
