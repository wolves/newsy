<p align="center"><img src="https://raw.githubusercontent.com/wolves/newsy/main/logo.svg" /></p>

# Newsy

*A news aggregation service that puts the user back in control.*

Newsy is a simple service that enables the aggregation of news from sources of the user's choosing. Articles from these sources can then be filtered by category for more focused consumption. No more fancy unseen algorithms deciding what news makes it into your feed. 

## Overview
Newsy provides an API that can be imported and used in your projects and a
command-line application for interfacing with your personal newsfeed.

*The following sub-sections provide a general description of the components, if
you are looking for a more detailed documentation, it can be found in the [Usage
section below](#usage-docs).*

### Newsy Service

The news service is in charge of managing the overall news system. It is
responsible for:

- Starting
  - Starting the service activates any saved state from the previous session as
    well as activating a process that will periodically update that state with
    current articles at an interval configured by the user.
- Stopping
  - Stopping cancels all running processes and takes care of resource cleanup
    for the service. When the service is stopped it also saves the last state
    at the time of stopping.
- Setting the State autosave frequency
- Searching Article history
- Error delivery and management
- Saving State
  - The historical state of articles is saved in JSON format in a local file on
    the system. This also provides the user with the ability to search and
    access previously ingested articles.
- Loading State
  - Any previouly saved state is loaded from the local JSON backup file at the
    time the service is started.

### Subscriptions

A user can subscribe to the news service via one or more categories of their choosing in order to
receive news that is based in that category. The user can unsubscribe from a category at anytime.
([*Usage Docs*](#subscriber-docs))

### Source/Sources

A news `Source` can be added and removed at anytime and are also able to
publish to any `Category` they wish to define at any interval they choose.
([*Usage Docs*](#source-docs))

### Article/Articles

An `Article` represents a single story unit delivered by a `Source` to a
`Subscriber` that is subscribed to its associated `Category`
([*Usage Docs*](#article-docs))

### Category/Categories

A `Category` is a user-friendly classification describing the subject matter
and topics found within its corresponding `Article` and is used to provide
those stories to those that have subscribed to it.
([*Usage Docs*](#category-docs))

### Errors

A set of custom Errors has also been defined to provide more descriptive error
messaging in the case that one is returned.
([*Usage Docs*](#error-docs))

## Installation

### Go Install
*Coming Soon*

### Source (via go get)
*Coming Soon*

### Package Manager
*Stretch Goal - Maybe Coming Soon* :)

# Usage
<a id="usage-docs"></a>

## Newsy Service

```Go
type Newsy struct {
  cancel       context.CancelFunc
  errs         chan error
  stopped      bool
  stateFile    string
  saveInterval time.Time
}
```

### Start
```Go
func (n *Newsy) Start(ctx context.Context) (context.Context, error) {}
```
`Newsy.Start` starts up the news service along with any required processes and
loads your previous state from the `Newsy.stateFile` if any previous state
exists. Once the service is started snapshots of the existing service state are
saved to `Newsy.stateFile` at an interval defined by `Newsy.SaveInterval()`.

### Stop
```Go
func (n *Newsy) Stop() {}
```
`Newsy.Stop` stops the news service along with canceling all running processes.
As part of the stopping process, the services existing state is saved to the
`Newsy.stateFile` and all other running processes are signaled to cleanup or
cleaned up by the main `Newsy` service itself.

Once `Newsy` is stopped any further interaction with it from other components
of the program will trigger an `ErrNewsyServiceStopped` error.

### SaveInterval
```Go
func (n *Newsy) SaveInterval(t time.Time) {}
```
`Newsy.SaveInterval` takes a `Time` value for the frequency which the service's
state should be auto-saved to the file defined by `Newsy.stateFile`. It also
sets that value for the `Newsy.saveInterval` field.

### Search
```Go
func (n *Newsy) Search(categories ...Category, ids ...ArticleId) (Articles, error) {}
```
`Newsy.Search` takes one or more `Categories` along with one or more
`ArticleIds` and searches the news service history for any matching `Articles`.

### Errors
```Go
func (n *Newsy) Errors() []error {}
```
`Newsy.Errors` provides insight into existing errors within the service.

### saveState
```Go
func (n *Newsy) saveState() {}
```
`Newsy.saveState` saves the current state of the news service in JSON format to
the file defined by `Newsy.stateFile`. It is utilized by both the auto-save
functionality as well as the state saving process that occurs when `Newsy.Stop`
is called.

### loadState
```Go
func (n *Newsy) loadState() {}
```
`Newsy.loadState` takes the JSON state file defined by `Newsy.stateFile` and
loads it into the news service. This takes place at the time the service is
started via `Newsy.Start`.


## Subscriber
<a id="subscriber-docs"></a>

*Equivialant to a "user"?*
- If Yes does that mean it is the same as the "end user" that can stop the Newsy service?
- If No how does it differ? Does the end user somehow become a "Subscriber" by
  subscribing to a category through the Newsy service

**Operating under the assumption that a User subscribes to categories to become a "Subscriber"**

`Subscriber.Subscribe(categories ...categories)`

`Subscriber.Unsubscribe(categories ...categories)`


## Source
<a id="source-docs"></a>

```Go
type NewsSource interface {
  Publish()
  Schedule()
}

// NOTE: Do I want to sub categorize different source types? (eg. FileSource / WebSource / RssSource)
type Source struct {
  name     string
  location string
  schedule time.Time
  articles Articles
}
```

### Publish

*__TODO__: Sort out the proper method signature*
```Go
func (s NewsSource) Publish() {}
```
`Source.Publish()` distributes the `Source`'s new `Articles` for `Newsy` to distribute to the proper `Subscribers`

### Schedule
*(Not quite sure if this will be needed)*

*__TODO__: Sort out the proper method signature*
```Go
func (s NewsSource) Schedule() {}
```
`Source.Schedule()` provides the frequency/interval at which the `Source`
regularly publishes its `Articles` 


## Article
<a id="article-docs"></a>

```Go
type articleId int

type Article struct {
  id         articleId
  title      string
  author     string
  content    string
  categories Categories
  source     Source
}

type Articles []Article
```

## Category/Categories
<a id="category-docs"></a>

```Go
type Category string

type Categories []Category
```
`Category` is a string that is used to define the `Category` of an `Article`

### CategoryArticles

```Go
type CategoryArticles map[Category]Articles
```
`CategoryArticles` is a map that uses a `Category` as a key and returns associated `Articles`


## Errors *(WIP)*
<a id="error-docs"></a>

*These are the initial errors I can assume will be implemented. As the
program is built out more may be required and/or these existing `Error`'s may
need modification*

### ErrInvalidSource

```Go
type ErrInvalidSource string
```
`ErrInvalidSource` is returned when the news `Source` is invalid due to an incorrect definition

### ErrInvalidArticle

```Go
type ErrInvalidArticle string
```
`ErrInvalidArticle` is returned when an `Article` is invalid due to an invalid attribute/field

### ErrInvalidCategory

```Go
type ErrInvalidCategory string
```
`ErrInvalidCategory` is returned when a `Category` is created with an incorrect type

### ErrNewsyServiceStopped  

```Go
type ErrNewsyServiceStopped struct{}
```
`ErrNewsyServiceStopped` is returned if a `Newsy` component attempts to interact with it after it has been stopped


# CLI
*Coming Soon...*

# Testing

After cloning the repository, tests can be run with the following:
```Go
go test ./... -v -cover -race
```
