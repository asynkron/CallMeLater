# Call Me Later

Schedule HTTP requests later in time. Also known as Delayed requests or Delayed sending.

## Usage

Make a http request just like you would normally do to call some service. Replace the URL with the URL of this service
and pass the following extra headers:

* `X-Later-Request-URL` the URL to call
* `X-Later-Request-When` when to call it, e.g. `1h` to schedule the call 1 hour from now. or `10m5s` for 10 minutes and
  5 seconds.
* `X-Later-Response-URL` optional response callback URL to send any response data to
* `X-Later-Response-Method` optional response callback method to use for the response callback

Any headers or content will be forwarded as is to the target Request URL once the scheduled time is reached.

Any response from the called service can be sent to a callback URL using the `X-Later-Response-URL`
and `X-Later-Response-Method`.

## Persistence

The service will persist all requests using a RequestStorage interface. The requests are then stored until they are
eventually called once the `When` criteria is met.

## Community

Pull requests are welcome. Support for more persistence providers would be welcome.

Example call:

```
curl --location --request POST 'http://localhost:10000/later' \
--header 'X-Later-Request-Url: http://github.com' \
--header 'X-Later-When: 10s' \
--header 'X-Later-Response-Url: http://github.com'' \
--header 'X-Later-Response-Method: POST' \
--header 'Content-Type: text/plain' \
--data-raw 'Hello'
```