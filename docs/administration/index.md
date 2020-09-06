# Administration

It is possible to set up and use your own instance of Cell. Itself, Cell is incredibly lightweight; its two main dependencies include PostgreSQL and Redis, both of which shouldn't add _too_ much to the load. It is also possible to scale Cell horizontally by using `locketd`. Do note that this manual does not cover securing your Cell instance. You will likely want to protect your PostgreSQL and Redis instances with some form of authentication.

This manual assumes that you're running Ubuntu.
