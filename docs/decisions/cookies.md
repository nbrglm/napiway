# why no cookie support?

We do not support cookies, yet, as they are not integral to the purpose of this SDK and Server generation tool. Our primary focus is on creating SDKs that are used between servers, microservices, etc. where cookies are less relevant.

But, if we do expand our scope to include web applications or browser-based interactions, they will be implemented as below:

- When added, use typed cookie wrappers, NOT raw+typed duplication
  ```go
  type StringCookie struct {
      Val string
      Raw *http.Cookie
  }
  ```

- Request/response structs use the wrapper
- Generator handles string ↔ typed conversion
- Raw is metadata only (path, domain, expiry)