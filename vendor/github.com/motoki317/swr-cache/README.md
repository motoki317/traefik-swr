# swr cache server

A dead-simple swr (stale-while-revalidate) style cache server.

- Coalesces requests and requests the upstream in the background while in "grace" period.
- Dead-simple as in, ignores all cache-control headers and just caches GET / HEAD requests.
No cache conditions supported (yet).
