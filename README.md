# Call Me Later

Schedule HTTP requests later in time. Also known as Delayed requests or Delayed sending.

## Usage

Call-me-later is "almost" a proxy. or rather a delayed proxy. Call your services "in the future".

Everything stays the same, with some minor tweaks. Call-me-later will always respond with a 202 accepted for any
properly formed request, you will not get a direct response back. Call-me-later needs a few header values in order to
know what to do with the request.

* `X-Later-Request-URL` the URL to call
* `X-Later-Request-When` when to call it, e.g. `1h` to schedule the call 1 hour from now. or `10m5s` for 10 minutes and
  5 seconds.
* `X-Later-Response-URL` optional response callback URL to send any response data to
* `X-Later-Response-Method` optional response callback method to use for the response callback

Any headers and content will be forwarded as is to the target Request URL once the scheduled time is reached. The HTTP
method of the call to the Call-me-later service will be captured and used for the request to the target URL.

Any response from the called service can optionally be sent to a callback URL using the `X-Later-Response-URL`
and `X-Later-Response-Method`. The same pattern applies there, any headers and content will be forwarded as is to the
target URL.

## Persistence

The service will persist all requests using a RequestStorage interface. The requests are then stored until they are
eventually called once the `When` criteria is met.

## Community

Pull requests are welcome. Support for more persistence providers would be welcome.

Example call:

```
curl --location --request GET 'http://localhost:10000/later' \
--header 'X-Later-Request-Url: http://github.com' \
--header 'X-Later-When: 10s' \
--header 'X-Later-Response-Url: http://some-callback-url.com' \
--header 'X-Later-Response-Method: POST' \
```