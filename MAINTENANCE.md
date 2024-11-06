# Maintenance Guide

## Service Interface

Every service in Queuer must implement the following interface:

```go
type Service interface {
    json.Unmarshaler
    fmt.Stringer
    Run() (results chan []byte, errs chan error)
}
```

### Service Implementation

Each service needs to implement this interface, where:

- `UnmarshalJSON` is used for parsing the incoming task data.
- `String` returns the name of the service (as a string).
- `Run` processes the task, returning channels for both results and errors.

### Example Service Implementation

Here’s a simple example of an `Adder` service:

```go
package services

import (
    "encoding/json"
    "errors"
)

func NewAdder() *Adder {
    return &Adder{}
}

type Adder struct {
    Addends []int `json:"addends"`
}

func (a *Adder) UnmarshalJSON(data []byte) error {
    temp := struct {
        Addends []*int `json:"addends"`
    }{}

    if err := json.Unmarshal(data, &temp); err != nil {
        return err
    }

    for _, addend := range temp.Addends {
        if addend == nil {
            return errors.New("addend is required")
        }
        a.Addends = append(a.Addends, *addend)
    }

    return nil
}

func (a *Adder) String() string {
    return "adder"
}

func (a *Adder) Run() (chan []byte, chan error) {
    results := make(chan []byte)
    errs := make(chan error)
    go func() {
        sum := 0

        for _, addend := range a.Addends {
            sum += addend
        }

        result, err := json.Marshal(sum)
        if err != nil {
            errs <- err
            return
        }

        results <- result
    }()

    return results, errs
}
```

### Registering Services

All services should be added to the `internal/service.go` file in the `ToService` function. This function maps the string name of each service to its corresponding service initialization function.

Example:

```go
func ToService(s string) Service {
    switch s {
    case "squarer":
        return services.NewSquarer()
    case "adder":
        return services.NewAdder()
    case "longrunner":
        return services.NewLongRunner()
    default:
        return nil
    }
}
```

This allows the program to dynamically map a service string (e.g., `"adder"`) to the corresponding service implementation. The `ToService` function should be updated whenever new services are added.

### Service Location

All service implementations should be placed in the `internal/services/` directory. This keeps the code modular and easy to maintain.

## Service Lifecycle

- **Initialization**: Services are initialized via the `ToService` function, which is responsible for creating an instance of the service when requested.
- **Task Processing**: Each service will process tasks asynchronously in its own goroutine (as shown in the `Run()` method example). Results and errors are sent over channels.

By following this structure, Queuer can be easily extended to include new services simply by implementing the service interface and adding it to `internal/service.go`.

