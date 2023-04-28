---
draft: false
title: "Mapping Aardwolf with Graphviz and Golang"
aliases: ["Mapping Aardwolf with Graphviz"]
series: ["mapping-aardwolf"]
date: "2023-04-06"
author: "Nick Dumas"
cover: ""
keywords: [""]
description: "Maxing out your CPU for fun and profit with dense graphs, or how I'm attempting to follow through on my plan to work on projects with more visual outputs"
showFullContent: false
tags:
- graphviz
- graph
- aardwolf
- golang
---

## Textual Cartography
Aardwolf has a fairly active developer community, people who write and maintain plugins and try to map the game world and its contents.

I saw one user, Danj, talking about their work on mapping software and my interest was piqued.

The MUSHclient [bundle](https://github.com/fiendish/aardwolfclientpackage/wiki/) provided by Fiendish has a mapper that tracks your movement through ==rooms== and ==exits==. This data is leveraged by a bunch of plugins in a variety of ways, none of which are super relevant to this post.

In practice, I know that I can't possibly compete with existing solutions like the [Gaardian Roominator](http://rooms.gaardian.com/index.php) and the beta SVG version that I don't have a link to at the moment. That doesn't stop me from wanting to gets my hands on the data and see if I can do anything cool with it.
## The Data
The mapper's map data is stored in an sqlite database, and the schema is pretty straightforward. There's a few tables we care about: [[#Areas]], [[#Rooms]], and [[#Exits]].

These tables look like they have a lot of columns, but most of them end up being irrelevant in the context of trying to create a graph representing the rooms and exits connecting them.

The `exits` table is just a join table on `rooms`, so in theory it should be pretty trivial to assemble a list of vertices ( rooms ) and edges ( exits ) and pump them into graphviz, right?

### Areas
```sql
sqlite> .schema areas
CREATE TABLE areas(
uid TEXT NOT NULL,
name TEXT,
texture TEXT,
color TEXT,
flags TEXT NOT NULL DEFAULT '',
`id` integer,
`created_at` datetime,
`updated_at` datetime,
`deleted_at` datetime,
PRIMARY KEY(uid));
CREATE INDEX `idx_areas_deleted_at` ON `areas`(`deleted_at`);
```
###  Rooms
```sql
sqlite> .schema rooms
CREATE TABLE rooms(
uid TEXT NOT NULL,
name TEXT,
area TEXT,
building TEXT,
terrain TEXT,
info TEXT,
notes TEXT,
x INTEGER,
y INTEGER,
z INTEGER,
norecall INTEGER,
noportal INTEGER,
ignore_exits_mismatch INTEGER NOT NULL DEFAULT 0,
`id` integer,
`created_at` datetime,
`updated_at` datetime,
`deleted_at` datetime,
`flags` text,
PRIMARY KEY(uid));
CREATE INDEX rooms_area_index ON rooms (area);
CREATE INDEX `idx_rooms_deleted_at` ON `rooms`(`deleted_at`);
```

It wasn't until writing this and cleaning up that `CREATE TABLE` statement to be readable did I notice that rooms have integer IDs. That may be useful for solving the problems I'll describe shortly.
###  Exits
```sql
sqlite> .schema exits
CREATE TABLE exits(
dir TEXT NOT NULL,
fromuid TEXT NOT NULL,
touid TEXT NOT NULL,
level STRING NOT NULL DEFAULT '0',
PRIMARY KEY(fromuid, dir));
CREATE INDEX exits_touid_index ON exits (touid);
```
## Almost Right
Getting the edges and vertices into graphviz ended up being pretty trivial. The part that took me the longest was learning how to do database stuff in Go. So far I'd managed to interact with flat files and HTTP requests for getting my data, but I knew that wouldn't last forever.
### A relational tangent
In brief, the Go database workflow has some steps in common:
1) import `database/sql`
2) import your database driver
3) open the database or establish a connection to the server
4) Make a query
5) Scan() into a value
6) use the value

There's some variance with points 5 and 6 on whether you want exactly one or some other number of results ( `Query` vs `QueryRow`) .

To demonstrate, here's a pared down sample of what I'm using in my `aardmapper`.

```go {title="main.go"}
import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

type Area struct {
  Uid, Name, Flags, Color, Texture sql.NullString
}

func main() {
  db, _ := sql.Open("sqlite3", fn)
  // error handling is elided for brevity. do not ignore errors.
}
  a = Area{}
  if err := row.Scan(&a.Uid, &a.Name, &a.Flags, &a.Color, &a.Texture); err != nil {
    if err == sql.ErrNoRows {
      fmt.Fatalf("no area found: %w", err)
    }
  }
  // do stuff with your queried Area
```
##  The graph must grow
Once I was able to query rooms and exits from the database, I was on the fast track. The graphviz API is relatively straightforward when you're using Go:
```go {title="mapper.go"}
gv := graphviz.New()
g := gv.Graph()
for _, room := range rooms { // creation of rooms elided
    origin, _ := g.CreateNode("RoomID_AAAAA")
    dest, _ := g.CreateNode("RoomID_AAAAB")
    edge, _ := g.CreateEdge("connecting AAAAA to AAAAB", origin, dest)
}
// Once again, error handling has been elided for brevity. **Do not ignore errors**.
```

This ended up working great. The rooms and exits matched up to vertices and edges the way I expected.

The only problem was that rendering the entire thing on my DigitalOcean droplet will apparently take more than 24 hours. I had to terminate the process at around 16 hours because I got impatient.
## The lay of the land
This first, naive implementation mostly does the trick. It works really well for smaller quantities of rooms. Below you can see a PNG and SVG rendering of 250 rooms, and the code used to generate it.

```go
  if err := gv.RenderFilename(g, graphviz.SVG, "/var/www/renders.ndumas.com/aardmaps/name.ext"); err != nil {
    log.Fatal(err)
  }
```

{{< figure src="[[Resources/attachments/250-rooms.svg]]" title="250 Rooms (SVG)" alt="a disorderly grid of squares representing rooms connected to each other in a video game" caption="SVG scales a lot better" >}}


{{< figure src="[[Resources/attachments/250-rooms.png]]" title="250 Rooms (PNG)" alt="a disorderly grid of squares representing rooms connected to each other in a video game" caption="Raster images can be simpler and more performant to render" >}}


## What's next?
The current iteration of rendering is really crude:
- The rooms are displayed using their numeric IDs, not human friendly names.
- Rooms are grouped by area, creating subgraphs to describe them will help interpreting the map and probably help rendering.
- The current iteration is very slow

I've also been contemplating the idea of rendering each area one at a time, and then manipulating the resulting SVG to display connections that cross between areas. This would almost certainly be infinitely faster than trying to render 30,00 vertices and 80,000 edges simultaneously.

All my code can be found [here](https://code.ndumas.com/ndumas/aardmapper). It's still early in prototyping so I don't have any stable builds or tags yet.
