# go-couch - CouchService Go SDK

<p align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/rozanecm/go-couch)](https://goreportcard.com/report/github.com/rozanecm/go-couch)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
</p>

![Logo](logo.png)

Welcome to the Go SDK for interacting with CouchDB services. This SDK provides a convenient way to manage
databases, documents, and design documents using Go programming language.

<div style="background-color: #ffffcc; padding: 10px; border-radius: 5px;">

⚠️ **Warning: Breaking Changes and Testing Status**

Please note that this project is under active development and testing. As a result, breaking changes may occur between
versions. While we strive for stability and reliability, we cannot guarantee error-free execution in all scenarios.

We recommend reviewing the release notes and testing updates in a controlled environment before deploying this project
in a production setting.

Thank you for your understanding and collaboration.
</div>

## Features

- **Database Management**: Create, retrieve, update, and delete databases.
- **Document Operations**: Create, retrieve, update, and delete documents within databases.
- **Design Document Management**: Create and update design documents with views.
- **Query Views**: Perform queries on database views using customizable parameters to retrieve filtered and sorted
  results.

## Installation

To use the CouchService Go SDK, you need to have Go installed on your system. Then, you can install the SDK using the
following command:

```bash
go get github.com/rozanecm/go-couch
```

## Usage

Here are some examples demonstrating how to use the CouchService Go SDK:

### Create a CouchService Instance

```go
import (
"context"
"github.com/rozanecm/go-couch"
)

func main() {
baseURL := "http://localhost:5984" // Example CouchDB URL
cs := couchservice.GetInstance(baseURL, "username", "password")

// Use cs to perform operations like GetDB, CreateDoc, etc.
}
```

### Create a Database

```go
ctx := context.Background()
dbName := "example_db"
createIfNotExist := true

db, err := cs.GetDB(ctx, dbName, createIfNotExist)
if err != nil {
panic(err)
}
```

### Create a Document

```go
ctx := context.Background()
docData := map[string]interface{}{
"name": "John Doe",
"age": 30,
}

resp, err := db.CreateDoc(ctx, docData)
if err != nil {
panic(err)
}

fmt.Println("Document created successfully with ID:", resp.ID)
```

### Update a Document

```go
ctx := context.Background()
docID := "document_id"
updatedData := map[string]interface{}{
"_id": docID,
"_rev": "current_revision",
"key": "new_value",
}

err := db.UpdateDoc(ctx, docID, updatedData)
if err != nil {
panic(err)
}

fmt.Println("Document updated successfully")
```

## Contributing Guidelines

We welcome contributions to improve and extend the functionality of this SDK. When contributing, please follow these
guidelines:

- Respect the existing code style and structure.
- Add unit tests for new features or bug fixes.
- Provide clear and informative commit messages.
- Follow the Contributor Covenant Code of Conduct.

## Usage of ViewResponse Struct

The `ViewResponse` struct is designed to represent the response JSON object returned from a database view. This struct
can be used as a generic resultVar in the `View` function of the `Database` type. The `View` function unmarshalls the
view results directly into a struct provided by the user, as long as the struct meets the following requirements:

- It must have a field named "Rows" holding a slice of structs with "id" and "key" JSON fields.
- If `params.IncludeDocs` is true, the struct must also have a "doc" JSON field.

Thank you for contributing to make our SDK better!

---

Feel free to reach out if you have any questions or need assistance with the SDK.

Happy coding!
