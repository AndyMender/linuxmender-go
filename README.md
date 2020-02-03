Description
---
A small Web development project for my personal website/blog. *It's largely a work-in-progress!*

Introduction
---

This project is being built using:
- front-end: HTML/CSS with probably [Materialize](https://materializecss.com) and [beego](https://github.com/astaxie/beego)
- back-end: [Go](https://golang.org/)
- Web serving (WIP): [caddy](https://caddyserver.com/) and [Docker](https://www.docker.com/)

The project follows the [MVC design pattern](https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller), though as implemented by the Ruby on Rails project, rather than by Django. I also took some liberty in the way routing controller objects are spawned and generic HTTP responses handled.

Preview
---
If you're curious how the Web page looks at the moment, it can be built and run using Docker (Dockerfile provided). Make sure to expose and map port `8080` properly :).

Requirements
---
* https://github.com/smartystreets/goconvey/convey
* https://github.com/astaxie/beego
