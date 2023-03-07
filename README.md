# Test task for GO

## 👋  Summary

The goal is to create an HTTP API for uploading, optimizing, and serving images.

1. The API should expose an endpoint for uploading an image. After uploading, the image should be sent for optimization via a queue (e.g. RabbitMQ) to prevent excessive system load in case of many parallel image uploads and increase the system durability.
2. Uploaded images should be taken from the queue one by one and optimized using the `github.com/h2non/bimg` go package (or `github.com/nfnt/resize` package). For each original image, three smaller-size image variants should be generated and saved, with 75%, 50%, and 25% quality.
3. The API should expose an endpoint for downloading an image by ID. The endpoint should allow specifying the image optimization level using query parameters (e.g. `?quality=100/75/50/25`).

## 🤔  Evaluation criteria

1. **Functionality.** The developed solution should function as described in the "Summary" section. However, if you think that you can create a solution better than described in the "Summary" section, you are welcome to do so.
2. **Code simplicity**. The architecture should be simple and easy to understand, the code should be well-formatted and consistent. Usage of code formatters (like gofmt) and linters (like golangci-lint) is encouraged.

___

## 💥 Realization

Workflow:

![App infrastructure](./img/app.svg)

App structure:

```note
/internal
    /app                    // Application that has all required components

    /delivery               // Package that include infrastructure for communicating by HTTP protocols
        /http                   // http
            /handler                // user defined handler
            /rest                   // REST API methods for handler
            /server                 // http server
        /rabbitmq               // amqp
            /client                 // RabbitMQ client

    /domain                 // Domain business logic
        /dto                    // Data Transfer Object
        /repositories           // Interfaces for services (use-cases)

    /infrastructure         // Actual implementation of components
        /file                   // Local file storage (using standard pkg os / filepath / io/ioutil)
        /worker                 // Background job / service that proceed the image from MessageBroker
            /compressor             // as a part of background job

    /services               // Services that App uses
        /image                  // Image service
        /publisher              // Part of MessageBroker (only Send to...) for HTTP API
```

## ✅ Usage

```note
make rabbit       // Up the RabbitMQ instance from docker
make rabbit-stop  // Stop and delete the container

go run cmd/main.go
```

### OUTPUT

You can see the [examples.log](./examples.log) file for actual output of the application

> Pay attention that in the log I've put the ID, URL and etc. data that are not for production!
>
> It's a test task, so I just skip this part of hiding the USER information (but what someone can steal? ImageID? XD)
>
> REMOVE the ID, URL and other information or CROP it

### How can I improve the application

1. WRITE UNIT/INTEGRATION/E2E TESTS (but it takes TOO MUCH time)
2. Add database for register the errors in background job, if user wants to know which problems his/her image has
3. Redirect errors from the Goroutines
4. Remove HARDCODE and add ENV setup or flags of app
5. USE the specific exchanger/channel/queue of RabbitMQ
6. Use relative path
7. Use caching (if need)
8. . . .

### 🟡 Notes

I don't use the `github.com/h2non/bimg` package due to its dependency for Linux.

I've written this application on Windows platform. I could use the WSL2 for example, but it to long for open it.

Hope it won't be a problem :)

## 💯 Testing

I've using the [Thunder client](https://marketplace.visualstudio.com/items?itemName=rangav.vscode-thunder-client) (`rangav.vscode-thunder-client`) extension from VS Code to proceed HTTP requests

### ❗ NO AUTOMATION TESTS

The HTTP methods collections that have been used placed there: [collection.json](./thunder-collection_RabbitMQ%20Image.json)

> PAY attention: in POST method you need replace the path to the image file in the body

```json
"body": {
        "type": "formdata",
        "raw": "",
        "form": [],
        "files": [
            {
                "name": "image",
                "value": "<PATH TO IMAGE>"
            }
        ]
    },
```

> For GET HTTP method:

```json
"url": "localhost:8080/img/<UUID OF IMAGE>"
```

> Use the json file for import:

![Import Collection for Thunder Client](./img/ImportThunderCollection.png)

## 👁‍🗨 Feedback

### **Disadvantages**

1. [Name variables in camelCase](https://github.com/andrsj/go-rabbit-image/blob/db9c4599659a6815e665133daa6a53f95c10c78f/internal/app/app.go#L27).
2. [Name local variables in lower case, in order not to violate the principle of encapsulation (when it is small, it is not possible to import)](https://github.com/andrsj/go-rabbit-image/blob/db9c4599659a6815e665133daa6a53f95c10c78f/internal/app/app.go#L30).
3. [Do not write excessive comments that only add noise to the code](https://github.com/andrsj/go-rabbit-image/blob/db9c4599659a6815e665133daa6a53f95c10c78f/internal/delivery/rabbitmq/client/client.go#L75). TIP: [Tips for writing comments from Clean Code](https://gist.github.com/wojteklu/73c6914cc446146b8b533c0988cf8d29#comments-rules).
4. Go prefers large files. `http/rest/api/get.go` and `http/rest/api/post.go` can be combined into one file. If other endpoints were to be added in the future, it would not be convenient to keep all GETs in one file and all POSTs in another.
5. Use a single style for error checking: [First](https://github.com/andrsj/go-rabbit-image/blob/db9c4599659a6815e665133daa6a53f95c10c78f/internal/services/image/storage/storage.go#L32), [Second](https://github.com/andrsj/go-rabbit-image/blob/db9c4599659a6815e665133daa6a53f95c10c78f/internal/services/image/storage/storage.go#L51).
6. `domain/repositories // Interfaces for services (use-cases)`. Confusing, because there is [pattern repository](https://medium.com/@pererikbergman/repository-design-pattern-e28c0f3e4a30) which is responsible for accessing the data.
7. **If it were necessary to use a queue in another project with an implementation of RabbitMq, how would it happen?**
8. By design, the `infrastructure` level s responsible for implementations, but there is a `Compressor` interface.
It is better to create a separate `compressor` package with the `resize` implementation as well as the queue.
9. The `infrastructure/worker/utils.go` file named `utils` is considered bad practice because it collects different functions from different places that could actually be put into separate packages without creating chaos in that file. It could be exported, for example, to `pkg/image`.
10. It is not obvious that in [`Worker.Start()`](https://github.com/andrsj/go-rabbit-image/blob/db9c4599659a6815e665133daa6a53f95c10c78f/internal/infrastructure/worker/work.go#L18) the logic of working with the image is prescribed. It would be better if the worker was responsible only for receiving messages and transferring them to another handler, which he does not know about at all. That is, a callback would be passed to `Start`, for example, which would call a service for image processing.
11. **What are your business rules in the application? Where would you put these rules?**
12. Image quality validation occurs at the http level. If `ReadImageFromStorage` is not called via http, the validation will no longer work.

### **Improvements**

Queue refactor (refers to item 7):

1. Create a separate `pkg/queue' package.
2. Move the `MessageBroker` interface to the `pkg/queue/queue.go` package (now it is in `domain/repositories/queue/interface.go`).
3. `MessageDTO` is a generic message that should be sent by the queue and it is correct, but it should be moved closer to where it is used in `pkg/queue/queue.go` (it is currently in `domain/dto/dto.go`) .
4. Move the rabbitmq implementation to `pkg/queue/rabbitmq.go`. Now the queue can be easily reused.

### **Advantages:**

1. No errors from golangci-lint.
2. There is a cool README, with a diagram and suggestions for improvement.
3. There is a docker.
4. The test is not abandoned by one commit. Although the name of the commits could be made more meaningful by using [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).
5. There are interfaces with further implementation.
6. There is logging at all levels.
7. `domain/repositories/queue`. The queue is transport. And this is correct, because data is not stored in the queue, data is simply passed through it in the same way as rest. Input-Output.

### **Literature**

1. [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
2. [SOLID Go Design by Dave Cheney](https://dave.cheney.net/2016/08/20/solid-go-design)
3. [A primer on the clean architecture pattern and its principles](https://www.techtarget.com/searchapparchitecture/tip/A-primer-on-the-clean-architecture-pattern-and-its-principles)
4. [Clean code by Robert Martin](https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882)
5. [Practical Go by Dave Cheney](https://dave.cheney.net/practical-go/presentations/qcon-china.html)
